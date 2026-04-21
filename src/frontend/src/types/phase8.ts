export type PhaseTrack = "FRONTEND" | "BACKEND";
export type TrackExecutionStatus = "PENDING" | "RUNNING" | "REVIEW" | "COMPLETED" | "ERROR";

export interface PhaseTrackStatus {
  track: PhaseTrack;
  status: TrackExecutionStatus;
  execution_id?: string;
  feedbacks_used: number;
  feedbacks_limit: number;
  updated_at?: string;
}

export interface ArtifactVersion {
  version: number;
  content: string;
  created_at: string;
}

export interface PhaseArtifact {
  id: string;
  phase_number: number;
  track: PhaseTrack;
  type: string;
  title: string;
  current_content: string;
  versions?: ArtifactVersion[];
  updated_at: string;
}

export interface TriadStepRuntime {
  step: "PRODUCER" | "REVIEWER" | "REFINER";
  agent_name: string;
  provider: string;
  model: string;
  status: TrackExecutionStatus;
  tokens_used?: number;
  duration_ms?: number;
  partial_output?: string;
}

export interface TriadTrackRuntime {
  track: PhaseTrack;
  status: TrackExecutionStatus;
  steps: TriadStepRuntime[];
  updated_at: string;
}

export interface TrackFeedbackItem {
  id: string;
  track: PhaseTrack;
  content: string;
  created_at: string;
}

export interface NotificationItem {
  id: string;
  project_id: string;
  phase_number: number;
  type: string;
  message: string;
  read: boolean;
  created_at: string;
}
