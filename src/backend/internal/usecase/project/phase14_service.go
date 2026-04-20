package project

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	domain "github.com/develop-agent/backend/internal/domain/project"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Phase14Service struct {
	projects domain.ProjectRepository
	files    domain.CodeFileRepository
}

func NewPhase14Service(projects domain.ProjectRepository, files domain.CodeFileRepository) *Phase14Service {
	return &Phase14Service{projects: projects, files: files}
}

func (s *Phase14Service) Run(ctx context.Context, projectID, ownerID string, in domain.Phase14RunInput) (*domain.Phase14DeliveryReport, error) {
	p, err := s.projects.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if p.OwnerUserID.Hex() != ownerID {
		return nil, fmt.Errorf("project not found")
	}

	brief, err := s.buildBrief(ctx, p, in)
	if err != nil {
		return nil, err
	}

	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, err
	}

	producerPrompt := buildLandingProducerPrompt(*brief)
	reviewerPrompt := buildLandingReviewerPrompt(*brief)
	refinerPrompt := buildLandingRefinerPrompt(*brief)
	landingHTML := buildLandingHTML(*brief, "benefit-focused")
	conversion := scoreLandingConversion(landingHTML)
	convReport := buildConversionReport(*brief, conversion)
	seoChecklist, seoFindings := buildSEOChecklist(*brief, landingHTML)

	briefJSON, _ := json.MarshalIndent(brief, "", "  ")
	artifacts := []domain.CodeFile{
		mkPhase14File(pid, "docs/prompts/landing/producer.md", "TASK-14-001", "markdown", 14, producerPrompt),
		mkPhase14File(pid, "docs/prompts/landing/reviewer.md", "TASK-14-001", "markdown", 14, reviewerPrompt),
		mkPhase14File(pid, "docs/prompts/landing/refiner.md", "TASK-14-001", "markdown", 14, refinerPrompt),
		mkPhase14File(pid, "artifacts/landing/brief.json", "TASK-14-002", "json", 14, string(briefJSON)),
		mkPhase14File(pid, "artifacts/landing/landing_page.html", "TASK-14-004", "html", 14, landingHTML),
		mkPhase14File(pid, "artifacts/landing/CONVERSION_REPORT.md", "TASK-14-006", "markdown", 14, convReport),
		mkPhase14File(pid, "artifacts/landing/SEO_CHECKLIST.md", "TASK-14-008", "markdown", 14, seoChecklist),
	}

	variantReports := make([]domain.LandingPageVariantReport, 0, 3)
	if in.GenerateVariants {
		count := in.VariantCount
		if count <= 0 {
			count = 3
		}
		if count > 3 {
			count = 3
		}
		styles := []string{"benefit-focused", "problem-focused", "curiosity-driven"}
		for i := 0; i < count; i++ {
			name := fmt.Sprintf("variant_%d", i+1)
			html := buildLandingHTML(*brief, styles[i])
			path := filepath.ToSlash(fmt.Sprintf("artifacts/landing/variants/%s.html", name))
			artifacts = append(artifacts, mkPhase14File(pid, path, "TASK-14-007", "html", 14, html))
			variantReports = append(variantReports, domain.LandingPageVariantReport{
				Name:            name,
				Path:            path,
				ConversionScore: scoreLandingConversion(html),
			})
		}
	}

	artifactPaths := make([]string, 0, len(artifacts))
	for i := range artifacts {
		if err := s.files.Upsert(ctx, &artifacts[i]); err != nil {
			return nil, err
		}
		artifactPaths = append(artifactPaths, artifacts[i].Path)
	}
	sort.Strings(artifactPaths)

	p.UpdatedAt = time.Now().UTC()
	if err := s.projects.Update(ctx, p); err != nil {
		return nil, err
	}

	return &domain.Phase14DeliveryReport{
		GeneratedAt:       time.Now().UTC(),
		ProjectID:         p.ID.Hex(),
		BriefSource:       brief.Source,
		OutputFormat:      brief.OutputFormat,
		ConversionScore:   conversion,
		ArtifactPaths:     artifactPaths,
		Variants:          variantReports,
		PrioritizedIssues: seoFindings,
	}, nil
}

func (s *Phase14Service) buildBrief(ctx context.Context, p *domain.Project, in domain.Phase14RunInput) (*domain.LandingPageBrief, error) {
	if in.UseLinkedProject {
		if p.LinkedProjectID == nil {
			return nil, fmt.Errorf("linked project not configured")
		}
		return s.extractBriefFromLinkedProject(ctx, p.LinkedProjectID.Hex(), in.ManualBrief)
	}
	brief, err := manualBriefToLanding(in.ManualBrief)
	if err != nil {
		return nil, err
	}
	brief.Source = "manual"
	return brief, nil
}

func (s *Phase14Service) extractBriefFromLinkedProject(ctx context.Context, linkedProjectID string, fallback domain.LandingPageManualBrief) (*domain.LandingPageBrief, error) {
	files, err := s.files.ListByProject(ctx, linkedProjectID)
	if err != nil {
		return nil, err
	}
	vision := findFileContent(files, "VISION.md")
	spec := findFileContent(files, "SPEC.md")
	userPrompt := findFileContent(files, "PROMPT")

	brief := &domain.LandingPageBrief{Source: "linked_project", OutputFormat: normalizeOutputFormat(fallback.OutputFormat)}
	brief.ProductName = firstNonEmpty(
		extractField(vision, "nome do produto"),
		extractField(vision, "product name"),
		fallback.ProductName,
	)
	brief.Tagline = firstNonEmpty(extractField(vision, "tagline"), extractField(vision, "slogan"))
	brief.ProblemSolved = firstNonEmpty(extractField(vision, "problema"), extractField(vision, "problem"), fallback.ProblemSolved)
	brief.TargetAudience = firstNonEmpty(extractField(vision, "público-alvo"), extractField(vision, "target audience"), fallback.TargetAudience)
	brief.UniqueValueProposed = firstNonEmpty(extractField(vision, "proposta de valor"), extractField(vision, "unique value"), fallback.UniqueValueProposed)
	brief.KeyFeatures = extractBulletList(spec, 5)
	if len(brief.KeyFeatures) == 0 {
		brief.KeyFeatures = dedupeAndTrim(fallback.KeyFeatures)
	}
	brief.BusinessModel = firstNonEmpty(extractField(spec, "modelo de negócio"), extractField(spec, "business model"))
	brief.RelevantIntegrations = extractIntegrations(spec)
	brief.ColorPalette = firstNonEmptyList(extractPalette(userPrompt), fallback.ColorPalette)
	brief.Theme = firstNonEmpty(extractField(userPrompt, "tema"), extractField(userPrompt, "theme"), fallback.Theme)
	brief.CommunicationTone = firstNonEmpty(extractField(userPrompt, "tom"), extractField(userPrompt, "tone"), fallback.CommunicationTone)
	brief.Language = firstNonEmpty(fallback.Language, "pt-BR")
	brief.PreferredTypography = firstNonEmpty(extractField(userPrompt, "tipografia"), fallback.PreferredTypography, "Inter")
	brief.PrimaryKeyword = firstNonEmpty(fallback.PrimaryKeyword, slugKeyword(brief.ProductName))
	brief.PrimaryCTA = firstNonEmpty(fallback.PrimaryCTA, "Começar agora")
	brief.SecondaryCTA = firstNonEmpty(fallback.SecondaryCTA, "Falar com especialista")
	brief.SocialProofHighlight = firstNonEmpty(fallback.SocialProofHighlight, "+1.000 usuários ativos")

	if err := validateLandingBrief(*brief); err != nil {
		return nil, err
	}
	return brief, nil
}

func manualBriefToLanding(in domain.LandingPageManualBrief) (*domain.LandingPageBrief, error) {
	brief := &domain.LandingPageBrief{
		Source:               "manual",
		ProductName:          strings.TrimSpace(in.ProductName),
		ProblemSolved:        strings.TrimSpace(in.ProblemSolved),
		TargetAudience:       strings.TrimSpace(in.TargetAudience),
		UniqueValueProposed:  strings.TrimSpace(in.UniqueValueProposed),
		KeyFeatures:          dedupeAndTrim(in.KeyFeatures),
		ColorPalette:         dedupeAndTrim(in.ColorPalette),
		Theme:                firstNonEmpty(in.Theme, "light"),
		CommunicationTone:    firstNonEmpty(in.CommunicationTone, "profissional"),
		Language:             firstNonEmpty(in.Language, "pt-BR"),
		PreferredTypography:  firstNonEmpty(in.PreferredTypography, "Inter"),
		OutputFormat:         normalizeOutputFormat(in.OutputFormat),
		PrimaryKeyword:       firstNonEmpty(in.PrimaryKeyword, slugKeyword(in.ProductName)),
		PrimaryCTA:           firstNonEmpty(in.PrimaryCTA, "Começar agora"),
		SecondaryCTA:         firstNonEmpty(in.SecondaryCTA, "Agendar demo"),
		SocialProofHighlight: firstNonEmpty(in.SocialProofHighlight, "4.9/5 de satisfação"),
	}
	if err := validateLandingBrief(*brief); err != nil {
		return nil, err
	}
	return brief, nil
}

func validateLandingBrief(brief domain.LandingPageBrief) error {
	if strings.TrimSpace(brief.ProductName) == "" {
		return fmt.Errorf("manual brief: product_name is required")
	}
	if strings.TrimSpace(brief.ProblemSolved) == "" {
		return fmt.Errorf("manual brief: problem_solved is required")
	}
	if strings.TrimSpace(brief.TargetAudience) == "" {
		return fmt.Errorf("manual brief: target_audience is required")
	}
	if strings.TrimSpace(brief.UniqueValueProposed) == "" {
		return fmt.Errorf("manual brief: unique_value_proposed is required")
	}
	if len(brief.KeyFeatures) == 0 {
		return fmt.Errorf("manual brief: at least one key_feature is required")
	}
	if len(brief.KeyFeatures) > 5 {
		brief.KeyFeatures = brief.KeyFeatures[:5]
	}
	return nil
}

func findFileContent(files []*domain.CodeFile, needle string) string {
	needle = strings.ToLower(strings.TrimSpace(needle))
	for _, f := range files {
		if strings.Contains(strings.ToLower(f.Path), needle) {
			return f.Content
		}
	}
	return ""
}

func extractField(content, key string) string {
	if strings.TrimSpace(content) == "" {
		return ""
	}
	re := regexp.MustCompile(`(?im)^[-*\s]*` + regexp.QuoteMeta(key) + `\s*[:\-]\s*(.+)$`)
	if m := re.FindStringSubmatch(content); len(m) == 2 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func extractBulletList(content string, max int) []string {
	if strings.TrimSpace(content) == "" {
		return nil
	}
	re := regexp.MustCompile(`(?m)^\s*[-*]\s+(.+)$`)
	matches := re.FindAllStringSubmatch(content, max)
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) == 2 {
			out = append(out, strings.TrimSpace(m[1]))
		}
	}
	return dedupeAndTrim(out)
}

func extractIntegrations(content string) []string {
	re := regexp.MustCompile(`(?im)(slack|stripe|zapier|google|hubspot|salesforce|whatsapp|notion)`)
	matches := re.FindAllString(content, -1)
	return dedupeAndTrim(matches)
}

func extractPalette(content string) []string {
	re := regexp.MustCompile(`#[0-9a-fA-F]{6}`)
	return dedupeAndTrim(re.FindAllString(content, -1))
}

func dedupeAndTrim(items []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(items))
	for _, item := range items {
		v := strings.TrimSpace(item)
		if v == "" {
			continue
		}
		key := strings.ToLower(v)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, v)
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func firstNonEmptyList(values ...[]string) []string {
	for _, v := range values {
		if len(v) > 0 {
			return dedupeAndTrim(v)
		}
	}
	return nil
}

func normalizeOutputFormat(in string) string {
	v := strings.ToLower(strings.TrimSpace(in))
	switch v {
	case "next", "nextjs", "next.js":
		return "nextjs"
	default:
		return "html"
	}
}

func slugKeyword(name string) string {
	kw := strings.ToLower(strings.TrimSpace(name))
	if kw == "" {
		return "landing page"
	}
	kw = strings.ReplaceAll(kw, "_", " ")
	kw = strings.ReplaceAll(kw, "-", " ")
	return kw
}

func mkPhase14File(projectID bson.ObjectID, path, task, lang string, phase int, content string) domain.CodeFile {
	return domain.CodeFile{ProjectID: projectID, Path: filepath.ToSlash(path), TaskID: task, Language: lang, PhaseNumber: phase, Content: content}
}

func buildLandingProducerPrompt(brief domain.LandingPageBrief) string {
	return fmt.Sprintf(`# LANDING PAGE PRODUCER — HIGH CONVERSION

Você é o agente Produtor da Tríade para Fluxo B. Gere landing page completa e funcional.

## Brief do Projeto
- Produto: %s
- Problema: %s
- Público-alvo: %s
- UVP: %s
- Features: %s
- CTA principal: %s
- Formato de saída: %s

## Estrutura obrigatória
1. Hero (headline benefício + subheadline + CTA acima da dobra)
2. Social Proof (depoimentos, logos, métricas)
3. Features/Benefits
4. Como Funciona (3 passos)
5. FAQ
6. CTA final

## Regras
- Copy clara, voz ativa, orientada ao benefício.
- Mobile-first, semântica HTML correta.
- Incluir meta tags básicas SEO e alt text em imagens.
- Código de produção, sem placeholders vazios.
`, brief.ProductName, brief.ProblemSolved, brief.TargetAudience, brief.UniqueValueProposed, strings.Join(brief.KeyFeatures, "; "), brief.PrimaryCTA, brief.OutputFormat)
}

func buildLandingReviewerPrompt(brief domain.LandingPageBrief) string {
	return fmt.Sprintf(`# LANDING PAGE REVIEWER — CRO SPECIALIST

Analise a landing page para %s.

Checklist obrigatório:
- CRO: headline clara, CTA visível, prova social, baixa fricção.
- SEO: H1 único, meta description, alt texts.
- Acessibilidade: WCAG 2.1 básico (contraste, labels, foco).
- Performance: imagens otimizadas e CSS enxuto.

Devolva melhorias priorizadas (P0/P1/P2) com justificativa objetiva.
`, brief.ProductName)
}

func buildLandingRefinerPrompt(brief domain.LandingPageBrief) string {
	return fmt.Sprintf(`# LANDING PAGE REFINER — DELIVERY

Aplique 100%% das melhorias do revisor mantendo o posicionamento de %s.

Condições de entrega:
- HTML/CSS/JS válido de produção.
- SEO e acessibilidade corrigidos.
- Responsividade validada para 375/768/1440.
- Não remover seções obrigatórias do funnel.
`, brief.ProductName)
}

func buildLandingHTML(brief domain.LandingPageBrief, variant string) string {
	headline := fmt.Sprintf("%s sem complicação para %s", brief.ProductName, brief.TargetAudience)
	sub := brief.UniqueValueProposed
	cta := brief.PrimaryCTA
	if variant == "problem-focused" {
		headline = fmt.Sprintf("Pare de sofrer com %s", strings.ToLower(brief.ProblemSolved))
	}
	if variant == "curiosity-driven" {
		headline = fmt.Sprintf("O método usado por times de alta performance com %s", brief.ProductName)
		cta = "Quero ver como funciona"
	}
	features := strings.Join(brief.KeyFeatures, "</li><li>")

	return fmt.Sprintf(`<!doctype html>
<html lang="pt-BR">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>%s | Landing Page</title>
  <meta name="description" content="%s" />
  <meta property="og:title" content="%s" />
  <meta property="og:description" content="%s" />
  <meta property="og:image" content="https://example.com/og-image.png" />
  <meta name="twitter:card" content="summary_large_image" />
  <link rel="canonical" href="https://example.com" />
  <style>body{font-family:%s;margin:0;padding:0;line-height:1.5}.wrap{max-width:960px;margin:0 auto;padding:24px}button{background:#111;color:#fff;border:0;padding:12px 18px;border-radius:8px}</style>
</head>
<body>
  <main class="wrap">
    <section id="hero">
      <h1>%s</h1>
      <p>%s</p>
      <button>%s</button>
    </section>
    <section id="social-proof">
      <h2>Resultados comprovados</h2>
      <p>%s</p>
      <img src="https://example.com/logo.png" alt="Logos de clientes" />
    </section>
    <section id="features">
      <h2>Benefícios principais</h2>
      <ul><li>%s</li></ul>
    </section>
    <section id="how-it-works">
      <h2>Como funciona</h2>
      <ol><li>Cadastre-se</li><li>Configure em minutos</li><li>Veja resultados</li></ol>
    </section>
    <section id="faq">
      <h2>Perguntas frequentes</h2>
      <details><summary>Em quanto tempo vejo resultado?</summary><p>Nas primeiras semanas.</p></details>
    </section>
    <section id="cta-final">
      <h2>Pronto para começar?</h2>
      <button>%s</button>
    </section>
  </main>
  <script type="application/ld+json">{"@context":"https://schema.org","@type":"SoftwareApplication","name":"%s"}</script>
</body>
</html>
`, brief.ProductName, brief.UniqueValueProposed, brief.ProductName, brief.UniqueValueProposed, brief.PreferredTypography, headline, sub, cta, brief.SocialProofHighlight, features, brief.SecondaryCTA, brief.ProductName)
}

func scoreLandingConversion(html string) float64 {
	checks := []struct {
		pts int
		ok  bool
	}{
		{15, strings.Contains(html, "<h1>")},
		{10, strings.Contains(strings.ToLower(html), "meta name=\"description\"")},
		{15, strings.Contains(strings.ToLower(html), "<button")},
		{10, strings.Contains(strings.ToLower(html), "social-proof")},
		{15, strings.Contains(strings.ToLower(html), "benefícios") || strings.Contains(strings.ToLower(html), "beneficios")},
		{10, strings.Count(strings.ToLower(html), "button") <= 6},
		{10, strings.Contains(strings.ToLower(html), "viewport")},
		{10, strings.Contains(strings.ToLower(html), "<style>")},
		{5, strings.Contains(strings.ToLower(html), "canonical")},
	}
	total := 0
	for _, c := range checks {
		if c.ok {
			total += c.pts
		}
	}
	if total > 100 {
		return 100
	}
	return float64(total)
}

func buildConversionReport(brief domain.LandingPageBrief, score float64) string {
	return fmt.Sprintf(`# CONVERSION_REPORT

- Produto: **%s**
- Score final: **%.0f/100**

## Critérios avaliados
- Headline clara e benefício explícito
- CTA principal acima da dobra
- Prova social e redução de fricção
- Responsividade mobile
- SEO on-page básico

## Recomendação
Priorizar teste A/B do Hero (headline + CTA) para elevar conversão em tráfego frio.
`, brief.ProductName, score)
}

func buildSEOChecklist(brief domain.LandingPageBrief, html string) (string, []string) {
	checks := []struct {
		name string
		ok   bool
	}{
		{"H1 único", strings.Count(strings.ToLower(html), "<h1") == 1},
		{"Meta title 50-60", len(brief.ProductName) >= 10},
		{"Meta description 150-160", strings.Contains(strings.ToLower(html), "meta name=\"description\"")},
		{"Open Graph tags", strings.Contains(strings.ToLower(html), "property=\"og:title\"")},
		{"Twitter card", strings.Contains(strings.ToLower(html), "twitter:card")},
		{"Schema.org", strings.Contains(strings.ToLower(html), "application/ld+json")},
		{"Alt text imagens", strings.Contains(strings.ToLower(html), "alt=")},
		{"Canonical", strings.Contains(strings.ToLower(html), "rel=\"canonical\"")},
	}
	findings := make([]string, 0)
	var b strings.Builder
	b.WriteString("# SEO_CHECKLIST\n\n")
	b.WriteString("| Item | Status |\n|---|---|\n")
	for _, c := range checks {
		status := "OK"
		if !c.ok {
			status = "MISSING"
			findings = append(findings, c.name)
		}
		b.WriteString(fmt.Sprintf("| %s | %s |\n", c.name, status))
	}
	if len(findings) == 0 {
		findings = append(findings, "Nenhuma pendência crítica de SEO")
	}
	return b.String(), findings
}
