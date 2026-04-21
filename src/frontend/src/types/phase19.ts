export interface AdminQualityReport {
  generated_at: string;
  project_sample_size: number;
  test_coverage_percent: number;
  triad_success_rate_percent: number;
  judge_average_by_phase: Record<string, number>;
  avg_execution_minutes_by_phase: Record<string, number>;
  average_cost_by_flow_type: Record<string, number>;
  platform_uptime_30d_percent: number;
  projects_completed: number;
  projects_abandoned: number;
  notes?: string[];
}

