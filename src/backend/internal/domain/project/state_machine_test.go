package project

import (
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestProjectStateMachine_ProjectTransitions(t *testing.T) {
	sm := NewProjectStateMachine()
	p, err := NewProject("Projeto", "desc", FlowSoftware, bson.NewObjectID(), false, nil)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if err := sm.TransitionProjectStatus(p, ProjectInProgress, "start", "u1"); err != nil {
		t.Fatalf("expected valid transition, got %v", err)
	}
	if err := sm.TransitionProjectStatus(p, ProjectCompleted, "finish", "u1"); err != nil {
		t.Fatalf("expected valid transition, got %v", err)
	}
	if len(p.TransitionHistory) != 2 {
		t.Fatalf("expected 2 transitions, got %d", len(p.TransitionHistory))
	}
}

func TestProjectStateMachine_InvalidTransition(t *testing.T) {
	sm := NewProjectStateMachine()
	p, _ := NewProject("Projeto", "desc", FlowSoftware, bson.NewObjectID(), false, nil)

	if err := sm.TransitionProjectStatus(p, ProjectCompleted, "invalid", "u1"); err == nil {
		t.Fatal("expected error for invalid transition")
	}
}

func TestProjectStateMachine_PhaseTransitions(t *testing.T) {
	sm := NewProjectStateMachine()
	p, _ := NewProject("Projeto", "desc", FlowSoftware, bson.NewObjectID(), false, nil)

	if err := sm.TransitionPhaseStatus(p, 1, PhaseInProgress, "start", "u1"); err != nil {
		t.Fatalf("expected phase start, got %v", err)
	}
	if err := sm.TransitionPhaseStatus(p, 1, PhaseReview, "triad done", "u1"); err != nil {
		t.Fatalf("expected phase review, got %v", err)
	}
	if err := sm.TransitionPhaseStatus(p, 1, PhaseCompleted, "approved", "u1"); err != nil {
		t.Fatalf("expected phase completion, got %v", err)
	}
	if p.CurrentPhaseNumber != 2 {
		t.Fatalf("expected current phase 2, got %d", p.CurrentPhaseNumber)
	}
}
