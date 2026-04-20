package project

import "time"

type MarketingManualBrief struct {
	ProductName          string   `json:"product_name"`
	Tagline              string   `json:"tagline,omitempty"`
	ProblemSolved        string   `json:"problem_solved"`
	TargetAudience       string   `json:"target_audience"`
	MainBenefits         []string `json:"main_benefits"`
	Differentials        []string `json:"differentials,omitempty"`
	BusinessModel        string   `json:"business_model,omitempty"`
	Pricing              string   `json:"pricing,omitempty"`
	MarketType           string   `json:"market_type,omitempty"`
	CommunicationTone    string   `json:"communication_tone,omitempty"`
	PrimaryCTA           string   `json:"primary_cta,omitempty"`
	SecondaryCTA         string   `json:"secondary_cta,omitempty"`
	CompetitorReferences []string `json:"competitor_references,omitempty"`
}

type Phase15RunInput struct {
	UseLinkedProject bool                 `json:"use_linked_project"`
	ManualBrief      MarketingManualBrief `json:"manual_brief"`
	Channels         []string             `json:"channels,omitempty"`
	MonthlyBudgetUSD float64              `json:"monthly_budget_usd,omitempty"`
}

type MarketingBrief struct {
	Source               string   `json:"source"`
	ProductName          string   `json:"product_name"`
	Tagline              string   `json:"tagline,omitempty"`
	ProblemSolved        string   `json:"problem_solved"`
	TargetAudience       string   `json:"target_audience"`
	MainBenefits         []string `json:"main_benefits"`
	Differentials        []string `json:"differentials,omitempty"`
	BusinessModel        string   `json:"business_model,omitempty"`
	Pricing              string   `json:"pricing,omitempty"`
	MarketType           string   `json:"market_type,omitempty"`
	CommunicationTone    string   `json:"communication_tone"`
	PrimaryCTA           string   `json:"primary_cta"`
	SecondaryCTA         string   `json:"secondary_cta,omitempty"`
	CompetitorReferences []string `json:"competitor_references,omitempty"`
}

type MarketingChannelSummary struct {
	Channel      string  `json:"channel"`
	Pieces       int     `json:"pieces"`
	ExpectedCTR  string  `json:"expected_ctr"`
	ExpectedConv string  `json:"expected_conversion"`
	BudgetUSD    float64 `json:"budget_usd"`
}

type Phase15DeliveryReport struct {
	GeneratedAt      time.Time                 `json:"generated_at"`
	ProjectID        string                    `json:"project_id"`
	BriefSource      string                    `json:"brief_source"`
	Channels         []string                  `json:"channels"`
	TotalPieces      int                       `json:"total_pieces"`
	ArtifactPaths    []string                  `json:"artifact_paths"`
	ChannelSummaries []MarketingChannelSummary `json:"channel_summaries"`
	Warnings         []string                  `json:"warnings,omitempty"`
}

type MarketingWebhookInput struct {
	URL string `json:"url"`
}

type MarketingWebhookDelivery struct {
	Timestamp      time.Time `json:"timestamp"`
	Status         string    `json:"status"`
	ResponseStatus int       `json:"response_status"`
	Error          string    `json:"error,omitempty"`
}

type MarketingWebhookResult struct {
	URL         string                   `json:"url"`
	ValidatedAt time.Time                `json:"validated_at"`
	LastTest    MarketingWebhookDelivery `json:"last_test"`
}
