package project

import "time"

type LandingPageManualBrief struct {
	ProductName          string   `json:"product_name"`
	ProblemSolved        string   `json:"problem_solved"`
	TargetAudience       string   `json:"target_audience"`
	UniqueValueProposed  string   `json:"unique_value_proposed"`
	KeyFeatures          []string `json:"key_features"`
	ColorPalette         []string `json:"color_palette,omitempty"`
	Theme                string   `json:"theme,omitempty"`
	CommunicationTone    string   `json:"communication_tone,omitempty"`
	Language             string   `json:"language,omitempty"`
	PreferredTypography  string   `json:"preferred_typography,omitempty"`
	OutputFormat         string   `json:"output_format,omitempty"`
	PrimaryKeyword       string   `json:"primary_keyword,omitempty"`
	PrimaryCTA           string   `json:"primary_cta,omitempty"`
	SecondaryCTA         string   `json:"secondary_cta,omitempty"`
	SocialProofHighlight string   `json:"social_proof_highlight,omitempty"`
}

type Phase14RunInput struct {
	UseLinkedProject bool                   `json:"use_linked_project"`
	ManualBrief      LandingPageManualBrief `json:"manual_brief"`
	GenerateVariants bool                   `json:"generate_variants"`
	VariantCount     int                    `json:"variant_count"`
}

type LandingPageBrief struct {
	Source               string   `json:"source"`
	ProductName          string   `json:"product_name"`
	Tagline              string   `json:"tagline,omitempty"`
	ProblemSolved        string   `json:"problem_solved"`
	TargetAudience       string   `json:"target_audience"`
	UniqueValueProposed  string   `json:"unique_value_proposed"`
	KeyFeatures          []string `json:"key_features"`
	BusinessModel        string   `json:"business_model,omitempty"`
	RelevantIntegrations []string `json:"relevant_integrations,omitempty"`
	ColorPalette         []string `json:"color_palette,omitempty"`
	Theme                string   `json:"theme,omitempty"`
	CommunicationTone    string   `json:"communication_tone,omitempty"`
	Language             string   `json:"language,omitempty"`
	PreferredTypography  string   `json:"preferred_typography,omitempty"`
	OutputFormat         string   `json:"output_format"`
	PrimaryKeyword       string   `json:"primary_keyword,omitempty"`
	PrimaryCTA           string   `json:"primary_cta"`
	SecondaryCTA         string   `json:"secondary_cta,omitempty"`
	SocialProofHighlight string   `json:"social_proof_highlight,omitempty"`
}

type LandingPageVariantReport struct {
	Name            string  `json:"name"`
	Path            string  `json:"path"`
	ConversionScore float64 `json:"conversion_score"`
}

type Phase14DeliveryReport struct {
	GeneratedAt       time.Time                  `json:"generated_at"`
	ProjectID         string                     `json:"project_id"`
	BriefSource       string                     `json:"brief_source"`
	OutputFormat      string                     `json:"output_format"`
	ConversionScore   float64                    `json:"conversion_score"`
	ArtifactPaths     []string                   `json:"artifact_paths"`
	Variants          []LandingPageVariantReport `json:"variants,omitempty"`
	PrioritizedIssues []string                   `json:"prioritized_issues,omitempty"`
}
