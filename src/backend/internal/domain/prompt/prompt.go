package prompt

import (
	"context"
	"errors"
	"regexp"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Group string

const (
	GroupGlobal          Group = "GLOBAL"
	GroupProjectCreation Group = "PROJECT_CREATION"
	GroupEngineering     Group = "ENGINEERING"
	GroupArchitecture    Group = "ARCHITECTURE"
	GroupPlanning        Group = "PLANNING"
	GroupDevelopment     Group = "DEVELOPMENT"
	GroupTesting         Group = "TESTING"
	GroupSecurity        Group = "SECURITY"
	GroupDocumentation   Group = "DOCUMENTATION"
	GroupDevOps          Group = "DEVOPS"
	GroupLandingPage     Group = "LANDING_PAGE"
	GroupMarketing       Group = "MARKETING"
)

const (
	MaxPromptsPerGroup = 50
	MaxPromptLength    = 2000
	WarnTokens         = 4000
)

type UserPrompt struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectID `bson:"user_id" json:"user_id"`
	Title     string        `bson:"title" json:"title"`
	Content   string        `bson:"content" json:"content"`
	Group     Group         `bson:"group" json:"group"`
	Priority  int           `bson:"priority" json:"priority"`
	Enabled   bool          `bson:"enabled" json:"enabled"`
	Tags      []string      `bson:"tags,omitempty" json:"tags,omitempty"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

type ListFilter struct {
	UserID  string
	Group   Group
	Enabled *bool
}

type ReorderItem struct {
	ID       string `json:"id"`
	Priority int    `json:"priority"`
}

type UserPromptRepository interface {
	EnsureIndexes(ctx context.Context) error
	Create(ctx context.Context, p *UserPrompt) error
	FindByID(ctx context.Context, id string) (*UserPrompt, error)
	FindByUserAndGroup(ctx context.Context, userID string, group Group) ([]*UserPrompt, error)
	FindAllByUser(ctx context.Context, filter ListFilter) ([]*UserPrompt, error)
	CountByUserAndGroup(ctx context.Context, userID string, group Group) (int64, error)
	Update(ctx context.Context, p *UserPrompt) error
	Delete(ctx context.Context, userID, id string) error
	Reorder(ctx context.Context, userID string, items []ReorderItem) error
}

func NewUserPrompt(userID, title, content string, group Group, priority int, enabled bool, tags []string) (*UserPrompt, error) {
	uid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	if title == "" {
		return nil, errors.New("title is required")
	}
	if !group.IsValid() {
		return nil, errors.New("invalid group")
	}
	if err := ValidateContent(content); err != nil {
		return nil, err
	}
	if priority < 0 {
		return nil, errors.New("priority must be >= 0")
	}
	now := time.Now().UTC()
	return &UserPrompt{
		ID:        bson.NewObjectID(),
		UserID:    uid,
		Title:     title,
		Content:   content,
		Group:     group,
		Priority:  priority,
		Enabled:   enabled,
		Tags:      cleanTags(tags),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (g Group) IsValid() bool {
	switch g {
	case GroupGlobal, GroupProjectCreation, GroupEngineering, GroupArchitecture,
		GroupPlanning, GroupDevelopment, GroupTesting, GroupSecurity,
		GroupDocumentation, GroupDevOps, GroupLandingPage, GroupMarketing:
		return true
	default:
		return false
	}
}

func AllGroups() []Group {
	return []Group{
		GroupGlobal, GroupProjectCreation, GroupEngineering, GroupArchitecture,
		GroupPlanning, GroupDevelopment, GroupTesting, GroupSecurity,
		GroupDocumentation, GroupDevOps, GroupLandingPage, GroupMarketing,
	}
}

var injectionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)ignore\s+all\s+previous\s+instructions`),
	regexp.MustCompile(`(?i)disregard\s+the\s+system\s+prompt`),
	regexp.MustCompile(`(?i)reveal\s+system\s+prompt`),
}

func ValidateContent(content string) error {
	if content == "" {
		return errors.New("content is required")
	}
	if len([]rune(content)) > MaxPromptLength {
		return errors.New("content exceeds 2000 characters")
	}
	for _, p := range injectionPatterns {
		if p.MatchString(content) {
			return errors.New("content contains prohibited instruction pattern")
		}
	}
	return nil
}

func EstimateTokens(texts ...string) int {
	totalWords := 0
	for _, t := range texts {
		totalWords += len(strings.Fields(t))
	}
	if totalWords == 0 {
		return 0
	}
	return int(float64(totalWords) / 0.75)
}

func SortByPriority(items []*UserPrompt) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Priority == items[j].Priority {
			return items[i].CreatedAt.Before(items[j].CreatedAt)
		}
		return items[i].Priority < items[j].Priority
	})
}

func cleanTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	for _, tag := range tags {
		t := strings.TrimSpace(tag)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}
