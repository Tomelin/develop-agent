export type AgentProvider = "OPENAI" | "ANTHROPIC" | "GOOGLE" | "OLLAMA";

export type AgentStatus = "IDLE" | "RUNNING" | "PAUSED" | "QUEUED" | "ERROR" | "COMPLETED";

export type AgentSkill =
  | "PROJECT_CREATION"
  | "ENGINEERING"
  | "ARCHITECTURE"
  | "PLANNING"
  | "DEVELOPMENT_FRONTEND"
  | "DEVELOPMENT_BACKEND"
  | "TESTING"
  | "SECURITY"
  | "DOCUMENTATION"
  | "DEVOPS"
  | "LANDING_PAGE"
  | "MARKETING";

export interface Agent {
  id: string;
  name: string;
  description: string;
  provider: AgentProvider;
  model: string;
  system_prompts: string[];
  skills: AgentSkill[];
  enabled: boolean;
  api_key_ref?: string;
  status: AgentStatus;
  created_at: string;
  updated_at: string;
}

export interface AgentListResponse {
  items: Agent[];
  total: number;
  page: number;
  size: number;
  pages: number;
}
