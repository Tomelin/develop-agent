package auth

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/develop-agent/backend/internal/domain/user"
	pkgauth "github.com/develop-agent/backend/pkg/auth"
)

type Service struct {
	users        user.Repository
	tokenManager *pkgauth.TokenManager
	refreshStore pkgauth.RefreshStore
}

type LoginResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	AccessExpiresAt  string `json:"access_expires_at"`
	RefreshExpiresAt string `json:"refresh_expires_at"`
}

func NewService(users user.Repository, tokenManager *pkgauth.TokenManager, refreshStore pkgauth.RefreshStore) *Service {
	return &Service{users: users, tokenManager: tokenManager, refreshStore: refreshStore}
}

func (s *Service) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !u.Enabled {
		return nil, errors.New("user disabled")
	}
	if err := u.CheckPassword(password); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return s.issueTokens(ctx, u)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	userID, err := s.refreshStore.GetUserID(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if err := s.refreshStore.Delete(ctx, refreshToken); err != nil {
		return nil, err
	}

	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if !u.Enabled {
		return nil, errors.New("user disabled")
	}
	return s.issueTokens(ctx, u)
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	return s.refreshStore.Delete(ctx, refreshToken)
}

func (s *Service) ValidateAccessToken(token string) (*pkgauth.Claims, error) {
	return s.tokenManager.ParseAccessToken(token)
}

func (s *Service) issueTokens(ctx context.Context, u *user.User) (*LoginResponse, error) {
	accessToken, accessExp, err := s.tokenManager.GenerateAccessToken(u.ID.Hex(), u.OrganizationID.Hex(), string(u.OrganizationRole), u.Email, string(u.Role))
	if err != nil {
		return nil, err
	}
	refreshToken, refreshExp, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	if err := s.refreshStore.Save(ctx, refreshToken, u.ID.Hex(), refreshExp); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExp.Format("2006-01-02T15:04:05Z"),
		RefreshExpiresAt: refreshExp.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func ParseObjectID(id string) (bson.ObjectID, error) {
	return bson.ObjectIDFromHex(id)
}
