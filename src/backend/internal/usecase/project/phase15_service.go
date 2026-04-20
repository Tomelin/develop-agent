package project

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	domain "github.com/develop-agent/backend/internal/domain/project"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Phase15Service struct {
	projects   domain.ProjectRepository
	files      domain.CodeFileRepository
	httpClient *http.Client
}

func NewPhase15Service(projects domain.ProjectRepository, files domain.CodeFileRepository) *Phase15Service {
	return &Phase15Service{
		projects: projects,
		files:    files,
		httpClient: &http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

func (s *Phase15Service) Run(ctx context.Context, projectID, ownerID string, in domain.Phase15RunInput) (*domain.Phase15DeliveryReport, error) {
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

	channels := normalizeMarketingChannels(in.Channels)
	budget := in.MonthlyBudgetUSD
	if budget <= 0 {
		budget = 3000
	}

	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, err
	}

	artifacts := make([]domain.CodeFile, 0, 24)
	warnings := make([]string, 0, 2)
	briefJSON, _ := json.MarshalIndent(brief, "", "  ")

	artifacts = append(artifacts,
		mkPhase15File(pid, "docs/prompts/marketing/producer.md", "TASK-15-001", "markdown", 15, buildMarketingProducerPrompt(*brief, channels)),
		mkPhase15File(pid, "docs/prompts/marketing/reviewer.md", "TASK-15-001", "markdown", 15, buildMarketingReviewerPrompt(*brief, channels)),
		mkPhase15File(pid, "docs/prompts/marketing/refiner.md", "TASK-15-001", "markdown", 15, buildMarketingRefinerPrompt(*brief)),
		mkPhase15File(pid, "artifacts/marketing/MARKETING_BRIEF.json", "TASK-15-002", "json", 15, string(briefJSON)),
	)

	strategy := buildMarketingStrategy(*brief, channels, budget)
	calendarRows := buildCalendarRows(channels)
	calendarCSV := csvFromRows(calendarRows)
	calendarICS := buildCalendarICS(calendarRows)
	perf := buildPerformanceForecast(channels)

	artifacts = append(artifacts,
		mkPhase15File(pid, "artifacts/marketing/strategy/MARKETING_STRATEGY.md", "TASK-15-003", "markdown", 15, strategy),
		mkPhase15File(pid, "artifacts/marketing/strategy/CALENDAR.csv", "TASK-15-007", "csv", 15, calendarCSV),
		mkPhase15File(pid, "artifacts/marketing/strategy/CALENDAR.ics", "TASK-15-007", "ics", 15, calendarICS),
		mkPhase15File(pid, "artifacts/marketing/PERFORMANCE_FORECAST.md", "TASK-15-009", "markdown", 15, perf),
	)

	summaries := make([]domain.MarketingChannelSummary, 0, len(channels))
	totalPieces := 0
	for _, ch := range channels {
		piece, contentFiles := buildChannelPack(pid, ch, brief.PrimaryCTA)
		totalPieces += piece
		artifacts = append(artifacts, contentFiles...)
		summaries = append(summaries, buildChannelSummary(ch, piece, budget/float64(len(channels))))
	}

	for i := range artifacts {
		if err := s.files.Upsert(ctx, &artifacts[i]); err != nil {
			return nil, err
		}
	}

	if totalPieces == 0 {
		warnings = append(warnings, "no content pieces generated")
	}

	paths := make([]string, 0, len(artifacts))
	for _, a := range artifacts {
		paths = append(paths, a.Path)
	}
	sort.Strings(paths)

	p.UpdatedAt = time.Now().UTC()
	if err := s.projects.Update(ctx, p); err != nil {
		return nil, err
	}

	return &domain.Phase15DeliveryReport{
		GeneratedAt:      time.Now().UTC(),
		ProjectID:        p.ID.Hex(),
		BriefSource:      brief.Source,
		Channels:         channels,
		TotalPieces:      totalPieces,
		ArtifactPaths:    paths,
		ChannelSummaries: summaries,
		Warnings:         warnings,
	}, nil
}

func (s *Phase15Service) ExportPack(ctx context.Context, projectID, ownerID string, channels []string) ([]byte, string, int, error) {
	p, err := s.projects.FindByID(ctx, projectID)
	if err != nil {
		return nil, "", 0, err
	}
	if p.OwnerUserID.Hex() != ownerID {
		return nil, "", 0, fmt.Errorf("project not found")
	}

	files, err := s.files.ListByProject(ctx, projectID)
	if err != nil {
		return nil, "", 0, err
	}
	filteredChannels := normalizeMarketingChannels(channels)
	allowed := map[string]struct{}{}
	for _, ch := range filteredChannels {
		allowed[ch] = struct{}{}
	}

	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)
	pieces := 0
	for _, f := range files {
		if !strings.Contains(f.Path, "artifacts/marketing/") {
			continue
		}
		if len(allowed) > 0 && !shouldIncludePathForChannels(f.Path, allowed) {
			continue
		}
		w, err := zw.Create(filepath.ToSlash(strings.TrimPrefix(f.Path, "artifacts/marketing/")))
		if err != nil {
			return nil, "", 0, err
		}
		if _, err := io.WriteString(w, f.Content); err != nil {
			return nil, "", 0, err
		}
		if isPiecePath(f.Path) {
			pieces++
		}
	}
	if err := zw.Close(); err != nil {
		return nil, "", 0, err
	}

	name := fmt.Sprintf("marketing-pack-%s.zip", projectID)
	if len(filteredChannels) > 0 {
		name = fmt.Sprintf("marketing-pack-%s-%s.zip", projectID, strings.Join(filteredChannels, "-"))
	}
	return buf.Bytes(), name, pieces, nil
}

func (s *Phase15Service) ConfigureWebhook(ctx context.Context, projectID, ownerID string, in domain.MarketingWebhookInput) (*domain.MarketingWebhookResult, error) {
	p, err := s.projects.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if p.OwnerUserID.Hex() != ownerID {
		return nil, fmt.Errorf("project not found")
	}
	url := strings.TrimSpace(in.URL)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return nil, fmt.Errorf("invalid webhook URL")
	}

	delivery := domain.MarketingWebhookDelivery{Timestamp: time.Now().UTC(), Status: "failed"}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := s.httpClient.Do(req)
	if err == nil && resp != nil {
		delivery.ResponseStatus = resp.StatusCode
		_ = resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			delivery.Status = "success"
		} else {
			delivery.Error = "healthcheck did not return status 200"
		}
	} else if err != nil {
		delivery.Error = err.Error()
	}
	if delivery.Status != "success" {
		return nil, fmt.Errorf("webhook validation failed: %s", delivery.Error)
	}

	pid, _ := bson.ObjectIDFromHex(projectID)
	payload := map[string]any{
		"url":          url,
		"validated_at": time.Now().UTC(),
		"last_test":    delivery,
	}
	raw, _ := json.MarshalIndent(payload, "", "  ")
	if err := s.files.Upsert(ctx, &domain.CodeFile{ProjectID: pid, Path: "artifacts/marketing/webhooks/config.json", TaskID: "TASK-15-010", Language: "json", PhaseNumber: 15, Content: string(raw)}); err != nil {
		return nil, err
	}

	return &domain.MarketingWebhookResult{URL: url, ValidatedAt: time.Now().UTC(), LastTest: delivery}, nil
}

func (s *Phase15Service) buildBrief(ctx context.Context, p *domain.Project, in domain.Phase15RunInput) (*domain.MarketingBrief, error) {
	if in.UseLinkedProject {
		if p.LinkedProjectID == nil {
			return nil, fmt.Errorf("linked project not configured")
		}
		return s.extractBriefFromLinkedProject(ctx, p.LinkedProjectID.Hex(), in.ManualBrief)
	}
	brief, err := manualToMarketingBrief(in.ManualBrief)
	if err != nil {
		return nil, err
	}
	brief.Source = "manual"
	return brief, nil
}

func (s *Phase15Service) extractBriefFromLinkedProject(ctx context.Context, linkedProjectID string, fallback domain.MarketingManualBrief) (*domain.MarketingBrief, error) {
	files, err := s.files.ListByProject(ctx, linkedProjectID)
	if err != nil {
		return nil, err
	}
	vision := findFileContent(files, "VISION.md")
	spec := findFileContent(files, "SPEC.md")
	prompt := findFileContent(files, "PROMPT")

	brief := &domain.MarketingBrief{
		Source:            "linked_project",
		ProductName:       firstNonEmpty(extractField(vision, "nome do produto"), extractField(vision, "product name"), fallback.ProductName),
		Tagline:           firstNonEmpty(extractField(vision, "tagline"), fallback.Tagline),
		ProblemSolved:     firstNonEmpty(extractField(vision, "problema"), extractField(vision, "problem"), fallback.ProblemSolved),
		TargetAudience:    firstNonEmpty(extractField(vision, "público-alvo"), extractField(vision, "target audience"), fallback.TargetAudience),
		MainBenefits:      firstNonEmptyList(extractBulletList(spec, 5), fallback.MainBenefits),
		Differentials:     firstNonEmptyList(extractDifferentials(spec), fallback.Differentials),
		BusinessModel:     firstNonEmpty(extractField(spec, "modelo de negócio"), extractField(spec, "business model"), fallback.BusinessModel),
		Pricing:           firstNonEmpty(extractField(spec, "preço"), extractField(spec, "pricing"), fallback.Pricing),
		MarketType:        firstNonEmpty(extractField(spec, "mercado"), fallback.MarketType),
		CommunicationTone: firstNonEmpty(extractField(prompt, "tom"), fallback.CommunicationTone, "consultivo"),
		PrimaryCTA:        firstNonEmpty(fallback.PrimaryCTA, "Agendar estratégia"),
		SecondaryCTA:      firstNonEmpty(fallback.SecondaryCTA, "Receber diagnóstico"),
	}
	brief.CompetitorReferences = firstNonEmptyList(extractCompetitors(vision+"\n"+spec), fallback.CompetitorReferences)
	if err := validateMarketingBrief(*brief); err != nil {
		return nil, err
	}
	return brief, nil
}

func manualToMarketingBrief(in domain.MarketingManualBrief) (*domain.MarketingBrief, error) {
	brief := &domain.MarketingBrief{
		Source:               "manual",
		ProductName:          strings.TrimSpace(in.ProductName),
		Tagline:              strings.TrimSpace(in.Tagline),
		ProblemSolved:        strings.TrimSpace(in.ProblemSolved),
		TargetAudience:       strings.TrimSpace(in.TargetAudience),
		MainBenefits:         dedupeAndTrim(in.MainBenefits),
		Differentials:        dedupeAndTrim(in.Differentials),
		BusinessModel:        strings.TrimSpace(in.BusinessModel),
		Pricing:              strings.TrimSpace(in.Pricing),
		MarketType:           firstNonEmpty(in.MarketType, "B2B"),
		CommunicationTone:    firstNonEmpty(in.CommunicationTone, "consultivo"),
		PrimaryCTA:           firstNonEmpty(in.PrimaryCTA, "Solicitar plano"),
		SecondaryCTA:         firstNonEmpty(in.SecondaryCTA, "Falar com especialista"),
		CompetitorReferences: dedupeAndTrim(in.CompetitorReferences),
	}
	if err := validateMarketingBrief(*brief); err != nil {
		return nil, err
	}
	return brief, nil
}

func validateMarketingBrief(brief domain.MarketingBrief) error {
	if strings.TrimSpace(brief.ProductName) == "" {
		return fmt.Errorf("manual brief: product_name is required")
	}
	if strings.TrimSpace(brief.ProblemSolved) == "" {
		return fmt.Errorf("manual brief: problem_solved is required")
	}
	if strings.TrimSpace(brief.TargetAudience) == "" {
		return fmt.Errorf("manual brief: target_audience is required")
	}
	if len(brief.MainBenefits) == 0 {
		return fmt.Errorf("manual brief: at least one main_benefit is required")
	}
	return nil
}

func mkPhase15File(projectID bson.ObjectID, path, task, lang string, phase int, content string) domain.CodeFile {
	return domain.CodeFile{ProjectID: projectID, Path: filepath.ToSlash(path), TaskID: task, Language: lang, PhaseNumber: phase, Content: content}
}

func normalizeMarketingChannels(in []string) []string {
	if len(in) == 0 {
		return []string{"linkedin", "instagram", "google-ads"}
	}
	allowed := map[string]struct{}{"linkedin": {}, "instagram": {}, "google-ads": {}}
	out := make([]string, 0, 3)
	seen := map[string]struct{}{}
	for _, raw := range in {
		k := strings.ToLower(strings.TrimSpace(raw))
		if _, ok := allowed[k]; !ok {
			continue
		}
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, k)
	}
	if len(out) == 0 {
		return []string{"linkedin", "instagram", "google-ads"}
	}
	sort.Strings(out)
	return out
}

func buildMarketingProducerPrompt(brief domain.MarketingBrief, channels []string) string {
	return fmt.Sprintf(`# Produtor — Marketing Strategist

Crie uma estratégia de marketing multi-canal para %s.

## Regras obrigatórias
- Defina objetivos SMART.
- Estruture buyer personas com dores, desejos e objeções.
- Justifique os canais prioritários: %s.
- Inclua plano de conteúdo, budget por canal e KPIs.
- Gere conteúdo adaptado por canal (tom, formato, CTA).
`, brief.ProductName, strings.Join(channels, ", "))
}

func buildMarketingReviewerPrompt(brief domain.MarketingBrief, channels []string) string {
	return fmt.Sprintf(`# Revisor — CRO & Marketing Specialist

Revise a estratégia de %s com postura crítica.

Checklist:
1. Alinhamento da mensagem com o público (%s)
2. Coerência de funil (topo, meio, fundo)
3. Compliance para %s
4. Realismo de estimativas de ROI
5. Gaps ou riscos operacionais
`, brief.ProductName, brief.TargetAudience, strings.Join(channels, ", "))
}

func buildMarketingRefinerPrompt(brief domain.MarketingBrief) string {
	return fmt.Sprintf(`# Refinador — Campaign Delivery

Otimize o material final de %s garantindo:
- Clareza na proposta de valor
- Consistência de tom (%s)
- Calls-to-action acionáveis
- Entrega final sem ambiguidades
`, brief.ProductName, brief.CommunicationTone)
}

func buildMarketingStrategy(brief domain.MarketingBrief, channels []string, budget float64) string {
	channelsSection := make([]string, 0, len(channels))
	for _, ch := range channels {
		channelsSection = append(channelsSection, fmt.Sprintf("### %s\n- Público principal: %s\n- Frequência: 3 a 5 publicações/semana\n- CTA foco: %s\n", strings.ToUpper(ch), brief.TargetAudience, brief.PrimaryCTA))
	}

	return fmt.Sprintf(`# MARKETING_STRATEGY

## Executive Summary
- Produto: **%s**
- Objetivo principal: gerar pipeline qualificado em 90 dias
- Budget mensal sugerido: **US$ %.2f**
- Canais prioritários: %s

## Análise de Audiência
- Persona 1: decisor com dor em %s
- Persona 2: executor buscando previsibilidade de demanda
- Comportamento online: consumo de conteúdo educacional + comparativos

## Estratégia por Canal
%s
## Messaging Framework
- Mensagem central: "%s"
- Tom de voz: %s
- Variações: educação (TOFU), prova social (MOFU), oferta (BOFU)

## Calendário de Conteúdo
Ver `+"`CALENDAR.csv`"+` e `+"`CALENDAR.ics`"+`.

## Budget Allocation
- LinkedIn: 40%%
- Instagram: 25%%
- Google Ads: 35%%

## KPIs por Canal
- LinkedIn: CTR 0.40%%, CPL até US$ 45
- Instagram: CTR 0.90%%, CPM US$ 7-12
- Google Ads: CTR 3.5%%, CAC até US$ 85

## Plano de 90 dias
1. **Dias 1-30 (Awareness):** validação de mensagem e criativos
2. **Dias 31-60 (Consideração):** retargeting + provas de valor
3. **Dias 61-90 (Conversão):** campanhas com ofertas e SQL handoff
`, brief.ProductName, budget, strings.Join(channels, ", "), brief.ProblemSolved, strings.Join(channelsSection, "\n"), firstNonEmpty(brief.Tagline, brief.ProblemSolved), brief.CommunicationTone)
}

func buildChannelPack(projectID bson.ObjectID, channel, cta string) (int, []domain.CodeFile) {
	switch channel {
	case "linkedin":
		return 18, []domain.CodeFile{
			mkPhase15File(projectID, "artifacts/marketing/linkedin/organic/posts.md", "TASK-15-004", "markdown", 15, linkedinOrganic(cta)),
			mkPhase15File(projectID, "artifacts/marketing/linkedin/articles/articles.md", "TASK-15-004", "markdown", 15, linkedinArticles(cta)),
			mkPhase15File(projectID, "artifacts/marketing/linkedin/ads/campaigns.md", "TASK-15-004", "markdown", 15, linkedinAds(cta)),
		}
	case "instagram":
		return 29, []domain.CodeFile{
			mkPhase15File(projectID, "artifacts/marketing/instagram/feed/posts.md", "TASK-15-005", "markdown", 15, instagramFeed(cta)),
			mkPhase15File(projectID, "artifacts/marketing/instagram/stories/sequences.md", "TASK-15-005", "markdown", 15, instagramStories(cta)),
			mkPhase15File(projectID, "artifacts/marketing/instagram/reels/scripts.md", "TASK-15-005", "markdown", 15, instagramReels(cta)),
			mkPhase15File(projectID, "artifacts/marketing/instagram/ads/creatives.md", "TASK-15-005", "markdown", 15, instagramAds(cta)),
		}
	case "google-ads":
		return 12, []domain.CodeFile{
			mkPhase15File(projectID, "artifacts/marketing/google-ads/keywords.csv", "TASK-15-006", "csv", 15, googleAdsKeywordsCSV()),
			mkPhase15File(projectID, "artifacts/marketing/google-ads/ads.md", "TASK-15-006", "markdown", 15, googleAdsAds(cta)),
		}
	default:
		return 0, nil
	}
}

func buildChannelSummary(channel string, pieces int, budget float64) domain.MarketingChannelSummary {
	summary := domain.MarketingChannelSummary{Channel: channel, Pieces: pieces, BudgetUSD: budget}
	switch channel {
	case "linkedin":
		summary.ExpectedCTR = "0.3% - 0.5%"
		summary.ExpectedConv = "1.0% - 2.0%"
	case "instagram":
		summary.ExpectedCTR = "0.8% - 1.2%"
		summary.ExpectedConv = "0.7% - 1.5%"
	case "google-ads":
		summary.ExpectedCTR = "3.0% - 5.0%"
		summary.ExpectedConv = "2.5% - 6.0%"
	}
	return summary
}

func buildCalendarRows(channels []string) [][]string {
	rows := [][]string{{"date", "channel", "asset_type", "title", "best_time_utc", "cta"}}
	start := time.Now().UTC()
	for i := 0; i < 30; i++ {
		date := start.AddDate(0, 0, i).Format("2006-01-02")
		for _, ch := range channels {
			rows = append(rows, []string{date, ch, "post", fmt.Sprintf("%s content #%d", strings.ToUpper(ch), i+1), "14:00", "Gerar demanda"})
		}
	}
	return rows
}

func csvFromRows(rows [][]string) string {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	_ = w.WriteAll(rows)
	w.Flush()
	return buf.String()
}

func buildCalendarICS(rows [][]string) string {
	b := &strings.Builder{}
	b.WriteString("BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//Develop Agent//Marketing Calendar//EN\n")
	for i, row := range rows {
		if i == 0 || len(row) < 6 {
			continue
		}
		start, _ := time.Parse("2006-01-02", row[0])
		dt := start.Format("20060102")
		b.WriteString("BEGIN:VEVENT\n")
		b.WriteString("UID:" + fmt.Sprintf("%d-%s@develop-agent", i, row[1]) + "\n")
		b.WriteString("DTSTAMP:" + time.Now().UTC().Format("20060102T150405Z") + "\n")
		b.WriteString("DTSTART;VALUE=DATE:" + dt + "\n")
		b.WriteString("SUMMARY:" + row[3] + "\n")
		b.WriteString("DESCRIPTION:Canal " + row[1] + " | CTA " + row[5] + "\n")
		b.WriteString("END:VEVENT\n")
	}
	b.WriteString("END:VCALENDAR\n")
	return b.String()
}

func buildPerformanceForecast(channels []string) string {
	lines := []string{"# PERFORMANCE_FORECAST", "", "Estimativas baseadas em benchmarks setoriais (não garantias).", ""}
	for _, ch := range channels {
		switch ch {
		case "linkedin":
			lines = append(lines, "## LinkedIn", "- CTR: 0.3% (mín) | 0.4% (médio) | 0.6% (otimista)", "- CPL: US$ 32 | US$ 24 | US$ 18", "- Alcance orgânico por post: 400 | 900 | 1500")
		case "instagram":
			lines = append(lines, "## Instagram", "- CPM: US$ 12 | US$ 9 | US$ 7", "- CTR: 0.8% | 1.0% | 1.3%", "- Custo por seguidor: US$ 2.2 | US$ 1.5 | US$ 1.0")
		case "google-ads":
			lines = append(lines, "## Google Ads", "- CPC médio: US$ 1.8 | US$ 1.3 | US$ 0.9", "- CTR Search: 3.0% | 3.8% | 5.2%", "- Conversão: 2.5% | 3.8% | 6.0%")
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

func shouldIncludePathForChannels(path string, allowed map[string]struct{}) bool {
	if strings.Contains(path, "/strategy/") || strings.Contains(path, "PERFORMANCE_FORECAST") || strings.Contains(path, "/prompts/") {
		return true
	}
	for ch := range allowed {
		if strings.Contains(path, "/"+ch+"/") {
			return true
		}
	}
	return false
}

func isPiecePath(path string) bool {
	return strings.HasSuffix(path, ".md") && (strings.Contains(path, "/linkedin/") || strings.Contains(path, "/instagram/") || strings.Contains(path, "/google-ads/"))
}

func extractDifferentials(content string) []string {
	return extractBulletList(content, 5)
}

func extractCompetitors(content string) []string {
	re := regexp.MustCompile(`(?im)(hubspot|rd station|salesforce|pipedrive|mailchimp)`)
	return dedupeAndTrim(re.FindAllString(content, -1))
}

func linkedinOrganic(cta string) string {
	return "# LinkedIn Organic\n\n10 posts com hooks, prova social e CTA: " + cta + "\n"
}
func linkedinArticles(cta string) string {
	return "# LinkedIn Articles\n\n3 artigos longos voltados ao ICP com CTA final: " + cta + "\n"
}
func linkedinAds(cta string) string {
	return "# LinkedIn Ads\n\n5 variações de Sponsored Content + 3 Message Ads + targeting sugerido. CTA: " + cta + "\n"
}
func instagramFeed(cta string) string {
	return "# Instagram Feed\n\n12 legendas com hashtags de nicho e trending. CTA: " + cta + "\n"
}
func instagramStories(cta string) string {
	return "# Instagram Stories\n\n8 sequências de stories (3-5 slides) com enquete/FAQ/prova social.\n"
}
func instagramReels(cta string) string {
	return "# Instagram Reels\n\n5 scripts (30-60s) com hook inicial e CTA: " + cta + "\n"
}
func instagramAds(cta string) string {
	return "# Instagram Ads\n\n4 variações por objetivo (awareness, consideration, conversion). CTA: " + cta + "\n"
}
func googleAdsKeywordsCSV() string {
	return "intent,keyword,ad_group\ninformacional,como reduzir CAC,saaS awareness\ncomercial,software automacao marketing,lead gen\ntransacional,plataforma de marketing b2b,bofu\n"
}
func googleAdsAds(cta string) string {
	return "# Google Ads\n\nSearch + Display + extensões + negativas + budget escalonado. CTA: " + cta + "\n"
}
