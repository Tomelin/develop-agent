package project

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	domainproject "github.com/develop-agent/backend/internal/domain/project"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type fakeSecurityRunner struct {
	outputs map[string]string
	errors  map[string]error
}

func (f fakeSecurityRunner) Run(_ context.Context, dir, name string, args ...string) (string, error) {
	key := dir + "|" + name + " " + join(args)
	return f.outputs[key], f.errors[key]
}

func TestPhase7RunAudit(t *testing.T) {
	tmp := t.TempDir()
	backend := filepath.Join(tmp, "src/backend")
	if err := os.MkdirAll(filepath.Join(backend, "api/server"), 0o750); err != nil {
		t.Fatal(err)
	}
	serverFile := `package server
func CORSMiddleware() {
	_ = "Access-Control-Allow-Origin"
	_ = "*"
}`
	if err := os.WriteFile(filepath.Join(backend, "api/server/server.go"), []byte(serverFile), 0o600); err != nil {
		t.Fatal(err)
	}

	gosecOut := `{"Issues":[{"rule_id":"G101","details":"hardcoded credential","severity":"HIGH","file":"main.go","line":"22"}]}`
	govulnOut := `{"osv":"GO-2024-9999","module":"example"}`
	npmJSON, _ := json.Marshal(map[string]any{
		"vulnerabilities": map[string]any{
			"axios": map[string]any{"severity": "high", "fixAvailable": true},
		},
	})

	repo := &memCodeFileRepo{}
	svc := NewPhase7Service(repo, nil)
	svc.runner = fakeSecurityRunner{outputs: map[string]string{
		backend + "|gosec -fmt json ./...":                   gosecOut,
		backend + "|govulncheck -json ./...":                 govulnOut,
		tmp + "|trufflehog filesystem . --json":              "",
		filepath.Join(tmp, "frontend") + "|npm audit --json": string(npmJSON),
	}}

	report, err := svc.RunAudit(context.Background(), bson.NewObjectID().Hex(), bson.NewObjectID().Hex(), Phase7AuditInput{
		BackendDir:     backend,
		FrontendDir:    filepath.Join(tmp, "frontend"),
		ProjectRootDir: tmp,
	})
	if err != nil {
		t.Fatalf("run audit: %v", err)
	}
	if report.Summary.TotalFindings == 0 {
		t.Fatal("expected findings")
	}
	if report.Summary.HighCount == 0 {
		t.Fatal("expected at least one high finding")
	}
	if _, ok := repo.files["reports/SECURITY_AUDIT.md"]; !ok {
		t.Fatal("expected persisted security report")
	}
}

func TestEvaluateSecurityAutoRejection(t *testing.T) {
	findings := []domainproject.SecurityFinding{
		{ID: "SECRET-001", Category: "Secrets", CVSS: 10.0},
	}
	res := evaluateSecurityAutoRejection(findings, 2)
	if !res.Triggered || !res.ReturnedPhase5 {
		t.Fatalf("expected immediate auto rejection: %+v", res)
	}
}
