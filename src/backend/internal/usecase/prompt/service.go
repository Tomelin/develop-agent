package prompt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"

	domain "github.com/develop-agent/backend/internal/domain/prompt"
)

type Template struct {
	ID      string       `json:"id"`
	Title   string       `json:"title"`
	Content string       `json:"content"`
	Group   domain.Group `json:"group"`
	Tags    []string     `json:"tags,omitempty"`
}

type Service struct {
	repo       domain.UserPromptRepository
	aggregator *domain.PromptAggregator
	templates  []Template
}

func NewService(repo domain.UserPromptRepository) *Service {
	return &Service{repo: repo, aggregator: domain.NewPromptAggregator(), templates: defaultTemplates()}
}

func (s *Service) Create(ctx context.Context, userID, title, content string, group domain.Group, priority int, enabled bool, tags []string) (*domain.UserPrompt, error) {
	count, err := s.repo.CountByUserAndGroup(ctx, userID, group)
	if err != nil {
		return nil, err
	}
	if count >= domain.MaxPromptsPerGroup {
		return nil, fmt.Errorf("group %s reached maximum of %d prompts", group, domain.MaxPromptsPerGroup)
	}
	p, err := domain.NewUserPrompt(userID, title, content, group, priority, enabled, tags)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Update(ctx context.Context, userID, id, title, content string, group domain.Group, priority int, enabled bool, tags []string) (*domain.UserPrompt, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item.UserID.Hex() != userID {
		return nil, mongo.ErrNoDocuments
	}
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title is required")
	}
	if !group.IsValid() {
		return nil, errors.New("invalid group")
	}
	if err := domain.ValidateContent(strings.TrimSpace(content)); err != nil {
		return nil, err
	}
	item.Title = strings.TrimSpace(title)
	item.Content = strings.TrimSpace(content)
	item.Group = group
	item.Priority = priority
	item.Enabled = enabled
	item.Tags = tags
	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) Preview(ctx context.Context, userID string, group domain.Group, systemPrompts []string, ragContext, phaseInstruction, userInstruction string) ([]domain.Message, int) {
	global, _ := s.repo.FindByUserAndGroup(ctx, userID, domain.GroupGlobal)
	groupPrompts, _ := s.repo.FindByUserAndGroup(ctx, userID, group)
	messages := s.aggregator.Compose(domain.AggregationInput{
		AgentSystemPrompts: systemPrompts,
		GlobalPrompts:      global,
		GroupPrompts:       groupPrompts,
		RAGContext:         ragContext,
		PhaseInstruction:   phaseInstruction,
		UserInstruction:    userInstruction,
	})
	texts := make([]string, 0, len(messages))
	for _, m := range messages {
		texts = append(texts, m.Content)
	}
	return messages, domain.EstimateTokens(texts...)
}

func (s *Service) Templates() []Template {
	return s.templates
}

func (s *Service) CreateFromTemplate(ctx context.Context, userID, templateID string, priority int) (*domain.UserPrompt, error) {
	for _, t := range s.templates {
		if t.ID == templateID {
			return s.Create(ctx, userID, t.Title, t.Content, t.Group, priority, true, t.Tags)
		}
	}
	return nil, errors.New("template not found")
}

type ExportPayload struct {
	Version string                                `json:"version"`
	Groups  map[domain.Group][]*domain.UserPrompt `json:"groups"`
}

func (s *Service) Export(ctx context.Context, userID string) ([]byte, error) {
	items, err := s.repo.FindAllByUser(ctx, domain.ListFilter{UserID: userID})
	if err != nil {
		return nil, err
	}
	payload := ExportPayload{Version: "1.0", Groups: map[domain.Group][]*domain.UserPrompt{}}
	for _, g := range domain.AllGroups() {
		payload.Groups[g] = []*domain.UserPrompt{}
	}
	for _, item := range items {
		payload.Groups[item.Group] = append(payload.Groups[item.Group], item)
	}
	return json.MarshalIndent(payload, "", "  ")
}

func (s *Service) Import(ctx context.Context, userID string, raw []byte, replace bool) (int, error) {
	var payload ExportPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return 0, errors.New("invalid import payload")
	}
	if replace {
		existing, err := s.repo.FindAllByUser(ctx, domain.ListFilter{UserID: userID})
		if err != nil {
			return 0, err
		}
		for _, item := range existing {
			if err := s.repo.Delete(ctx, userID, item.ID.Hex()); err != nil {
				return 0, err
			}
		}
	}
	created := 0
	for group, items := range payload.Groups {
		if !group.IsValid() {
			continue
		}
		for _, item := range items {
			if _, err := s.Create(ctx, userID, item.Title, item.Content, group, item.Priority, item.Enabled, item.Tags); err == nil {
				created++
			}
		}
	}
	return created, nil
}

func defaultTemplates() []Template {
	return []Template{
		{ID: "backend-golang-gin", Title: "Stack Backend Golang+GIN", Group: domain.GroupGlobal, Content: "Use Go 1.24+, Gin framework, clean architecture and strong typing for all backend implementations.", Tags: []string{"backend", "golang", "gin"}},
		{ID: "frontend-react-tailwind", Title: "Stack Frontend React+TailwindCSS", Group: domain.GroupGlobal, Content: "Prefer React with TypeScript and TailwindCSS. Build responsive, accessible UI components with reusable patterns.", Tags: []string{"frontend", "react", "tailwind"}},
		{ID: "dark-mode-design", Title: "Dark Mode Design", Group: domain.GroupGlobal, Content: "Default to dark mode visual language with high contrast typography and consistent spacing rhythm.", Tags: []string{"design", "dark-mode"}},
		{ID: "postgresql-default", Title: "PostgreSQL Database", Group: domain.GroupEngineering, Content: "Use PostgreSQL as primary datastore. Favor normalized schema, proper indexes and migrations.", Tags: []string{"database", "postgres"}},
		{ID: "clean-architecture", Title: "Clean Architecture", Group: domain.GroupArchitecture, Content: "Apply clean architecture boundaries: domain isolated from infrastructure, explicit interfaces and dependency inversion.", Tags: []string{"architecture"}},
		{ID: "solid-principles", Title: "SOLID Principles", Group: domain.GroupDevelopment, Content: "During implementation, enforce SOLID principles, cohesive modules and readability-first code.", Tags: []string{"solid", "development"}},
		{ID: "tdd-first", Title: "TDD First", Group: domain.GroupTesting, Content: "Write tests before implementation whenever possible. Prioritize deterministic unit tests and clear assertions.", Tags: []string{"testing", "tdd"}},
		{ID: "owasp-priority", Title: "OWASP Security Priority", Group: domain.GroupSecurity, Content: "Prioritize OWASP Top 10 mitigations, strict input validation and principle of least privilege.", Tags: []string{"security", "owasp"}},
		{ID: "restful-api", Title: "RESTful API Design", Group: domain.GroupDevelopment, Content: "Design REST APIs with resource-oriented routes, proper status codes and consistent error payloads.", Tags: []string{"api", "rest"}},
		{ID: "graphql-api", Title: "GraphQL API Design", Group: domain.GroupDevelopment, Content: "When GraphQL is requested, create strongly-typed schema, pagination strategy and resolver-level authorization.", Tags: []string{"api", "graphql"}},
	}
}
