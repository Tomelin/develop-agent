export type TaskType = "FRONTEND" | "BACKEND" | "INFRA" | "TEST" | "DOC";
export type TaskComplexity = "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
export type TaskStatus = "TODO" | "IN_PROGRESS" | "DONE" | "BLOCKED";

export interface Task {
  id: string;
  project_id: string;
  phase_id: string;
  epic_id?: string;
  title: string;
  description: string;
  type: TaskType;
  complexity: TaskComplexity;
  estimated_hours: number;
  status: TaskStatus;
  assigned_agent_id?: string;
  created_at: string;
  updated_at: string;
}

export interface TaskListResponse {
  items: Task[];
  total: number;
  page: number;
  size: number;
  pages: number;
}
