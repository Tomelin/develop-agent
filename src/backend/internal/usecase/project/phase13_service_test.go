package project

import (
	"context"
	"testing"

	domainproject "github.com/develop-agent/backend/internal/domain/project"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestPhase13RunGeneratesArtifactsAndCompletesProject(t *testing.T) {
	owner := bson.NewObjectID()
	p, _ := domainproject.NewProject("Phase13", "", domainproject.FlowSoftware, owner, false, nil)
	repo := &memProjectRepo{project: p}
	files := &memCodeFileRepo{}
	svc := NewPhase13Service(repo, files)

	res, err := svc.Run(context.Background(), p.ID.Hex(), owner.Hex(), domainproject.Phase13RunInput{IncludeDevOps: true})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.ProjectStatus != string(domainproject.ProjectCompleted) {
		t.Fatalf("expected project completed, got %s", res.ProjectStatus)
	}
	if len(res.Artifacts) < 10 {
		t.Fatalf("expected many artifacts, got %d", len(res.Artifacts))
	}
	if _, ok := files.files["docs/api/openapi.yaml"]; !ok {
		t.Fatal("expected OpenAPI file")
	}
	if _, ok := files.files["PROJECT_SUMMARY.md"]; !ok {
		t.Fatal("expected project summary")
	}
}

func TestValidateOpenAPI(t *testing.T) {
	good := "openapi: 3.0.3\npaths: {}\n"
	if err := validateOpenAPI(good); err != nil {
		t.Fatalf("expected valid openapi, got %v", err)
	}
	bad := "info: {}\n"
	if err := validateOpenAPI(bad); err == nil {
		t.Fatal("expected validation error")
	}
}
