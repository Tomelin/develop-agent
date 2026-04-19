package interview

import (
	"context"
	"testing"

	domaininterview "github.com/develop-agent/backend/internal/domain/interview"
	domainproject "github.com/develop-agent/backend/internal/domain/project"
	"github.com/develop-agent/backend/pkg/agentsdk"
	"github.com/develop-agent/backend/pkg/agentsdk/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type memInterviewRepo struct {
	byProject map[string]*domaininterview.InterviewSession
}

func (m *memInterviewRepo) EnsureIndexes(context.Context) error { return nil }
func (m *memInterviewRepo) FindByProjectID(_ context.Context, projectID string) (*domaininterview.InterviewSession, error) {
	s, ok := m.byProject[projectID]
	if !ok {
		return nil, mongo.ErrNoDocuments
	}
	cp := *s
	cp.Messages = append([]domaininterview.SessionMessage(nil), s.Messages...)
	return &cp, nil
}
func (m *memInterviewRepo) UpsertByProjectID(_ context.Context, projectID string, session *domaininterview.InterviewSession) error {
	cp := *session
	cp.Messages = append([]domaininterview.SessionMessage(nil), session.Messages...)
	m.byProject[projectID] = &cp
	return nil
}

type memProjectRepo struct{ p *domainproject.Project }

func (m *memProjectRepo) EnsureIndexes(context.Context) error                  { return nil }
func (m *memProjectRepo) Create(context.Context, *domainproject.Project) error { return nil }
func (m *memProjectRepo) FindByID(context.Context, string) (*domainproject.Project, error) {
	return m.p, nil
}
func (m *memProjectRepo) FindByOwner(context.Context, domainproject.ProjectListFilter) ([]*domainproject.Project, int64, error) {
	return nil, 0, nil
}
func (m *memProjectRepo) FindDashboardByOwner(context.Context, domainproject.ProjectListFilter) ([]*domainproject.Project, int64, error) {
	return nil, 0, nil
}
func (m *memProjectRepo) Update(_ context.Context, p *domainproject.Project) error {
	m.p = p
	return nil
}
func (m *memProjectRepo) Archive(context.Context, string, string) error { return nil }
func (m *memProjectRepo) UpdatePhase(context.Context, string, string, domainproject.PhaseExecution) error {
	return nil
}
func (m *memProjectRepo) UpdateSpecMD(context.Context, string, string, string) error { return nil }

func TestInterviewCompleteFlow(t *testing.T) {
	owner := bson.NewObjectID()
	p, _ := domainproject.NewProject("P", "", domainproject.FlowSoftware, owner, false, nil)
	repo := &memInterviewRepo{byProject: map[string]*domaininterview.InterviewSession{}}
	provider := mock.New(
		agentsdk.CompletionResponse{Message: agentsdk.Message{Role: agentsdk.RoleAssistant, Content: "pergunta 1"}},
		agentsdk.CompletionResponse{Message: agentsdk.Message{Role: agentsdk.RoleAssistant, Content: "pergunta 2"}},
		agentsdk.CompletionResponse{Message: agentsdk.Message{Role: agentsdk.RoleAssistant, Content: "pergunta 3"}},
		agentsdk.CompletionResponse{Message: agentsdk.Message{Role: agentsdk.RoleAssistant, Content: "# VISION.md\n\n## Sumário Executivo\nOK"}},
	)
	svc := NewService(repo, &memProjectRepo{p: p}, provider, NewBroker())
	projectID := p.ID.Hex()
	ownerID := owner.Hex()

	for i := 0; i < 3; i++ {
		_, _, err := svc.StreamMessage(context.Background(), projectID, ownerID, "msg", nil)
		if err != nil {
			t.Fatalf("stream message %d: %v", i, err)
		}
	}

	s, err := svc.Confirm(context.Background(), projectID, ownerID)
	if err != nil {
		t.Fatalf("confirm: %v", err)
	}
	if s.Status != domaininterview.StatusCompleted {
		t.Fatalf("expected completed, got %s", s.Status)
	}
	if s.VisionMD == "" {
		t.Fatal("expected vision document")
	}
}

func TestLimitAndMinimumInteractions(t *testing.T) {
	owner := bson.NewObjectID()
	p, _ := domainproject.NewProject("P", "", domainproject.FlowSoftware, owner, false, nil)
	repo := &memInterviewRepo{byProject: map[string]*domaininterview.InterviewSession{}}
	svc := NewService(repo, &memProjectRepo{p: p}, mock.New(), nil)
	projectID := p.ID.Hex()
	ownerID := owner.Hex()

	session, _ := svc.GetSession(context.Background(), projectID, ownerID)
	session.IterationCount = session.MaxIterations
	_ = repo.UpsertByProjectID(context.Background(), projectID, session)
	if _, _, err := svc.StreamMessage(context.Background(), projectID, ownerID, "oi", nil); err == nil {
		t.Fatal("expected iteration limit error")
	}

	repo.byProject = map[string]*domaininterview.InterviewSession{}
	if _, err := svc.Confirm(context.Background(), projectID, ownerID); err == nil {
		t.Fatal("expected minimum interactions error")
	}
}
