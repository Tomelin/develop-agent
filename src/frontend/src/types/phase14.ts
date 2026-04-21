export type LandingTone = "profissional" | "moderno" | "descontraído" | "inspirador";
export type LandingTheme = "light" | "dark";
export type LandingOutputFormat = "html" | "nextjs";

export interface LandingPageManualBrief {
  product_name: string;
  problem_solved: string;
  target_audience: string;
  unique_value_proposed: string;
  key_features: string[];
  color_palette: string[];
  theme: LandingTheme;
  communication_tone: LandingTone;
  language: "pt-BR" | "en-US" | "es";
  preferred_typography: string;
  output_format: LandingOutputFormat;
  primary_keyword?: string;
  primary_cta?: string;
  secondary_cta?: string;
  social_proof_highlight?: string;
}

export interface Phase14RunRequest {
  use_linked_project: boolean;
  generate_variants: boolean;
  variant_count: number;
  manual_brief: LandingPageManualBrief;
}

export interface LandingPageVariantReport {
  name: string;
  path: string;
  conversion_score: number;
}

export interface Phase14DeliveryReport {
  generated_at: string;
  project_id: string;
  brief_source: string;
  output_format: string;
  conversion_score: number;
  artifact_paths: string[];
  variants?: LandingPageVariantReport[];
  prioritized_issues?: string[];
}
