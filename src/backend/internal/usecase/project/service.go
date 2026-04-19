package project

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"

	domain "github.com/develop-agent/backend/internal/domain/project"
)

type Service struct {
	repo         domain.ProjectRepository
	tasks        domain.TaskRepository
	stateMachine *domain.ProjectStateMachine
	inheritance  *domain.InheritanceService
}

type CreateProjectInput struct {
	Name               string
	Description        string
	FlowType           domain.FlowType
	OwnerUserID        string
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

func (s *Service) CreateProject(ctx context.Context, in CreateProjectInput) (*domain.Project, error) {
	ownerID, err := bson.ObjectIDFromHex(in.OwnerUserID)
	if err != nil {
		return nil, errors.New("invalid owner user id")
	}

	var linkedID *bson.ObjectID
	if strings.TrimSpace(in.LinkedProjectID) != "" {
		lid, convErr := bson.ObjectIDFromHex(strings.TrimSpace(in.LinkedProjectID))
		if convErr != nil {
			return nil, errors.New("invalid linked project id")
		}
		linkedID = &lid
	}

	p, err := domain.NewProject(in.Name, in.Description, in.FlowType, ownerID, in.DynamicModeEnabled, linkedID)
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
	return p, nil
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
