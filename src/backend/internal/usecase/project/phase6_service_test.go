package project

import (
	"context"
	"testing"

	domainproject "github.com/develop-agent/backend/internal/domain/project"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type fakeRunner struct {
	outputs map[string]string
	errors  map[string]error
}

func (f fakeRunner) Run(_ context.Context, _, name string, args ...string) (string, error) {
	key := name + " " + join(args)
	return f.outputs[key], f.errors[key]
}

func join(in []string) string {
	out := ""
	for i, s := range in {
		if i > 0 {
			out += " "
		}
		out += s
	}
	return out
}

func TestParseCoverageOutputs(t *testing.T) {
	goTestOut := "ok  github.com/acme/pkg/user 0.1s coverage: 82.5% of statements\nok  github.com/acme/pkg/auth 0.1s coverage: 70.0% of statements"
	funcOut := "github.com/acme/pkg/user/service.go:CreateUser\t100.0%\ngithub.com/acme/pkg/auth/service.go:Login\t60.0%\ntotal:\t(statements)\t80.0%"

	pkgs := parsePackageCoverage(goTestOut)
	if len(pkgs) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(pkgs))
	}
	if got := parseTotalCoverage(funcOut); got != 80.0 {
		t.Fatalf("expected total 80.0, got %.1f", got)
	}
	fns := parseFunctionCoverage(funcOut)
	if len(fns) != 2 {
		t.Fatalf("expected 2 functions, got %d", len(fns))
	}
}

func TestAnalyzeCoverageAndPersist(t *testing.T) {
	files := &memCodeFileRepo{}
	autoRejectCalled := false
	svc := NewPhase6Service(files, func(_ context.Context, _, _ string, _ domainproject.CatastrophicFailureReport) error {
		autoRejectCalled = true
		return nil
	})
	svc.runner = fakeRunner{outputs: map[string]string{
		"go test ./... -cover -coverprofile=coverage.out": "ok github.com/acme/pkg/a 0.1s coverage: 72.0% of statements",
		"go tool cover -func=coverage.out":                "github.com/acme/pkg/a/a.go:A\t72.0%\ntotal:\t(statements)\t72.0%",
	}}

	projectID := bson.NewObjectID().Hex()
	report, below, err := svc.AnalyzeCoverage(context.Background(), projectID, bson.NewObjectID().Hex(), ".", 80)
	if err != nil {
		t.Fatalf("analyze coverage: %v", err)
	}
	if !below || report.TotalPercent != 72.0 {
		t.Fatalf("expected below threshold with 72.0, got below=%t total=%.1f", below, report.TotalPercent)
	}
	if !autoRejectCalled {
		t.Fatal("expected auto reject callback")
	}
	if len(files.files) == 0 {
		t.Fatal("expected persisted coverage artifact")
	}
}

func TestClassifyFailure(t *testing.T) {
	if kind := classifyFailure("foo_test.go:10: undefined: something\nFAIL\tbuild failed"); kind != domainproject.TestFailureTestImplementation {
		t.Fatalf("expected test implementation failure, got %s", kind)
	}
	if kind := classifyFailure("expected 200 got 500"); kind != domainproject.TestFailureProjectBug {
		t.Fatalf("expected project bug, got %s", kind)
	}
}
