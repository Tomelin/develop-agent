package project

import (
	"context"
	"testing"
	"time"

	domain "github.com/develop-agent/backend/internal/domain/project"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type fakeProjectRepo struct{ p *domain.Project }

func (f *fakeProjectRepo) EnsureIndexes(ctx context.Context) error             { return nil }
func (f *fakeProjectRepo) Create(ctx context.Context, p *domain.Project) error { return nil }
func (f *fakeProjectRepo) FindByID(ctx context.Context, id string) (*domain.Project, error) {
	return f.p, nil
}
func (f *fakeProjectRepo) FindByOwner(ctx context.Context, filter domain.ProjectListFilter) ([]*domain.Project, int64, error) {
	return nil, 0, nil
}
func (f *fakeProjectRepo) FindDashboardByOwner(ctx context.Context, filter domain.ProjectListFilter) ([]*domain.Project, int64, error) {
	return nil, 0, nil
}
func (f *fakeProjectRepo) ListRecent(ctx context.Context, limit int64) ([]*domain.Project, error) {
	return []*domain.Project{f.p}, nil
}
func (f *fakeProjectRepo) Update(ctx context.Context, p *domain.Project) error          { f.p = p; return nil }
func (f *fakeProjectRepo) Archive(ctx context.Context, id string, ownerID string) error { return nil }
func (f *fakeProjectRepo) UpdatePhase(ctx context.Context, projectID, ownerID string, phase domain.PhaseExecution) error {
	return nil
}
func (f *fakeProjectRepo) UpdateSpecMD(ctx context.Context, projectID, ownerID, specMD string) error {
	return nil
}

type fakeCodeFileRepo struct{ files []*domain.CodeFile }

func (f *fakeCodeFileRepo) EnsureIndexes(ctx context.Context) error { return nil }
func (f *fakeCodeFileRepo) Upsert(ctx context.Context, file *domain.CodeFile) error {
	cp := *file
	f.files = append(f.files, &cp)
	return nil
}
func (f *fakeCodeFileRepo) ListByProject(ctx context.Context, projectID string) ([]*domain.CodeFile, error) {
	return f.files, nil
}
func (f *fakeCodeFileRepo) FindByID(ctx context.Context, projectID, fileID string) (*domain.CodeFile, error) {
	return nil, nil
}
func (f *fakeCodeFileRepo) CountAutoRejections(ctx context.Context, projectID string) (int64, error) {
	return 0, nil
}
func (f *fakeCodeFileRepo) IncrementAutoRejections(ctx context.Context, projectID string) error {
	return nil
}

func TestPhase19ServiceRun(t *testing.T) {
	owner := bson.NewObjectID()
	p := &domain.Project{ID: bson.NewObjectID(), OwnerUserID: owner, UpdatedAt: time.Now().UTC()}
	projects := &fakeProjectRepo{p: p}
	files := &fakeCodeFileRepo{}
	svc := NewPhase19Service(projects, files)

	res, err := svc.Run(context.Background(), p.ID.Hex(), owner.Hex(), domain.Phase19RunInput{IncludeFlowA: true, IncludeFlowB: true, IncludeFlowC: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.QualityReport.TotalE2EScenarios != 20 {
		t.Fatalf("unexpected scenarios: %d", res.QualityReport.TotalE2EScenarios)
	}
	if len(res.Artifacts) != 10 {
		t.Fatalf("unexpected artifacts: %d", len(res.Artifacts))
	}
}

func TestAdminQualityReportBuild(t *testing.T) {
	owner := bson.NewObjectID()
	now := time.Now().UTC()
	started := now.Add(-12 * time.Minute)
	completed := now
	p := &domain.Project{
		ID:           bson.NewObjectID(),
		OwnerUserID:  owner,
		FlowType:     domain.FlowSoftware,
		Status:       domain.ProjectCompleted,
		TotalCostUSD: 55.5,
		Phases:       []domain.PhaseExecution{{PhaseNumber: 5, StartedAt: &started, CompletedAt: &completed}},
	}
	projects := &fakeProjectRepo{p: p}
	files := &fakeCodeFileRepo{files: []*domain.CodeFile{{Path: "artifacts/quality/e2e/flow-a.spec.ts"}}}
	svc := NewAdminQualityReportService(projects, files)

	report, err := svc.Build(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.ProjectSampleSize != 1 || report.TestCoveragePercent != 100 {
		t.Fatalf("unexpected report: %+v", report)
	}
}
