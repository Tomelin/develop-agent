import { Agent, AgentProvider } from "./agent";

export type TriadRole = "PRODUCER" | "REVIEWER" | "REFINER";
export type SelectionMode = "DYNAMIC" | "FIXED";

export interface TriadAgentSummary {
  id: string;
  name: string;
  provider: AgentProvider;
  model: string;
  avatar_url?: string;
}

export interface TriadSelection {
  phase_name: string;
  phase_number: number;
  execution_id?: string;
  mode: SelectionMode;
  producer: TriadAgentSummary;
  reviewer: TriadAgentSummary;
  refiner: TriadAgentSummary;
  selection_timestamp: string;
}

export interface TriadSelectionLog {
  id: string;
  project_id: string;
  phase_number: number;
  execution_id: string;
  mode: SelectionMode;
  candidate_agents: TriadAgentSummary[];
  selected_triad: {
    producer: TriadAgentSummary;
    reviewer: TriadAgentSummary;
    refiner: TriadAgentSummary;
  };
  selection_reason: string;
  timestamp: string;
}

export interface DynamicModeConfig {
  enabled: boolean;
  fixed_agents?: Partial<Record<TriadRole, string>>;
}

export interface DynamicModePreview {
  eligible_agents: number;
  triad: {
    producer: TriadAgentSummary;
    reviewer: TriadAgentSummary;
    refiner: TriadAgentSummary;
  };
  notes: string[];
}

export interface DiversityProviderUsage {
  provider: AgentProvider;
  usage_percentage: number;
  count: number;
}

export interface DiversityRoleDistribution {
  role: TriadRole;
  model: string;
  uses: number;
}

export interface DiversityMetrics {
  project_id: string;
  diversity_score: number;
  providers: DiversityProviderUsage[];
  models: string[];
  full_diversity_triads: number;
  repeated_provider_triads: number;
  role_distribution: DiversityRoleDistribution[];
}

export interface FeatureFlag {
  key: string;
  enabled: boolean;
  description: string;
}

export interface AdminPlatformSettings {
  workers: {
    max_concurrency: number;
    agent_timeout_seconds: number;
    triad_timeout_seconds: number;
  };
  models: {
    default_model: string;
    spec_generation_model: string;
  };
  limits: {
    max_projects_per_user: number;
    max_parallel_phases_per_user: number;
    max_spec_tokens: number;
  };
  retry: {
    max_attempts: number;
    backoff_seconds: number;
  };
}

export interface PhaseAgentMatrixRow {
  phase_key: string;
  phase_label: string;
  producer_id?: string | null;
  reviewer_id?: string | null;
  refiner_id?: string | null;
  dynamic?: boolean;
}

export interface PhaseAgentMatrixResponse {
  rows: PhaseAgentMatrixRow[];
}

export interface CostPreviewResponse {
  monthly_estimated_usd: number;
  note: string;
}

export type ProviderPalette = Record<AgentProvider, { bg: string; text: string; ring: string }>;

export const providerPalette: ProviderPalette = {
  OPENAI: { bg: "bg-blue-500/15", text: "text-blue-400", ring: "ring-blue-400/40" },
  ANTHROPIC: { bg: "bg-orange-500/15", text: "text-orange-400", ring: "ring-orange-400/40" },
  GOOGLE: { bg: "bg-emerald-500/15", text: "text-emerald-400", ring: "ring-emerald-400/40" },
  OLLAMA: { bg: "bg-purple-500/15", text: "text-purple-400", ring: "ring-purple-400/40" },
};

export const phaseMatrixDefaults: PhaseAgentMatrixRow[] = [
  ...Array.from({ length: 9 }).map((_, i) => ({
    phase_key: `PHASE_${i + 1}`,
    phase_label: `Fase ${String(i + 1).padStart(2, "0")}`,
    producer_id: null,
    reviewer_id: null,
    refiner_id: null,
    dynamic: true,
  })),
  { phase_key: "FLOW_B", phase_label: "Fluxo B (Landing)", producer_id: null, reviewer_id: null, refiner_id: null, dynamic: true },
  { phase_key: "FLOW_C", phase_label: "Fluxo C (Marketing)", producer_id: null, reviewer_id: null, refiner_id: null, dynamic: true },
];

export const TRIAD_ROLES: TriadRole[] = ["PRODUCER", "REVIEWER", "REFINER"];

export const roleLabel: Record<TriadRole, string> = {
  PRODUCER: "Produtor",
  REVIEWER: "Revisor",
  REFINER: "Refinador",
};

export const dynamicFeatureFlagsSeed: FeatureFlag[] = [
  { key: "DYNAMIC_MODE_ENABLED", enabled: true, description: "Disponibiliza o Modo Dinâmico para os projetos." },
  { key: "FLOW_B_ENABLED", enabled: true, description: "Habilita o fluxo de Landing Page." },
  { key: "FLOW_C_ENABLED", enabled: true, description: "Habilita o fluxo de Marketing." },
  { key: "DEVOPS_PHASE_ENABLED", enabled: false, description: "Habilita a Fase 9 (DevOps)." },
  { key: "BILLING_PANEL_ENABLED", enabled: true, description: "Exibe o painel de billing para usuários." },
  { key: "AUTO_REJECTION_ENABLED", enabled: true, description: "Permite gatilho de rejeição automática entre fases." },
];

export const roleKeyMap: Record<TriadRole, "producer_id" | "reviewer_id" | "refiner_id"> = {
  PRODUCER: "producer_id",
  REVIEWER: "reviewer_id",
  REFINER: "refiner_id",
};

export const asTriadSummary = (agent: Agent): TriadAgentSummary => ({
  id: agent.id,
  name: agent.name,
  provider: agent.provider,
  model: agent.model,
});
