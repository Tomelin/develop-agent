package project

import (
	"context"
	"testing"

	domainproject "github.com/develop-agent/backend/internal/domain/project"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestRoadmapSchemaValidatorValidate(t *testing.T) {
	validator := &RoadmapSchemaValidator{}
	valid := []byte(`{"project_id":"p1","phases":[{"id":"phase-1","name":"Setup","description":"d","order":1,"epics":[{"id":"epic-1","title":"Core","description":"d","tasks":[{"id":"task-1","title":"API","description":"d","type":"BACKEND","complexity":"MEDIUM","estimated_hours":8,"track":"BACKEND","dependencies":[]},{"id":"task-2","title":"UI","description":"d","type":"FRONTEND","complexity":"LOW","estimated_hours":5,"track":"FRONTEND","dependencies":["task-1"]}]}]}]}`)
	if _, err := validator.Validate(valid); err != nil {
		t.Fatalf("expected valid roadmap, got error: %+v", err)
	}
}

func TestRoadmapSchemaValidatorRejectsInvalidDependencyAndHours(t *testing.T) {
	validator := &RoadmapSchemaValidator{}
	invalid := []byte(`{"project_id":"p1","phases":[{"id":"phase-1","name":"Setup","description":"d","order":1,"epics":[{"id":"epic-1","title":"Core","description":"d","tasks":[{"id":"task-1","title":"API","description":"d","type":"BACKEND","complexity":"MEDIUM","estimated_hours":0,"track":"BACKEND","dependencies":["task-999"]}]}]}]}`)
	_, err := validator.Validate(invalid)
	if err == nil || len(err.Issues) < 2 {
		t.Fatalf("expected multiple validation issues, got: %+v", err)
	}
}

type memTaskRepo struct {
	tasks []*domainproject.Task
}

func (m *memTaskRepo) EnsureIndexes(context.Context) error { return nil }
func (m *memTaskRepo) BulkCreate(_ context.Context, tasks []*domainproject.Task) error {
	m.tasks = append(m.tasks, tasks...)
	return nil
}
func (m *memTaskRepo) ListByProject(context.Context, domainproject.TaskListFilter) ([]*domainproject.Task, error) {
	return m.tasks, nil
}
func (m *memTaskRepo) RoadmapSummary(context.Context, string) (*domainproject.RoadmapSummary, error) {
	return &domainproject.RoadmapSummary{}, nil
}
func (m *memTaskRepo) UpdateStatus(context.Context, string, string, domainproject.TaskStatus) error {
	return nil
}

func TestApproveRoadmapPhaseIngestsAndCompletesPhase(t *testing.T) {
	owner := bson.NewObjectID()
	p, _ := domainproject.NewProject("P", "", domainproject.FlowSoftware, owner, false, nil)
	p.Phases[0].Status = domainproject.PhaseCompleted
	p.Phases[3].Status = domainproject.PhaseInProgress

	projectRepo := &memProjectRepo{project: p}
	taskRepo := &memTaskRepo{}
	svc := NewService(projectRepo, taskRepo)

	roadmap := []byte(`{"project_id":"` + p.ID.Hex() + `","phases":[{"id":"phase-1","name":"Setup","description":"d","order":1,"epics":[{"id":"epic-1","title":"Core","description":"d","tasks":[{"id":"task-1","title":"API","description":"d","type":"BACKEND","complexity":"MEDIUM","estimated_hours":8,"track":"BACKEND","dependencies":[]}]}]}]}`)

	phase, result, err := svc.ApproveRoadmapPhase(context.Background(), p.ID.Hex(), owner.Hex(), roadmap)
	if err != nil {
		t.Fatalf("unexpected error approving roadmap: %v", err)
	}
	if phase.Status != domainproject.PhaseCompleted {
		t.Fatalf("expected phase completed, got %s", phase.Status)
	}
	if result.TaskCount != 1 || len(taskRepo.tasks) != 1 {
		t.Fatalf("expected one ingested task, got result=%+v tasks=%d", result, len(taskRepo.tasks))
	}
	if projectRepo.project.RoadmapJSON == "" {
		t.Fatal("expected project roadmap_json to be persisted")
	}
}
