package project

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	domain "github.com/develop-agent/backend/internal/domain/project"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type commandRunner interface {
	Run(ctx context.Context, dir, name string, args ...string) (string, error)
}

type execRunner struct{}

func (execRunner) Run(ctx context.Context, dir, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

type Phase6Service struct {
	files      domain.CodeFileRepository
	runner     commandRunner
	autoReject func(ctx context.Context, projectID, ownerID string, report domain.CatastrophicFailureReport) error
}

func NewPhase6Service(files domain.CodeFileRepository, autoReject func(context.Context, string, string, domain.CatastrophicFailureReport) error) *Phase6Service {
	return &Phase6Service{files: files, runner: execRunner{}, autoReject: autoReject}
}

func (s *Phase6Service) AnalyzeCoverage(ctx context.Context, projectID, ownerID, backendDir string, threshold float64) (*domain.Phase6CoverageReport, bool, error) {
	if threshold <= 0 {
		threshold = 80
	}
	if backendDir == "" {
		backendDir = "."
	}

	goTestOut, err := s.runner.Run(ctx, backendDir, "go", "test", "./...", "-cover", "-coverprofile=coverage.out")
	if err != nil {
		return nil, false, fmt.Errorf("go test coverage failed: %w\n%s", err, goTestOut)
	}
	funcOut, err := s.runner.Run(ctx, backendDir, "go", "tool", "cover", "-func=coverage.out")
	if err != nil {
		return nil, false, fmt.Errorf("go tool cover failed: %w\n%s", err, funcOut)
	}

	report := &domain.Phase6CoverageReport{
		GeneratedAt:      time.Now().UTC(),
		ThresholdPercent: threshold,
		Packages:         parsePackageCoverage(goTestOut),
		Functions:        parseFunctionCoverage(funcOut),
		TotalPercent:     parseTotalCoverage(funcOut),
		RawGoTestOutput:  goTestOut,
	}
	below := report.TotalPercent < threshold

	if err := s.persistCoverageArtifacts(ctx, projectID, report); err != nil {
		return nil, below, err
	}
	if below && s.autoReject != nil {
		desc := []string{fmt.Sprintf("Coverage %.2f%% below threshold %.2f%%", report.TotalPercent, threshold)}
		_ = s.autoReject(ctx, projectID, ownerID, domain.CatastrophicFailureReport{
			CoveragePercent:     report.TotalPercent,
			SourcePhase:         6,
			FailureDescriptions: desc,
		})
	}
	return report, below, nil
}

func (s *Phase6Service) ValidateTests(ctx context.Context, backendDir, frontendDir string) (*domain.Phase6ValidationResult, error) {
	if backendDir == "" {
		backendDir = "."
	}
	result := &domain.Phase6ValidationResult{}
	goOut, err := s.runner.Run(ctx, backendDir, "go", "test", "./...", "-run", ".")
	if err == nil {
		result.GoTestPassed = true
	} else {
		result.FailureKind = classifyFailure(goOut)
		result.Details = strings.TrimSpace(goOut)
		return result, nil
	}

	raceOut, err := s.runner.Run(ctx, backendDir, "go", "test", "./...", "-race", "-run", ".")
	if err == nil {
		result.GoRacePassed = true
	} else {
		result.FailureKind = classifyFailure(raceOut)
		result.Details = strings.TrimSpace(raceOut)
		return result, nil
	}

	if frontendDir == "" {
		result.FrontendTestPassed = true
		result.FailureKind = domain.TestFailureNone
		return result, nil
	}
	frontOut, err := s.runner.Run(ctx, frontendDir, "npm", "test", "--", "--watchAll=false")
	if err == nil {
		result.FrontendTestPassed = true
		result.FailureKind = domain.TestFailureNone
		return result, nil
	}
	result.FailureKind = classifyFailure(frontOut)
	result.Details = strings.TrimSpace(frontOut)
	return result, nil
}

func (s *Phase6Service) BuildQualityReport(report *domain.Phase6CoverageReport, validation *domain.Phase6ValidationResult) string {
	var b strings.Builder
	b.WriteString("# QUALITY_REPORT\n\n")
	if report != nil {
		b.WriteString(fmt.Sprintf("- Cobertura total: **%.2f%%** (threshold: %.2f%%)\n", report.TotalPercent, report.ThresholdPercent))
		b.WriteString(fmt.Sprintf("- Pacotes analisados: **%d**\n", len(report.Packages)))
		b.WriteString(fmt.Sprintf("- Funções analisadas: **%d**\n", len(report.Functions)))
	}
	if validation != nil {
		b.WriteString("\n## Validação de Execução\n\n")
		b.WriteString(fmt.Sprintf("- go test: **%t**\n", validation.GoTestPassed))
		b.WriteString(fmt.Sprintf("- go test -race: **%t**\n", validation.GoRacePassed))
		b.WriteString(fmt.Sprintf("- frontend tests: **%t**\n", validation.FrontendTestPassed))
		b.WriteString(fmt.Sprintf("- classificação de falha: **%s**\n", validation.FailureKind))
	}
	if report != nil {
		b.WriteString("\n## Pacotes com menor cobertura\n\n")
		pkgs := append([]domain.PackageCoverage(nil), report.Packages...)
		sort.Slice(pkgs, func(i, j int) bool { return pkgs[i].Percent < pkgs[j].Percent })
		limit := 5
		if len(pkgs) < limit {
			limit = len(pkgs)
		}
		for i := 0; i < limit; i++ {
			b.WriteString(fmt.Sprintf("- `%s`: %.2f%%\n", pkgs[i].Package, pkgs[i].Percent))
		}
	}
	return b.String()
}

func (s *Phase6Service) PersistQualityReport(ctx context.Context, projectID, content string) error {
	projectOID, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return err
	}
	return s.files.Upsert(ctx, &domain.CodeFile{
		ProjectID:   projectOID,
		Path:        "reports/QUALITY_REPORT.md",
		TaskID:      "TASK-11-008",
		Language:    "markdown",
		PhaseNumber: 6,
		Content:     content,
	})
}

func (s *Phase6Service) persistCoverageArtifacts(ctx context.Context, projectID string, report *domain.Phase6CoverageReport) error {
	projectOID, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return err
	}
	coverageMD := s.BuildQualityReport(report, nil)
	if err := s.files.Upsert(ctx, &domain.CodeFile{
		ProjectID:   projectOID,
		Path:        "reports/phase6/coverage_report.md",
		TaskID:      "TASK-11-003",
		Language:    "markdown",
		PhaseNumber: 6,
		Content:     coverageMD,
	}); err != nil {
		return err
	}
	return nil
}

var rePkgCoverage = regexp.MustCompile(`coverage:\s*([0-9]+(?:\.[0-9]+)?)%`)
var rePercent = regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)%`)

func parsePackageCoverage(out string) []domain.PackageCoverage {
	lines := strings.Split(out, "\n")
	items := make([]domain.PackageCoverage, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "coverage:") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		match := rePkgCoverage.FindStringSubmatch(line)
		if len(match) != 2 {
			continue
		}
		pct, _ := strconv.ParseFloat(match[1], 64)
		pkg := fields[1]
		items = append(items, domain.PackageCoverage{Package: pkg, Percent: pct})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Package < items[j].Package })
	return items
}

func parseFunctionCoverage(out string) []domain.FunctionCoverage {
	lines := strings.Split(out, "\n")
	items := make([]domain.FunctionCoverage, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "total:") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		sourceFn := parts[0]
		pctMatch := rePercent.FindStringSubmatch(line)
		if len(pctMatch) != 2 {
			continue
		}
		pct, _ := strconv.ParseFloat(pctMatch[1], 64)
		name := sourceFn
		if idx := strings.LastIndex(sourceFn, "."); idx >= 0 && idx < len(sourceFn)-1 {
			name = sourceFn[idx+1:]
		}
		items = append(items, domain.FunctionCoverage{Name: name, Source: sourceFn, Percent: pct})
	}
	return items
}

func parseTotalCoverage(out string) float64 {
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "total:") {
			continue
		}
		match := rePercent.FindStringSubmatch(line)
		if len(match) == 2 {
			pct, _ := strconv.ParseFloat(match[1], 64)
			return pct
		}
	}
	return 0
}

func classifyFailure(out string) domain.TestFailureKind {
	l := strings.ToLower(out)
	if l == "" {
		return domain.TestFailureNone
	}
	if strings.Contains(l, "_test.go") && (strings.Contains(l, "undefined") || strings.Contains(l, "build failed") || strings.Contains(l, "compile")) {
		return domain.TestFailureTestImplementation
	}
	return domain.TestFailureProjectBug
}
