package project

import (
	"context"
	"errors"
	"testing"
	"time"

	domainproject "github.com/develop-agent/backend/internal/domain/project"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type memCodeFileRepo struct {
	files          map[string]*domainproject.CodeFile
	autoRejections int64
}

func (m *memCodeFileRepo) EnsureIndexes(context.Context) error { return nil }
func (m *memCodeFileRepo) Upsert(_ context.Context, file *domainproject.CodeFile) error {
	if m.files == nil {
		m.files = map[string]*domainproject.CodeFile{}
	}
	clone := *file
	if clone.ID.IsZero() {
		clone.ID = bson.NewObjectID()
	}
	clone.UpdatedAt = time.Now().UTC()
	m.files[clone.Path] = &clone
	return nil
}
func (m *memCodeFileRepo) ListByProject(context.Context, string) ([]*domainproject.CodeFile, error) {
	out := make([]*domainproject.CodeFile, 0, len(m.files))
	for _, f := range m.files {
		out = append(out, f)
	}
	return out, nil
}
func (m *memCodeFileRepo) FindByID(_ context.Context, _ string, fileID string) (*domainproject.CodeFile, error) {
	for _, f := range m.files {
		if f.ID.Hex() == fileID {
			return f, nil
		}
	}
	return nil, errors.New("not found")
}
func (m *memCodeFileRepo) CountAutoRejections(context.Context, string) (int64, error) {
	return m.autoRejections, nil
}
func (m *memCodeFileRepo) IncrementAutoRejections(context.Context, string) error {
	m.autoRejections++
	return nil
}

func TestSetExecutionMode(t *testing.T) {
	owner := bson.NewObjectID()
	p, _ := domainproject.NewProject("P", "", domainproject.FlowSoftware, owner, owner, false, nil)
	svc := NewDevelopmentService(&memProjectRepo{project: p}, &noopTaskRepo{}, &memCodeFileRepo{})

	mode, err := svc.SetExecutionMode(context.Background(), p.ID.Hex(), owner.Hex(), domainproject.ExecutionModeAutomatic)
	if err != nil {
		t.Fatalf("set mode: %v", err)
	}
	if mode != domainproject.ExecutionModeAutomatic {
		t.Fatalf("unexpected mode: %s", mode)
	}
}

func TestExecuteTaskGeneratesAndCompletes(t *testing.T) {
	owner := bson.NewObjectID()
	p, _ := domainproject.NewProject("P", "", domainproject.FlowSoftware, owner, owner, false, nil)
	taskID := bson.NewObjectID()
	taskRepo := &inMemoryTaskRepo{items: []*domainproject.Task{{
		ID:         taskID,
		ProjectID:  p.ID,
		Title:      "create user usecase",
		Type:       domainproject.TaskTypeBackend,
		Complexity: domainproject.ComplexityMedium,
		Status:     domainproject.TaskTodo,
		UpdatedAt:  time.Now().UTC(),
		CreatedAt:  time.Now().UTC(),
	}}}
	fileRepo := &memCodeFileRepo{}
	svc := NewDevelopmentService(&memProjectRepo{project: p}, taskRepo, fileRepo)

	if err := svc.ExecuteTask(context.Background(), p.ID.Hex(), owner.Hex(), taskID.Hex()); err != nil {
		t.Fatalf("execute task: %v", err)
	}
	if taskRepo.items[0].Status != domainproject.TaskDone {
		t.Fatalf("expected task done, got %s", taskRepo.items[0].Status)
	}
	if len(fileRepo.files) == 0 {
		t.Fatal("expected generated file")
	}
}

type inMemoryTaskRepo struct {
	items []*domainproject.Task
}

func (m *inMemoryTaskRepo) EnsureIndexes(context.Context) error                     { return nil }
func (m *inMemoryTaskRepo) BulkCreate(context.Context, []*domainproject.Task) error { return nil }
func (m *inMemoryTaskRepo) ListByProject(_ context.Context, filter domainproject.TaskListFilter) ([]*domainproject.Task, error) {
	if filter.Status == "" {
		return m.items, nil
	}
	out := make([]*domainproject.Task, 0)
	for _, i := range m.items {
		if i.Status == filter.Status {
			out = append(out, i)
		}
	}
	return out, nil
}
func (m *inMemoryTaskRepo) RoadmapSummary(context.Context, string) (*domainproject.RoadmapSummary, error) {
	return &domainproject.RoadmapSummary{}, nil
}
func (m *inMemoryTaskRepo) UpdateStatus(_ context.Context, _, taskID string, status domainproject.TaskStatus) error {
	for _, i := range m.items {
		if i.ID.Hex() == taskID {
			i.Status = status
			i.UpdatedAt = time.Now().UTC()
		}
	}
	return nil
}
