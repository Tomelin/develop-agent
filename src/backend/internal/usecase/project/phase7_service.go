package project

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	domain "github.com/develop-agent/backend/internal/domain/project"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Phase7AuditInput struct {
	BackendDir     string `json:"backend_dir"`
	FrontendDir    string `json:"frontend_dir"`
	ProjectRootDir string `json:"project_root_dir"`
	HighRetryCount int    `json:"high_retry_count"`
}

type Phase7Service struct {
	files      domain.CodeFileRepository
	runner     commandRunner
	autoReject func(ctx context.Context, projectID, ownerID string, report domain.CatastrophicFailureReport) error
}

func NewPhase7Service(files domain.CodeFileRepository, autoReject func(context.Context, string, string, domain.CatastrophicFailureReport) error) *Phase7Service {
	return &Phase7Service{files: files, runner: execRunner{}, autoReject: autoReject}
}

func (s *Phase7Service) RunAudit(ctx context.Context, projectID, ownerID string, input Phase7AuditInput) (*domain.SecurityAuditReport, error) {
	backend := strings.TrimSpace(input.BackendDir)
	if backend == "" {
		backend = "."
	}
	root := strings.TrimSpace(input.ProjectRootDir)
	if root == "" {
		root = backend
	}
	retries := input.HighRetryCount
	if retries <= 0 {
		retries = 2
	}

	report := &domain.SecurityAuditReport{GeneratedAt: time.Now().UTC()}
	findings := make([]domain.SecurityFinding, 0)

	gosecFindings := s.scanGosec(ctx, backend)
	report.StaticAnalysis = append(report.StaticAnalysis, gosecFindings...)
	findings = append(findings, gosecFindings...)

	govulnFindings := s.scanGovulncheck(ctx, backend)
	npmFindings := s.scanNpmAudit(ctx, strings.TrimSpace(input.FrontendDir))
	report.Dependencies = append(report.Dependencies, govulnFindings...)
	report.Dependencies = append(report.Dependencies, npmFindings...)
	findings = append(findings, govulnFindings...)
	findings = append(findings, npmFindings...)

	secretFindings := s.scanSecrets(ctx, root)
	findings = append(findings, secretFindings...)

	configFindings := s.checkSecurityConfiguration(ctx, backend)
	findings = append(findings, configFindings...)

	for i := range findings {
		if findings[i].Status == "" {
			findings[i].Status = domain.FindingStatusOpen
		}
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].CVSS == findings[j].CVSS {
			return findings[i].ID < findings[j].ID
		}
		return findings[i].CVSS > findings[j].CVSS
	})

	report.Findings = findings
	report.Summary = summarizeSecurityFindings(findings)
	if report.Summary.Score >= 80 {
		report.Status = "PASS"
	} else {
		report.Status = "FAIL"
	}
	report.AutoRejection = evaluateSecurityAutoRejection(findings, retries)

	if report.AutoRejection.Triggered && s.autoReject != nil {
		desc := append([]string{}, report.AutoRejection.Findings...)
		maxCVSS := 0.0
		for _, f := range findings {
			if f.CVSS > maxCVSS {
				maxCVSS = f.CVSS
			}
		}
		_ = s.autoReject(ctx, projectID, ownerID, domain.CatastrophicFailureReport{
			MaxCVSS:             maxCVSS,
			CredentialsExposed:  containsCategory(findings, "Secrets"),
			FailureDescriptions: desc,
			SourcePhase:         7,
		})
	}

	if err := s.persistSecurityReport(ctx, projectID, report); err != nil {
		return nil, err
	}

	return report, nil
}

func (s *Phase7Service) persistSecurityReport(ctx context.Context, projectID string, report *domain.SecurityAuditReport) error {
	projectOID, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return err
	}
	content := BuildSecurityAuditMarkdown(report)
	return s.files.Upsert(ctx, &domain.CodeFile{
		ProjectID:   projectOID,
		Path:        "reports/SECURITY_AUDIT.md",
		TaskID:      "TASK-12-007",
		Language:    "markdown",
		PhaseNumber: 7,
		Content:     content,
	})
}

func summarizeSecurityFindings(findings []domain.SecurityFinding) domain.SecurityAuditSummary {
	s := domain.SecurityAuditSummary{TotalFindings: len(findings), Score: 100}
	for _, f := range findings {
		switch f.Severity {
		case domain.SecuritySeverityCritical:
			s.CriticalCount++
			s.Score -= 20
		case domain.SecuritySeverityHigh:
			s.HighCount++
			s.Score -= 10
		case domain.SecuritySeverityMedium:
			s.MediumCount++
			s.Score -= 5
		case domain.SecuritySeverityLow:
			s.LowCount++
			s.Score -= 2
		}
	}
	if s.Score < 0 {
		s.Score = 0
	}
	return s
}

func evaluateSecurityAutoRejection(findings []domain.SecurityFinding, retries int) domain.SecurityAutoRejectionResult {
	res := domain.SecurityAutoRejectionResult{RetryCount: retries}
	highCount := 0
	for _, f := range findings {
		if strings.EqualFold(f.Category, "Secrets") {
			res.Triggered = true
			res.Reason = "secret-exposed"
			res.ReturnedPhase5 = true
			res.Findings = append(res.Findings, f.ID+": exposed secret")
			return res
		}
		if f.CVSS >= 9.0 {
			res.Triggered = true
			res.Reason = "critical-cvss"
			res.ReturnedPhase5 = true
			res.Findings = append(res.Findings, fmt.Sprintf("%s: cvss %.1f", f.ID, f.CVSS))
			return res
		}
		if strings.Contains(strings.ToUpper(f.Title), "IDOR") || strings.Contains(strings.ToUpper(f.Title), "BROKEN ACCESS CONTROL") {
			res.Triggered = true
			res.Reason = "broken-access-control"
			res.ReturnedPhase5 = true
			res.Findings = append(res.Findings, f.ID+": structural authz failure")
			return res
		}
		if f.CVSS >= 7.0 {
			highCount++
		}
	}
	if highCount > 0 && retries <= 0 {
		res.Triggered = true
		res.Reason = "high-vulns-exhausted-retries"
		res.ReturnedPhase5 = true
	}
	return res
}

func BuildSecurityAuditMarkdown(report *domain.SecurityAuditReport) string {
	var b strings.Builder
	b.WriteString("# SECURITY_AUDIT\n\n")
	b.WriteString(fmt.Sprintf("Gerado em: `%s`\n\n", report.GeneratedAt.Format(time.RFC3339)))
	b.WriteString("## Executive Summary\n\n")
	b.WriteString(fmt.Sprintf("- Score de segurança: **%d/100**\n", report.Summary.Score))
	b.WriteString(fmt.Sprintf("- Status geral: **%s**\n", report.Status))
	b.WriteString(fmt.Sprintf("- Findings: CRITICAL=%d, HIGH=%d, MEDIUM=%d, LOW=%d\n", report.Summary.CriticalCount, report.Summary.HighCount, report.Summary.MediumCount, report.Summary.LowCount))
	b.WriteString("\n## Findings\n\n")
	b.WriteString("| ID | Título | Severidade | CVSS | Status | Detector |\n")
	b.WriteString("|----|--------|------------|------|--------|----------|\n")
	for _, f := range report.Findings {
		b.WriteString(fmt.Sprintf("| %s | %s | %s | %.1f | %s | %s |\n", f.ID, sanitizeCell(f.Title), f.Severity, f.CVSS, f.Status, sanitizeCell(f.DetectedBy)))
	}
	b.WriteString("\n## Detalhes por Finding\n\n")
	for _, f := range report.Findings {
		b.WriteString(fmt.Sprintf("### %s — %s\n\n", f.ID, f.Title))
		b.WriteString(fmt.Sprintf("- Categoria: %s\n", f.Category))
		b.WriteString(fmt.Sprintf("- Severidade/CVSS: %s / %.1f\n", f.Severity, f.CVSS))
		if f.CVE != "" {
			b.WriteString(fmt.Sprintf("- CVE: %s\n", f.CVE))
		}
		if f.File != "" {
			b.WriteString(fmt.Sprintf("- Arquivo/Linha: `%s:%d`\n", f.File, f.Line))
		}
		if f.PoC != "" {
			b.WriteString(fmt.Sprintf("- PoC: %s\n", f.PoC))
		}
		if f.Remediation != "" {
			b.WriteString(fmt.Sprintf("- Remediação: %s\n", f.Remediation))
		}
		b.WriteString("\n")
	}
	b.WriteString("## Auto Rejection Trigger\n\n")
	b.WriteString(fmt.Sprintf("- Triggered: **%t**\n", report.AutoRejection.Triggered))
	b.WriteString(fmt.Sprintf("- Reason: `%s`\n", report.AutoRejection.Reason))
	b.WriteString(fmt.Sprintf("- Retornou para fase 5: **%t**\n", report.AutoRejection.ReturnedPhase5))
	return b.String()
}

func sanitizeCell(in string) string {
	in = strings.ReplaceAll(in, "|", "/")
	in = strings.ReplaceAll(in, "\n", " ")
	return strings.TrimSpace(in)
}

func containsCategory(findings []domain.SecurityFinding, category string) bool {
	for _, f := range findings {
		if strings.EqualFold(f.Category, category) {
			return true
		}
	}
	return false
}

func (s *Phase7Service) scanGosec(ctx context.Context, dir string) []domain.SecurityFinding {
	out, err := s.runner.Run(ctx, dir, "gosec", "-fmt", "json", "./...")
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	type issue struct {
		RuleID   string `json:"rule_id"`
		Details  string `json:"details"`
		Severity string `json:"severity"`
		File     string `json:"file"`
		Line     string `json:"line"`
	}
	var parsed struct {
		Issues []issue `json:"Issues"`
	}
	if json.Unmarshal([]byte(out), &parsed) != nil {
		return nil
	}
	findings := make([]domain.SecurityFinding, 0, len(parsed.Issues))
	for i, it := range parsed.Issues {
		line, _ := strconv.Atoi(it.Line)
		sev, cvss := gosecSeverity(it.Severity)
		findings = append(findings, domain.SecurityFinding{
			ID:          fmt.Sprintf("GOSEC-%03d", i+1),
			Title:       it.RuleID,
			Category:    "Static Analysis",
			Severity:    sev,
			CVSS:        cvss,
			Description: it.Details,
			File:        it.File,
			Line:        line,
			DetectedBy:  "gosec",
			Status:      domain.FindingStatusOpen,
		})
	}
	return findings
}

func gosecSeverity(in string) (domain.SecuritySeverity, float64) {
	switch strings.ToUpper(strings.TrimSpace(in)) {
	case "HIGH":
		return domain.SecuritySeverityHigh, 8.0
	case "MEDIUM":
		return domain.SecuritySeverityMedium, 5.5
	default:
		return domain.SecuritySeverityLow, 3.5
	}
}

func (s *Phase7Service) scanGovulncheck(ctx context.Context, dir string) []domain.SecurityFinding {
	out, err := s.runner.Run(ctx, dir, "govulncheck", "-json", "./...")
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	lines := strings.Split(out, "\n")
	findings := make([]domain.SecurityFinding, 0)
	idx := 1
	for _, line := range lines {
		if !strings.Contains(line, "GO-") || !strings.Contains(line, "module") {
			continue
		}
		cve := extractGovulnID(line)
		if cve == "" {
			continue
		}
		findings = append(findings, domain.SecurityFinding{
			ID:          fmt.Sprintf("GOVULN-%03d", idx),
			Title:       "Dependência Go vulnerável",
			Category:    "Dependencies",
			Severity:    domain.SecuritySeverityHigh,
			CVSS:        7.5,
			CVE:         cve,
			Description: strings.TrimSpace(line),
			DetectedBy:  "govulncheck",
			Status:      domain.FindingStatusOpen,
		})
		idx++
	}
	return findings
}

func extractGovulnID(line string) string {
	re := regexp.MustCompile(`GO-[0-9]{4}-[0-9]+`)
	return re.FindString(line)
}

func (s *Phase7Service) scanNpmAudit(ctx context.Context, dir string) []domain.SecurityFinding {
	if strings.TrimSpace(dir) == "" {
		return nil
	}
	out, err := s.runner.Run(ctx, dir, "npm", "audit", "--json")
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	var parsed map[string]any
	if json.Unmarshal([]byte(out), &parsed) != nil {
		return nil
	}
	vulns, ok := parsed["vulnerabilities"].(map[string]any)
	if !ok {
		return nil
	}
	findings := make([]domain.SecurityFinding, 0)
	idx := 1
	for pkg, raw := range vulns {
		item, _ := raw.(map[string]any)
		severity := fmt.Sprintf("%v", item["severity"])
		sev, cvss := npmSeverity(severity)
		title := fmt.Sprintf("Pacote npm vulnerável: %s", pkg)
		fix := fmt.Sprintf("%v", item["fixAvailable"])
		findings = append(findings, domain.SecurityFinding{
			ID:          fmt.Sprintf("NPM-%03d", idx),
			Title:       title,
			Category:    "Dependencies",
			Severity:    sev,
			CVSS:        cvss,
			Description: "Dependência reportada no npm audit",
			Remediation: "Aplicar atualização sugerida pelo npm audit: " + fix,
			DetectedBy:  "npm-audit",
			Status:      domain.FindingStatusOpen,
		})
		idx++
	}
	return findings
}

func npmSeverity(in string) (domain.SecuritySeverity, float64) {
	switch strings.ToLower(strings.TrimSpace(in)) {
	case "critical":
		return domain.SecuritySeverityCritical, 9.5
	case "high":
		return domain.SecuritySeverityHigh, 8.0
	case "moderate":
		return domain.SecuritySeverityMedium, 5.5
	default:
		return domain.SecuritySeverityLow, 3.0
	}
}

func (s *Phase7Service) scanSecrets(ctx context.Context, root string) []domain.SecurityFinding {
	if strings.TrimSpace(root) == "" {
		root = "."
	}
	findings := make([]domain.SecurityFinding, 0)
	if out, err := s.runner.Run(ctx, root, "gitleaks", "detect", "--no-git", "--report-format", "json", "--report-path", "/tmp/gitleaks.json"); err == nil || strings.Contains(strings.ToLower(out), "leaks") {
		parsed := parseGitleaksFromFile("/tmp/gitleaks.json")
		if len(parsed) > 0 {
			return parsed
		}
	}
	out, err := s.runner.Run(ctx, root, "trufflehog", "filesystem", ".", "--json")
	if err != nil || strings.TrimSpace(out) == "" {
		return findings
	}
	idx := 1
	for _, line := range strings.Split(out, "\n") {
		if !strings.Contains(line, "DetectorName") {
			continue
		}
		findings = append(findings, domain.SecurityFinding{
			ID:          fmt.Sprintf("SECRET-%03d", idx),
			Title:       "Secret hardcoded detectado",
			Category:    "Secrets",
			Severity:    domain.SecuritySeverityCritical,
			CVSS:        10.0,
			Description: strings.TrimSpace(line),
			PoC:         "Reprodução: execute trufflehog filesystem . --json",
			DetectedBy:  "trufflehog",
			Status:      domain.FindingStatusOpen,
		})
		idx++
	}
	return findings
}

func parseGitleaksFromFile(path string) []domain.SecurityFinding {
	data, err := os.ReadFile(path) // #nosec G304 -- path is generated internally for scanner report file
	if err != nil || len(data) == 0 {
		return nil
	}
	var items []map[string]any
	if json.Unmarshal(data, &items) != nil {
		return nil
	}
	findings := make([]domain.SecurityFinding, 0, len(items))
	for i, it := range items {
		line := intFromAny(it["StartLine"])
		file := fmt.Sprintf("%v", it["File"])
		desc := fmt.Sprintf("rule=%v", it["RuleID"])
		findings = append(findings, domain.SecurityFinding{
			ID:          fmt.Sprintf("SECRET-%03d", i+1),
			Title:       "Secret hardcoded detectado",
			Category:    "Secrets",
			Severity:    domain.SecuritySeverityCritical,
			CVSS:        10.0,
			Description: desc,
			File:        file,
			Line:        line,
			DetectedBy:  "gitleaks",
			Status:      domain.FindingStatusOpen,
		})
	}
	return findings
}

func intFromAny(v any) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	default:
		return 0
	}
}

func (s *Phase7Service) checkSecurityConfiguration(_ context.Context, backendDir string) []domain.SecurityFinding {
	findings := make([]domain.SecurityFinding, 0)
	serverFile := filepath.Join(backendDir, "api/server/server.go")
	content, err := os.ReadFile(serverFile) // #nosec G304 -- file path is deterministic inside backend source tree
	if err != nil {
		return findings
	}
	text := string(content)
	if strings.Contains(text, "Access-Control-Allow-Origin") && strings.Contains(text, "*") {
		findings = append(findings, domain.SecurityFinding{
			ID:          "CFG-001",
			Title:       "CORS com wildcard em produção",
			Category:    "Security Misconfiguration",
			Severity:    domain.SecuritySeverityHigh,
			CVSS:        7.2,
			Description: "Header Access-Control-Allow-Origin está com '*' e deve ser restrito por config.",
			File:        serverFile,
			Remediation: "Substituir wildcard por allowlist explícita de origens.",
			DetectedBy:  "manual-check",
			Status:      domain.FindingStatusOpen,
		})
	}
	securityHeaders := []string{"Content-Security-Policy", "X-Frame-Options", "X-Content-Type-Options", "Strict-Transport-Security", "Referrer-Policy"}
	for _, header := range securityHeaders {
		if !strings.Contains(text, header) {
			id := strings.ToUpper(strings.ReplaceAll(header, "-", ""))
			if len(id) > 6 {
				id = id[:6]
			}
			findings = append(findings, domain.SecurityFinding{
				ID:          "CFG-" + id,
				Title:       "Header de segurança ausente: " + header,
				Category:    "Security Misconfiguration",
				Severity:    domain.SecuritySeverityMedium,
				CVSS:        5.0,
				Description: "Middleware HTTP não define o header obrigatório de hardening.",
				File:        serverFile,
				Remediation: "Adicionar header no middleware de resposta HTTP.",
				DetectedBy:  "manual-check",
				Status:      domain.FindingStatusOpen,
			})
		}
	}
	return findings
}
