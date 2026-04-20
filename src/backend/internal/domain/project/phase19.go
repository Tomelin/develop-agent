package project

import "time"

type Phase19RunInput struct {
	IncludeFlowA bool `json:"include_flow_a"`
	IncludeFlowB bool `json:"include_flow_b"`
	IncludeFlowC bool `json:"include_flow_c"`
}

type LLMJudgePhaseScore struct {
	Phase        int      `json:"phase"`
	ArtifactType string   `json:"artifact_type"`
	Criteria     []string `json:"criteria"`
	Score        float64  `json:"score"`
	Notes        []string `json:"notes,omitempty"`
}

type Phase19QualityReport struct {
	GeneratedAt          time.Time            `json:"generated_at"`
	ProjectID            string               `json:"project_id"`
	FlowCoverage         map[string]bool      `json:"flow_coverage"`
	TotalE2EScenarios    int                  `json:"total_e2e_scenarios"`
	ContractsValidated   int                  `json:"contracts_validated"`
	JudgeScores          []LLMJudgePhaseScore `json:"judge_scores"`
	JudgeAverageByPhase  map[string]float64   `json:"judge_average_by_phase"`
	JudgeOverallAverage  float64              `json:"judge_overall_average"`
	StagingReady         bool                 `json:"staging_ready"`
	DeployChecklistReady bool                 `json:"deploy_checklist_ready"`
	TroubleshootingReady bool                 `json:"troubleshooting_ready"`
}

type Phase19DeliveryReport struct {
	GeneratedAt   time.Time            `json:"generated_at"`
	ProjectID     string               `json:"project_id"`
	Artifacts     []string             `json:"artifacts"`
	QualityReport Phase19QualityReport `json:"quality_report"`
}

type AdminQualityReport struct {
	GeneratedAt                time.Time          `json:"generated_at"`
	ProjectSampleSize          int                `json:"project_sample_size"`
	TestCoveragePercent        float64            `json:"test_coverage_percent"`
	TriadSuccessRatePercent    float64            `json:"triad_success_rate_percent"`
	JudgeAverageByPhase        map[string]float64 `json:"judge_average_by_phase"`
	AvgExecutionMinutesByPhase map[string]float64 `json:"avg_execution_minutes_by_phase"`
	AverageCostByFlowType      map[string]float64 `json:"average_cost_by_flow_type"`
	PlatformUptime30dPercent   float64            `json:"platform_uptime_30d_percent"`
	ProjectsCompleted          int                `json:"projects_completed"`
	ProjectsAbandoned          int                `json:"projects_abandoned"`
	Notes                      []string           `json:"notes,omitempty"`
}
