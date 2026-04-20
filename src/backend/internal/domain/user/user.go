package user

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

type OrganizationRole string

const (
	OrganizationRoleOwner  OrganizationRole = "OWNER"
	OrganizationRoleAdmin  OrganizationRole = "ADMIN"
	OrganizationRoleMember OrganizationRole = "MEMBER"
	OrganizationRoleViewer OrganizationRole = "VIEWER"
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type User struct {
	ID               bson.ObjectID     `bson:"_id,omitempty" json:"id"`
	OrganizationID   bson.ObjectID     `bson:"organization_id" json:"organization_id"`
	OrganizationRole OrganizationRole  `bson:"organization_role" json:"organization_role"`
	Name             string            `bson:"name" json:"name"`
	Email            string            `bson:"email" json:"email"`
	PasswordHash     string            `bson:"password_hash" json:"-"`
	Role             Role              `bson:"role" json:"role"`
	Prompts          map[string]string `bson:"prompts,omitempty" json:"prompts,omitempty"`
	Enabled          bool              `bson:"enabled" json:"enabled"`
	CreatedAt        time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time         `bson:"updated_at" json:"updated_at"`
	DeletedAt        *time.Time        `bson:"deleted_at,omitempty" json:"-"`
}

func New(name, email, password string, role Role, organizationID bson.ObjectID, organizationRole OrganizationRole) (*User, error) {
	if err := ValidateName(name); err != nil {
		return nil, err
	}
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}
	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &User{
		ID:               bson.NewObjectID(),
		OrganizationID:   organizationID,
		OrganizationRole: organizationRole,
		Name:             strings.TrimSpace(name),
		Email:            strings.ToLower(strings.TrimSpace(email)),
		PasswordHash:     string(hash),
		Role:             role,
		Prompts:          map[string]string{},
		Enabled:          true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}
	return nil
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

func (u *User) Sanitize() *User {
	cp := *u
	cp.PasswordHash = ""
	return &cp
}
