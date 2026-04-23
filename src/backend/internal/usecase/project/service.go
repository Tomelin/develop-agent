package project

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"

	domain "github.com/develop-agent/backend/internal/domain/project"
)

type Service struct {
	repo         domain.ProjectRepository
	tasks        domain.TaskRepository
	stateMachine *domain.ProjectStateMachine
	inheritance  *domain.InheritanceService
	publisher    PhasePublisher
}

type PhasePublisher interface {
	Publish(ctx context.Context, routingKey string, body []byte) error
}

type CreateProjectInput struct {
	Name               string
	Description        string
	FlowType           domain.FlowType
	OwnerUserID        string
	OrganizationID     string
	LinkedProjectID    string
	DynamicModeEnabled bool
}

func NewService(repo domain.ProjectRepository, tasks domain.TaskRepository) *Service {
	return &Service{
		repo:         repo,
		tasks:        tasks,
		stateMachine: domain.NewProjectStateMachine(),
		inheritance:  domain.NewInheritanceService(repo),
	}
}

func (s *Service) WithPublisher(publisher PhasePublisher) *Service {
	s.publisher = publisher
	return s
}

func (s *Service) CreateProject(ctx context.Context, in CreateProjectInput) (*domain.Project, error) {
	ownerID, err := bson.ObjectIDFromHex(in.OwnerUserID)
	if err != nil {
		return nil, errors.New("invalid owner user id")
	}
	organizationID, err := bson.ObjectIDFromHex(in.OrganizationID)
	if err != nil {
		return nil, errors.New("invalid organization id")
	}

	var linkedID *bson.ObjectID
	if strings.TrimSpace(in.LinkedProjectID) != "" {
		lid, convErr := bson.ObjectIDFromHex(strings.TrimSpace(in.LinkedProjectID))
		if convErr != nil {
			return nil, errors.New("invalid linked project id")
		}
		linkedID = &lid
	}

	p, err := domain.NewProject(in.Name, in.Description, in.FlowType, ownerID, organizationID, in.DynamicModeEnabled, linkedID)
	if err != nil {
		return nil, err
	}

	if in.FlowType == domain.FlowLandingPage || in.FlowType == domain.FlowMarketing {
		if linkedID == nil {
			return nil, errors.New("linked_project_id is required for LANDING_PAGE or MARKETING")
		}
		ctxSpec, buildErr := s.inheritance.BuildInitialContext(ctx, linkedID.Hex())
		if buildErr != nil {
			return nil, buildErr
		}
		p.SpecMD = ctxSpec
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	if _, err := s.StartPhase(ctx, p.ID.Hex(), in.OwnerUserID, 1); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, p.ID.Hex())
}

func (s *Service) Pause(ctx context.Context, projectID, ownerID string) (*domain.Project, error) {
	p, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if p.OwnerUserID.Hex() != ownerID {
		return nil, errors.New("project not found")
	}
	if err := s.stateMachine.TransitionProjectStatus(p, domain.ProjectPaused, "user action", ownerID); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Resume(ctx context.Context, projectID, ownerID string) (*domain.Project, error) {
	p, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if p.OwnerUserID.Hex() != ownerID {
		return nil, errors.New("project not found")
	}
	if err := s.stateMachine.TransitionProjectStatus(p, domain.ProjectInProgress, "user action", ownerID); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Archive(ctx context.Context, projectID, ownerID string) error {
	return s.repo.Archive(ctx, projectID, ownerID)
}

func (s *Service) StartPhase(ctx context.Context, projectID, ownerID string, phaseNumber int) (*domain.PhaseExecution, error) {
	p, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if p.OwnerUserID.Hex() != ownerID {
		return nil, errors.New("project not found")
	}
	if err := validatePhaseStartPreconditions(p, phaseNumber); err != nil {
		return nil, err
	}
	if p.Status == domain.ProjectDraft || p.Status == domain.ProjectPaused {
		if err := s.stateMachine.TransitionProjectStatus(p, domain.ProjectInProgress, "phase started", ownerID); err != nil {
			return nil, err
		}
	}
	phase, err := findPhaseByNumber(p, phaseNumber)
	if err != nil {
		return nil, err
	}

	if len(phase.Tracks) > 0 {
		for _, track := range []domain.Track{domain.TrackFrontend, domain.TrackBackend} {
			if err := s.stateMachine.TransitionTrackStatus(p, phaseNumber, track, domain.PhaseInProgress, "phase started", ownerID); err != nil {
				return nil, err
			}
			if s.publisher != nil {
				if err := s.publisher.Publish(ctx, fmt.Sprintf("phase.%d.%s", phaseNumber, strings.ToLower(string(track))), []byte(fmt.Sprintf(`{"project_id":"%s","owner_user_id":"%s","phase_number":%d,"track":"%s"}`, projectID, ownerID, phaseNumber, track))); err != nil {
					return nil, err
				}
			}
		}
	} else {
		if err := s.stateMachine.TransitionPhaseStatus(p, phaseNumber, domain.PhaseInProgress, "phase started", ownerID); err != nil {
			return nil, err
		}
		if s.publisher != nil {
			if err := s.publisher.Publish(ctx, fmt.Sprintf("phase.%d.full", phaseNumber), []byte(fmt.Sprintf(`{"project_id":"%s","owner_user_id":"%s","phase_number":%d,"track":"FULL"}`, projectID, ownerID, phaseNumber))); err != nil {
				return nil, err
			}
		}
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return findPhaseByNumber(p, phaseNumber)
}

func (s *Service) ApprovePhaseTrack(ctx context.Context, projectID, ownerID string, phaseNumber int, track domain.Track) (*domain.PhaseExecution, error) {
	p, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if p.OwnerUserID.Hex() != ownerID {
		return nil, errors.New("project not found")
	}
	phase, err := findPhaseByNumber(p, phaseNumber)
	if err != nil {
		return nil, err
	}

	if len(phase.Tracks) > 0 {
		if err := s.stateMachine.TransitionTrackStatus(p, phaseNumber, track, domain.PhaseReview, "track awaiting user approval", ownerID); err != nil {
			return nil, err
		}
		if err := s.stateMachine.TransitionTrackStatus(p, phaseNumber, track, domain.PhaseCompleted, "user approved track", ownerID); err != nil {
			return nil, err
		}
	} else {
		if err := s.stateMachine.TransitionPhaseStatus(p, phaseNumber, domain.PhaseCompleted, "user approved phase", ownerID); err != nil {
			return nil, err
		}
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return findPhaseByNumber(p, phaseNumber)
}

func (s *Service) ApproveRoadmapPhase(ctx context.Context, projectID, ownerID string, roadmapJSON []byte) (*domain.PhaseExecution, *RoadmapIngestResult, error) {
	p, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return nil, nil, err
	}
	if p.OwnerUserID.Hex() != ownerID {
		return nil, nil, errors.New("project not found")
	}
	phase, err := findPhaseByNumber(p, 4)
	if err != nil {
		return nil, nil, err
	}
	if phase.Status != domain.PhaseInProgress && phase.Status != domain.PhaseReview {
		return nil, nil, errors.New("phase 4 must be IN_PROGRESS or REVIEW before approval")
	}

	ingester := NewRoadmapIngester(s.tasks)
	result, canonical, err := ingester.Ingest(ctx, projectID, roadmapJSON)
	if err != nil {
		return nil, nil, err
	}

	if phase.Status == domain.PhaseInProgress {
		if err := s.stateMachine.TransitionPhaseStatus(p, 4, domain.PhaseReview, "roadmap validated", ownerID); err != nil {
			return nil, nil, err
		}
	}
	if err := s.stateMachine.TransitionPhaseStatus(p, 4, domain.PhaseCompleted, "roadmap approved and ingested", ownerID); err != nil {
		return nil, nil, err
	}
	p.RoadmapJSON = canonical

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, nil, err
	}
	updatedPhase, err := findPhaseByNumber(p, 4)
	if err != nil {
		return nil, nil, err
	}
	return updatedPhase, result, nil
}

func (s *Service) GetPhaseTracks(ctx context.Context, projectID, ownerID string, phaseNumber int) ([]domain.TrackExecution, error) {
	p, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if p.OwnerUserID.Hex() != ownerID {
		return nil, errors.New("project not found")
	}
	phase, err := findPhaseByNumber(p, phaseNumber)
	if err != nil {
		return nil, err
	}
	if len(phase.Tracks) == 0 {
		return []domain.TrackExecution{{Track: domain.TrackFull, Status: phase.Status, StartedAt: phase.StartedAt, CompletedAt: phase.CompletedAt}}, nil
	}
	return phase.Tracks, nil
}

func validatePhaseStartPreconditions(p *domain.Project, phaseNumber int) error {
	if phaseNumber == 2 {
		phase1, err := findPhaseByNumber(p, 1)
		if err != nil {
			return err
		}
		if phase1.Status != domain.PhaseCompleted {
			return errors.New("phase 1 must be COMPLETED before starting phase 2")
		}
	}
	if phaseNumber == 3 {
		phase2, err := findPhaseByNumber(p, 2)
		if err != nil {
			return err
		}
		if phase2.Status != domain.PhaseCompleted {
			return errors.New("phase 2 must be COMPLETED in both tracks before starting phase 3")
		}
	}
	return nil
}

func findPhaseByNumber(p *domain.Project, phaseNumber int) (*domain.PhaseExecution, error) {
	for i := range p.Phases {
		if p.Phases[i].PhaseNumber == phaseNumber {
			return &p.Phases[i], nil
		}
	}
	return nil, fmt.Errorf("phase %d not found", phaseNumber)
}
