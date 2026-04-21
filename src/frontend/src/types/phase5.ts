import { Task } from "./task";

export type Phase5ExecutionMode = "AUTOMATIC" | "MANUAL";

export interface Phase5Summary {
  total_tasks: number;
  done_tasks: number;
  in_progress_tasks: number;
  blocked_tasks: number;
  todo_tasks: number;
  backend_files: number;
  frontend_files: number;
  generated_lines_of_code: number;
  average_task_minutes: number;
  auto_rejections: number;
  total_phase_tokens: number;
  execution_mode: Phase5ExecutionMode;
  completion_percent: number;
  last_execution_unix_time?: number;
}

export interface Phase5CodeFile {
  id: string;
  project_id: string;
  path: string;
  content: string;
  task_id: string;
  language: string;
  version: string;
  phase_number: number;
  created_at: string;
  updated_at: string;
}

export interface Phase5CodeContextFile {
  path: string;
  language: string;
  purpose: string;
}

export interface Phase5CodeSymbol {
  name: string;
  kind: string;
  source: string;
  backend: boolean;
}

export interface Phase5CodeContext {
  files: Phase5CodeContextFile[];
  symbols: Phase5CodeSymbol[];
  dependencies: string[];
  environment_hints: string[];
  approx_tokens: number;
}

export interface Phase5FileListResponse {
  items: Phase5CodeFile[];
}

export interface Phase5TaskExecutionEvent {
  timestamp: string;
  taskId?: string;
  taskTitle?: string;
  level: "INFO" | "SUCCESS" | "ERROR";
  message: string;
}

export interface Phase5TaskWithFile extends Task {
  generatedFile?: Phase5CodeFile;
}
