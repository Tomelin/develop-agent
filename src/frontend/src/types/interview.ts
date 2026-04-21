export type InterviewRole = "USER" | "ASSISTANT";
export type InterviewStatus =
  | "ACTIVE"
  | "AWAITING_CONFIRMATION"
  | "COMPLETED"
  | "ABANDONED";

export interface InterviewMessage {
  id?: string;
  role: InterviewRole;
  content: string;
  timestamp: string;
}

export interface InterviewCoverageItem {
  key: string;
  title: string;
  status: "DONE" | "IN_PROGRESS" | "PENDING";
}

export interface InterviewSession {
  id: string;
  project_id: string;
  status: InterviewStatus;
  iteration_count: number;
  max_iterations: number;
  completed_at?: string;
  messages: InterviewMessage[];
  coverage?: InterviewCoverageItem[];
  vision_markdown?: string;
  vision_generated_at?: string;
}
