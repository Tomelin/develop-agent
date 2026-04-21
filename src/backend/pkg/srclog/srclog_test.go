package srclog

import (
	"testing"
)

func TestParseCallerName(t *testing.T) {
	tests := []struct {
		name     string
		fullName string
		expected string
	}{
		{
			name:     "pkg matching",
			fullName: "github.com/develop-agent/backend/pkg/logger.NewLogger",
			expected: "third party",
		},
		{
			name:     "infrastructure matching",
			fullName: "github.com/develop-agent/backend/internal/infrastructure/mongo.(*Repository).Insert",
			expected: "infra - mongo",
		},
		{
			name:     "business matching - usecase",
			fullName: "github.com/develop-agent/backend/internal/usecase.CreateUser",
			expected: "business - usecase",
		},
		{
			name:     "business matching - domain",
			fullName: "github.com/develop-agent/backend/internal/domain/user.(*Entity).Validate",
			expected: "business - user",
		},
		{
			name:     "unknown - external lib",
			fullName: "github.com/gin-gonic/gin.(*Engine).Run",
			expected: "unknown",
		},
		{
			name:     "unknown - main",
			fullName: "main.main",
			expected: "unknown",
		},
		{
			name:     "empty string",
			fullName: "",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCallerName(tt.fullName)
			if got != tt.expected {
				t.Errorf("ParseCallerName(%q) = %v, want %v", tt.fullName, got, tt.expected)
			}
		})
	}
}

// A simple test for GetComponent itself to ensure it doesn't panic
// and returns "third party" when called from this package tests
// because the test is in "pkg/srclog"
func TestGetComponent(t *testing.T) {
	got := GetComponent()
	if got != "third party" {
		t.Errorf("GetComponent() = %v, want %v (since it's called from inside pkg/)", got, "third party")
	}
}
