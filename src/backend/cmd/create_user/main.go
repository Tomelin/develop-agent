package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/develop-agent/backend/internal/domain/user"
	"github.com/develop-agent/backend/internal/infra/database/mongodb"
)

func main() {
	mongoURI := flag.String("mongo-uri", "mongodb://admin:password@localhost:27017", "MongoDB connection URI")
	dbName := flag.String("db", "develop_agent", "MongoDB database name")
	name := flag.String("name", "Admin", "User name")
	email := flag.String("email", "admin@example.com", "User email")
	password := flag.String("password", "ChangeMe123!", "User password")
	role := flag.String("role", string(user.RoleAdmin), "User global role: ADMIN or USER")
	orgID := flag.String("org-id", "", "Organization ID (hex ObjectID)")
	orgRole := flag.String("org-role", string(user.OrganizationRoleOwner), "User organization role: OWNER, ADMIN, MEMBER or VIEWER")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	adapter, err := mongodb.NewAdapter(*mongoURI)
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer func() {
		if err := adapter.Close(context.Background()); err != nil {
			log.Printf("warning: failed to close mongodb connection: %v", err)
		}
	}()

	repo := mongodb.NewUserRepository(adapter, *dbName)

	parsedRole, err := parseRole(*role)
	if err != nil {
		log.Fatal(err)
	}

	parsedOrgRole, err := parseOrganizationRole(*orgRole)
	if err != nil {
		log.Fatal(err)
	}

	parsedOrgID, err := parseOrganizationID(*orgID)
	if err != nil {
		log.Fatal(err)
	}

	u, err := user.New(*name, *email, *password, parsedRole, parsedOrgID, parsedOrgRole)
	if err != nil {
		log.Fatalf("failed to create user entity: %v", err)
	}

	if err := repo.Create(ctx, u); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Fatalf("user with email %q already exists", strings.ToLower(strings.TrimSpace(*email)))
		}
		log.Fatalf("failed to persist user: %v", err)
	}

	fmt.Printf("user created successfully: id=%s email=%s role=%s organization_role=%s\n", u.ID.Hex(), u.Email, u.Role, u.OrganizationRole)
}

func parseRole(raw string) (user.Role, error) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case string(user.RoleAdmin):
		return user.RoleAdmin, nil
	case string(user.RoleUser):
		return user.RoleUser, nil
	default:
		return "", fmt.Errorf("invalid role %q: expected ADMIN or USER", raw)
	}
}

func parseOrganizationRole(raw string) (user.OrganizationRole, error) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case string(user.OrganizationRoleOwner):
		return user.OrganizationRoleOwner, nil
	case string(user.OrganizationRoleAdmin):
		return user.OrganizationRoleAdmin, nil
	case string(user.OrganizationRoleMember):
		return user.OrganizationRoleMember, nil
	case string(user.OrganizationRoleViewer):
		return user.OrganizationRoleViewer, nil
	default:
		return "", fmt.Errorf("invalid org-role %q: expected OWNER, ADMIN, MEMBER or VIEWER", raw)
	}
}

func parseOrganizationID(raw string) (bson.ObjectID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return bson.NilObjectID, nil
	}

	id, err := bson.ObjectIDFromHex(trimmed)
	if err != nil {
		return bson.NilObjectID, fmt.Errorf("invalid org-id %q: %w", raw, err)
	}

	return id, nil
}
