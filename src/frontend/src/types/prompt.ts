export type PromptGroup =
  | "GLOBAL"
  | "PROJECT_CREATION"
  | "ENGINEERING"
  | "ARCHITECTURE"
  | "PLANNING"
  | "DEVELOPMENT"
  | "TESTING"
  | "SECURITY"
  | "DOCUMENTATION"
  | "DEVOPS"
  | "LANDING_PAGE"
  | "MARKETING";

export interface UserPrompt {
  id: string;
  user_id: string;
  title: string;
  content: string;
  group: PromptGroup;
  priority: number;
  enabled: boolean;
  tags: string[];
  created_at: string;
  updated_at: string;
}

export interface PromptListResponse {
  items: UserPrompt[];
  total?: number;
  page?: number;
  size?: number;
}

export interface PromptTemplate {
  id: string;
  title: string;
  content: string;
  group: PromptGroup;
  tags: string[];
  description?: string;
  category?: string;
}

export interface PromptPreviewBlock {
  source: "SYSTEM" | "GLOBAL" | "GROUP" | "RAG" | "PHASE_INSTRUCTION";
  title: string;
  content: string;
}

export interface PromptPreviewResponse {
  group: PromptGroup;
  agent_id?: string;
  blocks: PromptPreviewBlock[];
  composed_prompt: string;
  token_estimate?: number;
}

export interface PromptReorderItem {
  id: string;
  priority: number;
}

export interface PromptImportPayload {
  mode: "MERGE" | "REPLACE";
  prompts: Array<{
    title: string;
    content: string;
    group: PromptGroup;
    priority: number;
    enabled: boolean;
    tags: string[];
  }>;
}

export const PROMPT_GROUPS: Array<{ label: string; value: PromptGroup; description: string }> = [
  { label: "Global", value: "GLOBAL", description: "Diretrizes aplicadas em todas as phases" },
  { label: "Project Creation", value: "PROJECT_CREATION", description: "Definição e descoberta de produto" },
  { label: "Engineering", value: "ENGINEERING", description: "Regras funcionais e não-funcionais" },
  { label: "Architecture", value: "ARCHITECTURE", description: "Modelagem e padrões arquiteturais" },
  { label: "Planning", value: "PLANNING", description: "Roadmap, épicos e tasks" },
  { label: "Development", value: "DEVELOPMENT", description: "Implementação front/back" },
  { label: "Testing", value: "TESTING", description: "Estratégia de testes e TDD" },
  { label: "Security", value: "SECURITY", description: "Hardening, OWASP, compliance" },
  { label: "Documentation", value: "DOCUMENTATION", description: "Documentação técnica e operacional" },
  { label: "DevOps", value: "DEVOPS", description: "Pipelines, containers e deploy" },
  { label: "Landing Page", value: "LANDING_PAGE", description: "Páginas de conversão e copy" },
  { label: "Marketing", value: "MARKETING", description: "Campanhas e performance" },
];
