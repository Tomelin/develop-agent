package project

import (
	"errors"
	"fmt"
	"time"
)

type ProjectStateMachine struct{}

func NewProjectStateMachine() *ProjectStateMachine {
	return &ProjectStateMachine{}
}

func (sm *ProjectStateMachine) TransitionProjectStatus(p *Project, to ProjectStatus, reason, triggeredBy string) error {
	if p == nil {
		return errors.New("project is nil")
	}
	from := p.Status
	if !sm.isValidProjectTransition(from, to) {
		return fmt.Errorf("invalid project transition: %s -> %s", from, to)
	}

	p.Status = to
	now := time.Now().UTC()
	if to == ProjectArchived {
		p.ArchivedAt = &now
	}
	p.UpdatedAt = now
	p.TransitionHistory = append(p.TransitionHistory, TransitionRecord{
		Kind:      "project",
		From:      string(from),
		To:        string(to),
		Reason:    reason,
		Triggered: triggeredBy,
		At:        now,
	})
	return nil
}

func (sm *ProjectStateMachine) TransitionPhaseStatus(p *Project, phaseNumber int, to PhaseStatus, reason, triggeredBy string) error {
	if p == nil {
		return errors.New("project is nil")
	}
	phase, err := findPhase(p, phaseNumber)
	if err != nil {
		return err
	}
	from := phase.Status
	if !sm.isValidPhaseTransition(from, to) {
		return fmt.Errorf("invalid phase transition: %s -> %s", from, to)
	}

	now := time.Now().UTC()
	phase.Status = to
	if to == PhaseInProgress && phase.StartedAt == nil {
		phase.StartedAt = &now
	}
	if to == PhaseCompleted {
		phase.CompletedAt = &now
		if phaseNumber >= p.CurrentPhaseNumber {
			p.CurrentPhaseNumber = phaseNumber + 1
		}
	}
	p.UpdatedAt = now
	p.TransitionHistory = append(p.TransitionHistory, TransitionRecord{
		Kind:      "phase",
		From:      string(from),
		To:        string(to),
		Reason:    reason,
		Triggered: triggeredBy,
		At:        now,
		Phase:     phaseNumber,
	})
	return nil
}

func (sm *ProjectStateMachine) isValidProjectTransition(from, to ProjectStatus) bool {
	switch from {
	case ProjectDraft:
		return to == ProjectInProgress
	case ProjectInProgress:
		return to == ProjectPaused || to == ProjectCompleted || to == ProjectArchived
	case ProjectPaused:
		return to == ProjectInProgress
	default:
		return false
	}
}

func (sm *ProjectStateMachine) isValidPhaseTransition(from, to PhaseStatus) bool {
	switch from {
	case PhasePending:
		return to == PhaseInProgress
	case PhaseInProgress:
		return to == PhaseReview || to == PhaseRejected
	case PhaseReview:
		return to == PhaseInProgress || to == PhaseCompleted
	case PhaseRejected:
		return to == PhaseInProgress
	default:
		return false
	}
}

func findPhase(p *Project, phaseNumber int) (*PhaseExecution, error) {
	for i := range p.Phases {
		if p.Phases[i].PhaseNumber == phaseNumber {
			return &p.Phases[i], nil
		}
	}
	return nil, fmt.Errorf("phase %d not found", phaseNumber)
}
