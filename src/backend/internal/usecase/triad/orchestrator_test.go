package triad

import (
	"context"
	"testing"

	"github.com/develop-agent/backend/pkg/agentsdk"
	"github.com/develop-agent/backend/pkg/agentsdk/mock"
)

type inMemorySink struct{ events []Event }

func (s *inMemorySink) Emit(_ context.Context, e Event) { s.events = append(s.events, e) }

func TestOrchestratorRun(t *testing.T) {
	producer := mock.New(agentsdk.CompletionResponse{Message: agentsdk.Message{Role: agentsdk.RoleAssistant, Content: "draft"}})
	reviewer := mock.New(agentsdk.CompletionResponse{Message: agentsdk.Message{Role: agentsdk.RoleAssistant, Content: "fix naming"}})
	refiner := mock.New(agentsdk.CompletionResponse{Message: agentsdk.Message{Role: agentsdk.RoleAssistant, Content: "final artifact"}})
	sink := &inMemorySink{}

	o := Orchestrator{Producer: producer, Reviewer: reviewer, Refiner: refiner, Events: sink}

	got, err := o.Run(context.Background(), ExecutionInput{ProjectID: "p1", Prompt: "build api", Feedback: "be concise"})
	if err != nil {
		t.Fatalf("run triad: %v", err)
	}
	if got != "final artifact" {
		t.Fatalf("expected final artifact, got %q", got)
	}
	if len(sink.events) != 6 {
		t.Fatalf("expected 6 lifecycle events, got %d", len(sink.events))
	}
}
