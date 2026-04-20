package project

import (
	"context"
	"testing"

	domainproject "github.com/develop-agent/backend/internal/domain/project"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type memProjectRepo struct {
	project *domainproject.Project
}

func (m *memProjectRepo) EnsureIndexes(context.Context) error { return nil }
func (m *memProjectRepo) Create(_ context.Context, p *domainproject.Project) error {
	m.project = p
	return nil
}
func (m *memProjectRepo) FindByID(context.Context, string) (*domainproject.Project, error) {
	return m.project, nil
}
func (m *memProjectRepo) FindByOwner(context.Context, domainproject.ProjectListFilter) ([]*domainproject.Project, int64, error) {
	return nil, 0, nil
}
func (m *memProjectRepo) FindDashboardByOwner(context.Context, domainproject.ProjectListFilter) ([]*domainproject.Project, int64, error) {
	return nil, 0, nil
}
func (m *memProjectRepo) ListRecent(context.Context, int64) ([]*domainproject.Project, error) {
	if m.project == nil {
		return []*domainproject.Project{}, nil
	}
	return []*domainproject.Project{m.project}, nil
}
func (m *memProjectRepo) Update(_ context.Context, p *domainproject.Project) error {
	m.project = p
	return nil
}
func (m *memProjectRepo) Archive(context.Context, string, string) error { return nil }
func (m *memProjectRepo) UpdatePhase(context.Context, string, string, domainproject.PhaseExecution) error {
	return nil
}
func (m *memProjectRepo) UpdateSpecMD(context.Context, string, string, string) error { return nil }

type noopTaskRepo struct{}

func (n *noopTaskRepo) EnsureIndexes(context.Context) error { return nil }
func (n *noopTaskRepo) BulkCreate(context.Context, []*domainproject.Task) error {
	return nil
}
func (n *noopTaskRepo) ListByProject(context.Context, domainproject.TaskListFilter) ([]*domainproject.Task, error) {
	return nil, nil
}
func (n *noopTaskRepo) RoadmapSummary(context.Context, string) (*domainproject.RoadmapSummary, error) {
	return &domainproject.RoadmapSummary{}, nil
}
func (n *noopTaskRepo) UpdateStatus(context.Context, string, string, domainproject.TaskStatus) error {
	return nil
}

func TestStartPhase2RequiresPhase1Completed(t *testing.T) {
	owner := bson.NewObjectID()
	p, _ := domainproject.NewProject("P", "", domainproject.FlowSoftware, owner, false, nil)
	svc := NewService(&memProjectRepo{project: p}, &noopTaskRepo{})

	if _, err := svc.StartPhase(context.Background(), p.ID.Hex(), owner.Hex(), 2); err == nil {
		t.Fatal("expected precondition error when phase 1 is not completed")
	}
}

func TestSplitTracksApprovalCompletesPhaseOnlyWhenBothApproved(t *testing.T) {
	owner := bson.NewObjectID()
	p, _ := domainproject.NewProject("P", "", domainproject.FlowSoftware, owner, false, nil)
	p.Phases[0].Status = domainproject.PhaseCompleted // phase 1
	svc := NewService(&memProjectRepo{project: p}, &noopTaskRepo{})

	phase2, err := svc.StartPhase(context.Background(), p.ID.Hex(), owner.Hex(), 2)
	if err != nil {
		t.Fatalf("start phase 2: %v", err)
	}
	if len(phase2.Tracks) != 2 {
		t.Fatalf("expected 2 tracks, got %d", len(phase2.Tracks))
	}

	phase2, err = svc.ApprovePhaseTrack(context.Background(), p.ID.Hex(), owner.Hex(), 2, domainproject.TrackFrontend)
	if err != nil {
		t.Fatalf("approve frontend: %v", err)
	}
	if phase2.Status == domainproject.PhaseCompleted {
		t.Fatal("phase should not complete with only one track approved")
	}

	phase2, err = svc.ApprovePhaseTrack(context.Background(), p.ID.Hex(), owner.Hex(), 2, domainproject.TrackBackend)
	if err != nil {
		t.Fatalf("approve backend: %v", err)
	}
	if phase2.Status != domainproject.PhaseCompleted {
		t.Fatalf("expected phase completed when both tracks approved, got %s", phase2.Status)
	}
}

func TestPhase3RequiresPhase2CompletedInBothTracks(t *testing.T) {
	owner := bson.NewObjectID()
	p, _ := domainproject.NewProject("P", "", domainproject.FlowSoftware, owner, false, nil)
	p.Phases[0].Status = domainproject.PhaseCompleted // phase 1
	svc := NewService(&memProjectRepo{project: p}, &noopTaskRepo{})

	if _, err := svc.StartPhase(context.Background(), p.ID.Hex(), owner.Hex(), 3); err == nil {
		t.Fatal("expected phase 2 completion precondition")
	}
}
