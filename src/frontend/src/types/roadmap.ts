export interface RoadmapSummary {
  tasks_by_type: Record<string, number>;
  tasks_by_complexity: Record<string, number>;
  estimated_hours_by_type: Record<string, number>;
  estimated_hours_by_phase: Record<string, number>;
  total_phases: number;
  total_epics: number;
  critical_path_hours: number;
}

import { RoadmapTask } from "./task";

export interface RoadmapEpic {
  id: string;
  title: string;
  description: string;
  tasks: RoadmapTask[];
  status?: "PENDING" | "IN_PROGRESS" | "COMPLETED";
  progress?: number;
}

export interface RoadmapPhase {
  id: string;
  name: string;
  description: string;
  order: number;
  epics: RoadmapEpic[];
}

export interface RoadmapData {
  project_id: string;
  phases: RoadmapPhase[];
}
