export type MarketingChannel = "linkedin" | "instagram" | "google-ads";

export interface MarketingManualBrief {
  product_name: string;
  tagline?: string;
  problem_solved: string;
  target_audience: string;
  main_benefits: string[];
  differentials?: string[];
  business_model?: string;
  pricing?: string;
  market_type?: string;
  communication_tone?: string;
  primary_cta?: string;
  secondary_cta?: string;
  competitor_references?: string[];
}

export interface Phase15RunRequest {
  use_linked_project: boolean;
  channels: MarketingChannel[];
  monthly_budget_usd: number;
  manual_brief: MarketingManualBrief;
}

export interface MarketingChannelSummary {
  channel: MarketingChannel;
  pieces: number;
  expected_ctr: string;
  expected_conversion: string;
  budget_usd: number;
}

export interface Phase15DeliveryReport {
  generated_at: string;
  project_id: string;
  brief_source: string;
  channels: MarketingChannel[];
  total_pieces: number;
  artifact_paths: string[];
  channel_summaries: MarketingChannelSummary[];
  warnings?: string[];
}

export interface MarketingWebhookResult {
  url: string;
  validated_at: string;
  last_test: {
    timestamp: string;
    status: string;
    response_status: number;
    error?: string;
  };
}
