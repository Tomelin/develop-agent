package project

import (
	"context"
	"testing"

	domainproject "github.com/develop-agent/backend/internal/domain/project"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestPhase14RunWithManualBrief(t *testing.T) {
	owner := bson.NewObjectID()
	p, _ := domainproject.NewProject("LP", "", domainproject.FlowLandingPage, owner, false, nil)
	repo := &memProjectRepo{project: p}
	files := &memCodeFileRepo{}
	svc := NewPhase14Service(repo, files)

	res, err := svc.Run(context.Background(), p.ID.Hex(), owner.Hex(), domainproject.Phase14RunInput{
		ManualBrief: domainproject.LandingPageManualBrief{
			ProductName:         "ConversorX",
			ProblemSolved:       "Baixa conversão",
			TargetAudience:      "Times de marketing",
			UniqueValueProposed: "Mais leads em 7 dias",
			KeyFeatures:         []string{"Heatmap", "A/B Test", "Insights"},
		},
		GenerateVariants: true,
		VariantCount:     2,
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.BriefSource != "manual" {
		t.Fatalf("expected manual source, got %s", res.BriefSource)
	}
	if len(res.ArtifactPaths) < 8 {
		t.Fatalf("expected artifacts, got %d", len(res.ArtifactPaths))
	}
	if len(res.Variants) != 2 {
		t.Fatalf("expected 2 variants, got %d", len(res.Variants))
	}
	if _, ok := files.files["artifacts/landing/CONVERSION_REPORT.md"]; !ok {
		t.Fatal("expected conversion report artifact")
	}
}

func TestPhase14RunExtractsLinkedContext(t *testing.T) {
	owner := bson.NewObjectID()
	base, _ := domainproject.NewProject("Base", "", domainproject.FlowSoftware, owner, false, nil)
	landing, _ := domainproject.NewProject("Landing", "", domainproject.FlowLandingPage, owner, false, &base.ID)
	files := &projectAwareCodeFileRepo{store: map[string][]*domainproject.CodeFile{
		base.ID.Hex(): {
			{ProjectID: base.ID, Path: "VISION.md", Content: "- nome do produto: Produto Alfa\n- problema: onboarding lento\n- público-alvo: PMEs\n- proposta de valor: implantação em 1 dia"},
			{ProjectID: base.ID, Path: "SPEC.md", Content: "- Automação\n- Integração Slack\n- Dashboard\n- Alertas"},
			{ProjectID: base.ID, Path: "prompts/user.md", Content: "tema: dark\npaleta: #111111 #F4C430"},
		},
	}}
	repo := &multiProjectRepo{items: map[string]*domainproject.Project{landing.ID.Hex(): landing, base.ID.Hex(): base}}
	svc := NewPhase14Service(repo, files)

	res, err := svc.Run(context.Background(), landing.ID.Hex(), owner.Hex(), domainproject.Phase14RunInput{UseLinkedProject: true})
	if err != nil {
		t.Fatalf("run linked: %v", err)
	}
	if res.BriefSource != "linked_project" {
		t.Fatalf("expected linked source, got %s", res.BriefSource)
	}
	if _, ok := files.byPath(landing.ID.Hex())["artifacts/landing/brief.json"]; !ok {
		t.Fatal("expected brief artifact")
	}
}

type multiProjectRepo struct {
	items map[string]*domainproject.Project
}

func (m *multiProjectRepo) EnsureIndexes(context.Context) error { return nil }
func (m *multiProjectRepo) Create(_ context.Context, p *domainproject.Project) error {
	m.items[p.ID.Hex()] = p
	return nil
}
func (m *multiProjectRepo) FindByID(_ context.Context, id string) (*domainproject.Project, error) {
	return m.items[id], nil
}
func (m *multiProjectRepo) FindByOwner(context.Context, domainproject.ProjectListFilter) ([]*domainproject.Project, int64, error) {
	return nil, 0, nil
}
func (m *multiProjectRepo) FindDashboardByOwner(context.Context, domainproject.ProjectListFilter) ([]*domainproject.Project, int64, error) {
	return nil, 0, nil
}
func (m *multiProjectRepo) Update(_ context.Context, p *domainproject.Project) error {
	m.items[p.ID.Hex()] = p
	return nil
}
func (m *multiProjectRepo) Archive(context.Context, string, string) error { return nil }
func (m *multiProjectRepo) UpdatePhase(context.Context, string, string, domainproject.PhaseExecution) error {
	return nil
}
func (m *multiProjectRepo) UpdateSpecMD(context.Context, string, string, string) error { return nil }

type projectAwareCodeFileRepo struct {
	store map[string][]*domainproject.CodeFile
}

func (m *projectAwareCodeFileRepo) EnsureIndexes(context.Context) error { return nil }
func (m *projectAwareCodeFileRepo) Upsert(_ context.Context, file *domainproject.CodeFile) error {
	pid := file.ProjectID.Hex()
	items := m.store[pid]
	for i := range items {
		if items[i].Path == file.Path {
			items[i] = file
			m.store[pid] = items
			return nil
		}
	}
	m.store[pid] = append(items, file)
	return nil
}
func (m *projectAwareCodeFileRepo) ListByProject(_ context.Context, projectID string) ([]*domainproject.CodeFile, error) {
	return m.store[projectID], nil
}
func (m *projectAwareCodeFileRepo) FindByID(context.Context, string, string) (*domainproject.CodeFile, error) {
	return nil, nil
}
func (m *projectAwareCodeFileRepo) CountAutoRejections(context.Context, string) (int64, error) {
	return 0, nil
}
func (m *projectAwareCodeFileRepo) IncrementAutoRejections(context.Context, string) error { return nil }
func (m *projectAwareCodeFileRepo) byPath(projectID string) map[string]*domainproject.CodeFile {
	out := map[string]*domainproject.CodeFile{}
	for _, f := range m.store[projectID] {
		out[f.Path] = f
	}
	return out
}
