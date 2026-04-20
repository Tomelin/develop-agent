package project

import (
	"archive/zip"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	domainproject "github.com/develop-agent/backend/internal/domain/project"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestPhase15RunWithManualBrief(t *testing.T) {
	owner := bson.NewObjectID()
	p, _ := domainproject.NewProject("Marketing", "", domainproject.FlowMarketing, owner, false, nil)
	repo := &memProjectRepo{project: p}
	files := &memCodeFileRepo{}
	svc := NewPhase15Service(repo, files)

	res, err := svc.Run(context.Background(), p.ID.Hex(), owner.Hex(), domainproject.Phase15RunInput{
		ManualBrief: domainproject.MarketingManualBrief{
			ProductName:    "GrowthOS",
			ProblemSolved:  "Pipeline inconsistente",
			TargetAudience: "Heads de marketing B2B",
			MainBenefits:   []string{"Mais SQLs", "CAC previsível"},
		},
		Channels: []string{"linkedin", "google-ads"},
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.BriefSource != "manual" {
		t.Fatalf("expected manual source, got %s", res.BriefSource)
	}
	if res.TotalPieces == 0 {
		t.Fatal("expected generated pieces")
	}
	if _, ok := files.files["artifacts/marketing/strategy/MARKETING_STRATEGY.md"]; !ok {
		t.Fatal("expected strategy artifact")
	}
}

func TestPhase15ExportPackFilteredByChannel(t *testing.T) {
	owner := bson.NewObjectID()
	p, _ := domainproject.NewProject("Marketing", "", domainproject.FlowMarketing, owner, false, nil)
	repo := &memProjectRepo{project: p}
	files := &memCodeFileRepo{}
	svc := NewPhase15Service(repo, files)
	_, _ = svc.Run(context.Background(), p.ID.Hex(), owner.Hex(), domainproject.Phase15RunInput{
		ManualBrief: domainproject.MarketingManualBrief{ProductName: "GrowthOS", ProblemSolved: "x", TargetAudience: "y", MainBenefits: []string{"z"}},
	})

	raw, _, pieces, err := svc.ExportPack(context.Background(), p.ID.Hex(), owner.Hex(), []string{"linkedin"})
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	if pieces == 0 {
		t.Fatal("expected piece counter")
	}
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		t.Fatalf("zip reader: %v", err)
	}
	for _, f := range zr.File {
		if bytes.Contains([]byte(f.Name), []byte("instagram")) {
			t.Fatalf("unexpected instagram file in filtered export: %s", f.Name)
		}
	}
}

func TestPhase15ConfigureWebhook(t *testing.T) {
	owner := bson.NewObjectID()
	p, _ := domainproject.NewProject("Marketing", "", domainproject.FlowMarketing, owner, false, nil)
	repo := &memProjectRepo{project: p}
	files := &memCodeFileRepo{}
	svc := NewPhase15Service(repo, files)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	res, err := svc.ConfigureWebhook(context.Background(), p.ID.Hex(), owner.Hex(), domainproject.MarketingWebhookInput{URL: ts.URL})
	if err != nil {
		t.Fatalf("configure webhook: %v", err)
	}
	if res.LastTest.Status != "success" {
		t.Fatalf("expected success status, got %s", res.LastTest.Status)
	}
	if _, ok := files.files["artifacts/marketing/webhooks/config.json"]; !ok {
		t.Fatal("expected webhook config artifact")
	}
}
