package project

type RoadmapTask struct {
	ID             string         `json:"id"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Type           TaskType       `json:"type"`
	Complexity     TaskComplexity `json:"complexity"`
	EstimatedHours int            `json:"estimated_hours"`
	Track          Track          `json:"track"`
	Dependencies   []string       `json:"dependencies,omitempty"`
}

type RoadmapEpic struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Tasks       []RoadmapTask `json:"tasks"`
}

type RoadmapPhase struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Order       int           `json:"order"`
	Epics       []RoadmapEpic `json:"epics"`
}

type RoadmapDocument struct {
	ProjectID string         `json:"project_id"`
	Phases    []RoadmapPhase `json:"phases"`
}

type ValidationIssue struct {
	Path    string `json:"path"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type RoadmapValidationError struct {
	Issues []ValidationIssue `json:"issues"`
}

func (e *RoadmapValidationError) Error() string {
	if len(e.Issues) == 0 {
		return "roadmap validation failed"
	}
	return e.Issues[0].Message
}

func (e *RoadmapValidationError) HasIssues() bool { return len(e.Issues) > 0 }

type RoadmapSummary struct {
	TotalTasks              int64                    `json:"total_tasks"`
	TotalByType             map[TaskType]int64       `json:"total_by_type"`
	TotalByComplexity       map[TaskComplexity]int64 `json:"total_by_complexity"`
	HoursByType             map[TaskType]float64     `json:"hours_by_type"`
	HoursByPhase            map[string]float64       `json:"hours_by_phase"`
	PhaseCount              int64                    `json:"phase_count"`
	EpicCount               int64                    `json:"epic_count"`
	EstimatedCriticalPathHR float64                  `json:"estimated_critical_path_hours"`
}
