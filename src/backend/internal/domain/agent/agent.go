package agent

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Provider string

type Skill string

type Status string

const (
	ProviderOpenAI    Provider = "OPENAI"
	ProviderAnthropic Provider = "ANTHROPIC"
	ProviderGoogle    Provider = "GOOGLE"
	ProviderOllama    Provider = "OLLAMA"
)

const (
	SkillProjectCreation     Skill = "PROJECT_CREATION"
	SkillEngineering         Skill = "ENGINEERING"
	SkillArchitecture        Skill = "ARCHITECTURE"
	SkillPlanning            Skill = "PLANNING"
	SkillDevelopmentFrontend Skill = "DEVELOPMENT_FRONTEND"
	SkillDevelopmentBackend  Skill = "DEVELOPMENT_BACKEND"
	SkillTesting             Skill = "TESTING"
	SkillSecurity            Skill = "SECURITY"
	SkillDocumentation       Skill = "DOCUMENTATION"
	SkillDevOps              Skill = "DEVOPS"
	SkillLandingPage         Skill = "LANDING_PAGE"
	SkillMarketing           Skill = "MARKETING"
)

const (
	StatusIdle      Status = "IDLE"
	StatusRunning   Status = "RUNNING"
	StatusPaused    Status = "PAUSED"
	StatusQueued    Status = "QUEUED"
	StatusError     Status = "ERROR"
	StatusCompleted Status = "COMPLETED"
)

type Agent struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string        `bson:"name" json:"name"`
	Description   string        `bson:"description" json:"description"`
	Provider      Provider      `bson:"provider" json:"provider"`
	Model         string        `bson:"model" json:"model"`
	SystemPrompts []string      `bson:"system_prompts" json:"system_prompts"`
	Skills        []Skill       `bson:"skills" json:"skills"`
	Enabled       bool          `bson:"enabled" json:"enabled"`
	ApiKeyRef     string        `bson:"api_key_ref" json:"api_key_ref,omitempty"`
	Status        Status        `bson:"status" json:"status"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time    `bson:"deleted_at,omitempty" json:"-"`
}

type ListFilter struct {
	Enabled  *bool
	Skill    Skill
	Provider Provider
}

type Repository interface {
	Create(ctx context.Context, a *Agent) error
	FindByID(ctx context.Context, id string) (*Agent, error)
	FindByName(ctx context.Context, name string) (*Agent, error)
	Update(ctx context.Context, a *Agent) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter ListFilter) ([]*Agent, error)
	FindBySkill(ctx context.Context, skill Skill) ([]*Agent, error)
}

func New(name, description string, provider Provider, model string, systemPrompts []string, skills []Skill, enabled bool, apiKeyRef string) (*Agent, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	if !provider.IsValid() {
		return nil, errors.New("invalid provider")
	}
	if strings.TrimSpace(model) == "" {
		return nil, errors.New("model is required")
	}
	if len(skills) == 0 {
		return nil, errors.New("at least one skill is required")
	}
	for _, s := range skills {
		if !s.IsValid() {
			return nil, errors.New("invalid skill")
		}
	}

	now := time.Now().UTC()
	return &Agent{
		ID:            bson.NewObjectID(),
		Name:          name,
		Description:   strings.TrimSpace(description),
		Provider:      provider,
		Model:         strings.TrimSpace(model),
		SystemPrompts: cleanPrompts(systemPrompts),
		Skills:        skills,
		Enabled:       enabled,
		ApiKeyRef:     strings.TrimSpace(apiKeyRef),
		Status:        StatusIdle,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func cleanPrompts(prompts []string) []string {
	out := make([]string, 0, len(prompts))
	for _, p := range prompts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func (p Provider) IsValid() bool {
	switch p {
	case ProviderOpenAI, ProviderAnthropic, ProviderGoogle, ProviderOllama:
		return true
	default:
		return false
	}
}

func (s Skill) IsValid() bool {
	switch s {
	case SkillProjectCreation, SkillEngineering, SkillArchitecture, SkillPlanning,
		SkillDevelopmentFrontend, SkillDevelopmentBackend, SkillTesting, SkillSecurity,
		SkillDocumentation, SkillDevOps, SkillLandingPage, SkillMarketing:
		return true
	default:
		return false
	}
}

func (s Status) IsValid() bool {
	switch s {
	case StatusIdle, StatusRunning, StatusPaused, StatusQueued, StatusError, StatusCompleted:
		return true
	default:
		return false
	}
}
