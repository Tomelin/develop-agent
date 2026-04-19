package interview

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Status string

type MessageRole string

const (
	StatusActive               Status = "ACTIVE"
	StatusAwaitingConfirmation Status = "AWAITING_CONFIRMATION"
	StatusCompleted            Status = "COMPLETED"
	StatusAbandoned            Status = "ABANDONED"
)

const (
	RoleUser      MessageRole = "USER"
	RoleAssistant MessageRole = "ASSISTANT"
)

type SessionMessage struct {
	Role      MessageRole `bson:"role" json:"role"`
	Content   string      `bson:"content" json:"content"`
	Timestamp time.Time   `bson:"timestamp" json:"timestamp"`
}

type InterviewSession struct {
	ID             bson.ObjectID    `bson:"_id,omitempty" json:"id"`
	ProjectID      bson.ObjectID    `bson:"project_id" json:"project_id"`
	Messages       []SessionMessage `bson:"messages" json:"messages"`
	Status         Status           `bson:"status" json:"status"`
	IterationCount int              `bson:"iteration_count" json:"iteration_count"`
	MaxIterations  int              `bson:"max_iterations" json:"max_iterations"`
	VisionMD       string           `bson:"vision_md,omitempty" json:"vision_md,omitempty"`
	CompletedAt    *time.Time       `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	CreatedAt      time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time        `bson:"updated_at" json:"updated_at"`
}

func NewSession(projectID bson.ObjectID) *InterviewSession {
	now := time.Now().UTC()
	return &InterviewSession{
		ID:             bson.NewObjectID(),
		ProjectID:      projectID,
		Messages:       make([]SessionMessage, 0, 8),
		Status:         StatusActive,
		IterationCount: 0,
		MaxIterations:  10,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (s *InterviewSession) AddMessage(role MessageRole, content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("message content is required")
	}
	if role != RoleUser && role != RoleAssistant {
		return errors.New("invalid role")
	}
	s.Messages = append(s.Messages, SessionMessage{Role: role, Content: content, Timestamp: time.Now().UTC()})
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *InterviewSession) CanIterate() bool {
	return s.IterationCount < s.MaxIterations
}
