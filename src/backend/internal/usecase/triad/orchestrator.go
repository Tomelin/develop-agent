package triad

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/develop-agent/backend/internal/domain/agent"
	"github.com/develop-agent/backend/pkg/agentsdk"
)

type Stage string

const (
	StageProducer Stage = "producer"
	StageReviewer Stage = "reviewer"
	StageRefiner  Stage = "refiner"
)

type Event struct {
	ExecutionID string         `json:"execution_id"`
	ProjectID   string         `json:"project_id"`
	Stage       Stage          `json:"stage"`
	Status      agent.Status   `json:"status"`
	Message     string         `json:"message"`
	At          time.Time      `json:"at"`
	Meta        map[string]any `json:"meta,omitempty"`
}

type EventSink interface {
	Emit(ctx context.Context, event Event)
}

type ExecutionInput struct {
	ProjectID string
	Prompt    string
	Feedback  string
}

type Orchestrator struct {
	Producer agentsdk.Provider
	Reviewer agentsdk.Provider
	Refiner  agentsdk.Provider
	Events   EventSink
}

func (o *Orchestrator) Run(ctx context.Context, in ExecutionInput) (string, error) {
	execID := uuid.NewString()
	o.emit(ctx, execID, in.ProjectID, StageProducer, agent.StatusRunning, "producer started", nil)
	produced, err := o.complete(ctx, o.Producer, in.Prompt)
	if err != nil {
		o.emit(ctx, execID, in.ProjectID, StageProducer, agent.StatusError, err.Error(), nil)
		return "", err
	}
	o.emit(ctx, execID, in.ProjectID, StageProducer, agent.StatusCompleted, "producer completed", nil)

	o.emit(ctx, execID, in.ProjectID, StageReviewer, agent.StatusRunning, "reviewer started", nil)
	reviewPrompt := fmt.Sprintf("Review this artifact:\n%s", produced)
	review, err := o.complete(ctx, o.Reviewer, reviewPrompt)
	if err != nil {
		o.emit(ctx, execID, in.ProjectID, StageReviewer, agent.StatusError, err.Error(), nil)
		return "", err
	}
	o.emit(ctx, execID, in.ProjectID, StageReviewer, agent.StatusCompleted, "reviewer completed", nil)

	o.emit(ctx, execID, in.ProjectID, StageRefiner, agent.StatusRunning, "refiner started", nil)
	refinePrompt := fmt.Sprintf("Original:\n%s\n\nReview:\n%s\n\nUser Feedback:\n%s", produced, review, in.Feedback)
	refined, err := o.complete(ctx, o.Refiner, refinePrompt)
	if err != nil {
		o.emit(ctx, execID, in.ProjectID, StageRefiner, agent.StatusError, err.Error(), nil)
		return "", err
	}
	o.emit(ctx, execID, in.ProjectID, StageRefiner, agent.StatusCompleted, "refiner completed", map[string]any{"result_preview": refined})

	return refined, nil
}

func (o *Orchestrator) complete(ctx context.Context, provider agentsdk.Provider, content string) (string, error) {
	resp, err := provider.Complete(ctx, agentsdk.CompletionRequest{
		Messages: []agentsdk.Message{{Role: agentsdk.RoleUser, Content: content}},
	})
	if err != nil {
		return "", err
	}
	return resp.Message.Content, nil
}

func (o *Orchestrator) emit(ctx context.Context, executionID, projectID string, stage Stage, status agent.Status, msg string, meta map[string]any) {
	if o.Events == nil {
		return
	}
	o.Events.Emit(ctx, Event{
		ExecutionID: executionID,
		ProjectID:   projectID,
		Stage:       stage,
		Status:      status,
		Message:     msg,
		At:          time.Now().UTC(),
		Meta:        meta,
	})
}
