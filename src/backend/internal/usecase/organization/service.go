package organization

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	domainorg "github.com/develop-agent/backend/internal/domain/organization"
	"github.com/develop-agent/backend/internal/domain/user"
)

var (
	ErrForbidden = errors.New("forbidden")
)

type Service struct {
	orgs  domainorg.Repository
	users user.Repository
}

func NewService(orgs domainorg.Repository, users user.Repository) *Service {
	return &Service{orgs: orgs, users: users}
}

func (s *Service) GetOrganization(ctx context.Context, organizationID string) (*domainorg.Organization, error) {
	return s.orgs.FindByID(ctx, organizationID)
}

func (s *Service) ListMembers(ctx context.Context, organizationID string) ([]*user.User, error) {
	return s.users.ListByOrganization(ctx, organizationID)
}

type InviteInput struct {
	OrganizationID   string
	Name             string
	Email            string
	OrganizationRole user.OrganizationRole
}

func (s *Service) InviteMember(ctx context.Context, in InviteInput) (*user.User, error) {
	if !isOrganizationRoleAllowed(in.OrganizationRole) {
		return nil, errors.New("invalid organization role")
	}
	if strings.TrimSpace(in.Email) == "" {
		return nil, errors.New("email is required")
	}
	existing, err := s.users.FindByEmail(ctx, in.Email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if existing != nil {
		if existing.OrganizationID.Hex() != in.OrganizationID {
			return nil, errors.New("user already belongs to another organization")
		}
		existing.OrganizationRole = in.OrganizationRole
		existing.Enabled = true
		if err := s.users.Update(ctx, existing); err != nil {
			return nil, err
		}
		return existing.Sanitize(), nil
	}

	defaultName := strings.TrimSpace(in.Name)
	if defaultName == "" {
		defaultName = strings.Split(strings.ToLower(strings.TrimSpace(in.Email)), "@")[0]
	}
	orgID, err := bson.ObjectIDFromHex(in.OrganizationID)
	if err != nil {
		return nil, errors.New("invalid organization id")
	}
	newUser, err := user.New(defaultName, in.Email, "TempPass#123", user.RoleUser, orgID, in.OrganizationRole)
	if err != nil {
		return nil, err
	}
	newUser.Enabled = false
	if err := s.users.Create(ctx, newUser); err != nil {
		return nil, err
	}
	return newUser.Sanitize(), nil
}

func (s *Service) UpdateMemberRole(ctx context.Context, organizationID, memberID string, role user.OrganizationRole) (*user.User, error) {
	if !isOrganizationRoleAllowed(role) {
		return nil, errors.New("invalid organization role")
	}
	member, err := s.users.FindByID(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if member.OrganizationID.Hex() != organizationID {
		return nil, errors.New("member not found")
	}
	member.OrganizationRole = role
	if err := s.users.Update(ctx, member); err != nil {
		return nil, err
	}
	return member.Sanitize(), nil
}

func (s *Service) RemoveMember(ctx context.Context, organizationID, memberID string) error {
	member, err := s.users.FindByID(ctx, memberID)
	if err != nil {
		return err
	}
	if member.OrganizationID.Hex() != organizationID {
		return errors.New("member not found")
	}
	return s.users.SoftDelete(ctx, memberID)
}

func isOrganizationRoleAllowed(role user.OrganizationRole) bool {
	switch role {
	case user.OrganizationRoleOwner, user.OrganizationRoleAdmin, user.OrganizationRoleMember, user.OrganizationRoleViewer:
		return true
	default:
		return false
	}
}
