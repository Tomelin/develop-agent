package project

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CodeFile struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID   bson.ObjectID `bson:"project_id" json:"project_id"`
	Path        string        `bson:"path" json:"path"`
	Content     string        `bson:"content" json:"content"`
	TaskID      string        `bson:"task_id" json:"task_id"`
	Language    string        `bson:"language" json:"language"`
	Version     time.Time     `bson:"version" json:"version"`
	PhaseNumber int           `bson:"phase_number" json:"phase_number"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
}

type CodeSymbol struct {
	Name    string `json:"name"`
	Kind    string `json:"kind"`
	Source  string `json:"source"`
	Backend bool   `json:"backend"`
}

type CodeContextManifest struct {
	Files            []CodeContextFile `json:"files"`
	Symbols          []CodeSymbol      `json:"symbols"`
	Dependencies     []string          `json:"dependencies"`
	EnvironmentHints []string          `json:"environment_hints"`
	ApproxTokens     int               `json:"approx_tokens"`
}

type CodeContextFile struct {
	Path     string `json:"path"`
	Language string `json:"language"`
	Purpose  string `json:"purpose"`
}

type Phase5Summary struct {
	TotalTasks            int64   `json:"total_tasks"`
	DoneTasks             int64   `json:"done_tasks"`
	InProgressTasks       int64   `json:"in_progress_tasks"`
	BlockedTasks          int64   `json:"blocked_tasks"`
	TodoTasks             int64   `json:"todo_tasks"`
	BackendFiles          int64   `json:"backend_files"`
	FrontendFiles         int64   `json:"frontend_files"`
	GeneratedLinesOfCode  int64   `json:"generated_lines_of_code"`
	AverageTaskMinutes    float64 `json:"average_task_minutes"`
	AutoRejections        int64   `json:"auto_rejections"`
	TotalPhaseTokens      int64   `json:"total_phase_tokens"`
	ExecutionMode         string  `json:"execution_mode"`
	CompletionPercent     float64 `json:"completion_percent"`
	LastExecutionUnixTime int64   `json:"last_execution_unix_time,omitempty"`
}

type CatastrophicFailureReport struct {
	CoveragePercent      float64  `json:"coverage_percent"`
	MaxCVSS              float64  `json:"max_cvss"`
	CredentialsExposed   bool     `json:"credentials_exposed"`
	CompilationFailed    bool     `json:"compilation_failed"`
	FailureDescriptions  []string `json:"failure_descriptions"`
	SourcePhase          int      `json:"source_phase"`
	RequestedBySystemTag string   `json:"requested_by_system_tag"`
}
