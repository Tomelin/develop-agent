package user

import (
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestValidateEmail(t *testing.T) {
	if err := ValidateEmail("invalid"); err == nil {
		t.Fatal("expected invalid email error")
	}
	if err := ValidateEmail("john@doe.com"); err != nil {
		t.Fatalf("expected valid email, got %v", err)
	}
}

func TestValidatePassword(t *testing.T) {
	if err := ValidatePassword("123"); err == nil {
		t.Fatal("expected weak password error")
	}
	if err := ValidatePassword("12345678"); err != nil {
		t.Fatalf("expected valid password, got %v", err)
	}
}

func TestNewAndCheckPassword(t *testing.T) {
	u, err := New("John", "john@doe.com", "supersecret", RoleUser, bson.NewObjectID(), OrganizationRoleMember)
	if err != nil {
		t.Fatalf("unexpected new user error: %v", err)
	}
	if err := u.CheckPassword("supersecret"); err != nil {
		t.Fatalf("expected valid password check, got %v", err)
	}
}
