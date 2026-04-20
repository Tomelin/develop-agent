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

func (sm *ProjectStateMachine) TransitionTrackStatus(p *Project, phaseNumber int, track Track, to PhaseStatus, reason, triggeredBy string) error {
	if p == nil {
		return errors.New("project is nil")
	}
	phase, err := findPhase(p, phaseNumber)
	if err != nil {
		return err
	}
	targetTrack, err := findTrack(phase, track)
	if err != nil {
		return err
	}
	from := targetTrack.Status
	if !sm.isValidPhaseTransition(from, to) {
		return fmt.Errorf("invalid track transition: %s -> %s", from, to)
	}

	now := time.Now().UTC()
	targetTrack.Status = to
	if to == PhaseInProgress && targetTrack.StartedAt == nil {
		targetTrack.StartedAt = &now
	}
	if to == PhaseCompleted {
		targetTrack.CompletedAt = &now
	}

	phase.Status = derivePhaseStatusFromTracks(phase.Tracks)
	if phase.Status == PhaseCompleted {
		phase.CompletedAt = &now
		if phaseNumber >= p.CurrentPhaseNumber {
			p.CurrentPhaseNumber = phaseNumber + 1
		}
	}
	if phase.Status == PhaseInProgress && phase.StartedAt == nil {
		phase.StartedAt = &now
	}

	p.UpdatedAt = now
	p.TransitionHistory = append(p.TransitionHistory, TransitionRecord{
		Kind:      "phase_track",
		From:      string(from),
		To:        string(to),
		Reason:    reason,
		Triggered: triggeredBy,
		At:        now,
		Phase:     phaseNumber,
		Meta: map[string]any{
			"track": track,
		},
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

func findTrack(phase *PhaseExecution, track Track) (*TrackExecution, error) {
	for i := range phase.Tracks {
		if phase.Tracks[i].Track == track {
			return &phase.Tracks[i], nil
		}
	}
	return nil, fmt.Errorf("track %s not found", track)
}

func derivePhaseStatusFromTracks(tracks []TrackExecution) PhaseStatus {
	if len(tracks) == 0 {
		return PhasePending
	}
	allCompleted := true
	hasInProgress := false
	hasReview := false
	for _, track := range tracks {
		if track.Status != PhaseCompleted {
			allCompleted = false
		}
		if track.Status == PhaseInProgress {
			hasInProgress = true
		}
		if track.Status == PhaseReview {
			hasReview = true
		}
	}
	if allCompleted {
		return PhaseCompleted
	}
	if hasInProgress {
		return PhaseInProgress
	}
	if hasReview {
		return PhaseReview
	}
	return PhasePending
}
