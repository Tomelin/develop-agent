export type ProjectStatus = "DRAFT" | "IN_PROGRESS" | "PAUSED" | "COMPLETED" | "ARCHIVED";
export type FlowType = "A" | "B" | "C";
export type PhaseStatus = "PENDING" | "IN_PROGRESS" | "REVIEW" | "COMPLETED" | "REJECTED";

export interface Project {
  id: string;
  name: string;
  description: string;
  status: ProjectStatus;
  flow_type: FlowType;
  linked_project_id?: string;
  current_phase: number;
  progress_percentage: number;
  tokens_used: number;
  dynamic_mode: boolean;
  created_at: string;
  updated_at: string;
}

export interface ProjectListResponse {
  items: Project[];
  total: number;
  page: number;
  size: number;
  pages: number;
}

export interface ProjectCreateRequest {
  name: string;
  description: string;
  flow_type: FlowType;
  dynamic_mode: boolean;
  linked_project_id?: string;
}
