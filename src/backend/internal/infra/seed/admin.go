package seed

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"

	"github.com/develop-agent/backend/internal/domain/user"
	"github.com/develop-agent/backend/pkg/logger"
)

type AdminSeeder struct {
	users user.Repository
}

func NewAdminSeeder(users user.Repository) *AdminSeeder {
	return &AdminSeeder{users: users}
}

func (s *AdminSeeder) Run(ctx context.Context, forceReset bool) error {
	const adminEmail = "admin@agency.ai"

	u, err := s.users.FindByEmail(ctx, adminEmail)
	if err == nil && u != nil && !forceReset {
		return nil
	}
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) && !forceReset {
		return err
	}
	if u != nil && forceReset {
		_ = s.users.SoftDelete(ctx, u.ID.Hex())
	}

	password, err := randomPassword(12)
	if err != nil {
		return err
	}
	admin, err := user.New("Administrator", adminEmail, password, user.RoleAdmin)
	if err != nil {
		return err
	}
	if err := s.users.Create(ctx, admin); err != nil {
		return err
	}

	logger.Global().Warn("default admin user created",
		zap.String("email", adminEmail),
		zap.String("password", password),
		zap.String("action", "change password immediately"),
	)

	return nil
}

func randomPassword(length int) (string, error) {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf)[:length], nil
}
