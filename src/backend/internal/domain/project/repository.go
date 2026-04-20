package project

import "context"

type ProjectListFilter struct {
	OwnerID  string
	Status   ProjectStatus
	FlowType FlowType
	Page     int64
	Limit    int64
}

type ProjectRepository interface {
	EnsureIndexes(ctx context.Context) error
	Create(ctx context.Context, p *Project) error
	FindByID(ctx context.Context, id string) (*Project, error)
	FindByOwner(ctx context.Context, filter ProjectListFilter) ([]*Project, int64, error)
	FindDashboardByOwner(ctx context.Context, filter ProjectListFilter) ([]*Project, int64, error)
	Update(ctx context.Context, p *Project) error
	Archive(ctx context.Context, id string, ownerID string) error
	UpdatePhase(ctx context.Context, projectID, ownerID string, phase PhaseExecution) error
	UpdateSpecMD(ctx context.Context, projectID, ownerID, specMD string) error
}

type TaskListFilter struct {
	ProjectID  string
	Type       TaskType
	Complexity TaskComplexity
	Status     TaskStatus
}

type TaskRepository interface {
	EnsureIndexes(ctx context.Context) error
	BulkCreate(ctx context.Context, tasks []*Task) error
	ListByProject(ctx context.Context, filter TaskListFilter) ([]*Task, error)
	RoadmapSummary(ctx context.Context, projectID string) (*RoadmapSummary, error)
	UpdateStatus(ctx context.Context, projectID, taskID string, status TaskStatus) error
}

type CodeFileRepository interface {
	EnsureIndexes(ctx context.Context) error
	Upsert(ctx context.Context, file *CodeFile) error
	ListByProject(ctx context.Context, projectID string) ([]*CodeFile, error)
	FindByID(ctx context.Context, projectID, fileID string) (*CodeFile, error)
	CountAutoRejections(ctx context.Context, projectID string) (int64, error)
	IncrementAutoRejections(ctx context.Context, projectID string) error
}
