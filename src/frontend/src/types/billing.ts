export interface BillingGroupedItem {
  key: string;
  cost_usd: number;
  tokens: number;
  executions: number;
}

export interface BillingSummary {
  total_cost_usd: number;
  total_tokens: number;
  by_project: BillingGroupedItem[];
  by_model: BillingGroupedItem[];
}

export interface BillingProjectDetails {
  project_id: string;
  by_phase: BillingGroupedItem[];
  by_agent: BillingGroupedItem[];
  by_model: BillingGroupedItem[];
  total_usd: number;
}

export interface BillingPricingItem {
  provider: string;
  model: string;
  prompt_price_per_million_tokens: number;
  completion_price_per_million_tokens: number;
}

export interface BillingPricingTable {
  last_updated: string;
  models: BillingPricingItem[];
}

export interface BillingRecord {
  id: string;
  organization_id: string;
  project_id: string;
  user_id: string;
  phase_number: number;
  phase_name: string;
  triad_role: "PRODUCER" | "REVIEWER" | "REFINER";
  agent_id: string;
  agent_name: string;
  provider: string;
  model: string;
  prompt_tokens: number;
  completion_tokens: number;
  total_tokens: number;
  price_per_million_prompt_tokens: number;
  price_per_million_completion_tokens: number;
  estimated_cost_usd: number;
  duration_ms: number;
  is_auto_rejection: boolean;
  timestamp: string;
}
