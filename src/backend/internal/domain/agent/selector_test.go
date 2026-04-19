package agent

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type fakeRepo struct {
	items []*Agent
	err   error
}

func (f fakeRepo) Create(context.Context, *Agent) error                 { return nil }
func (f fakeRepo) FindByID(context.Context, string) (*Agent, error)     { return nil, nil }
func (f fakeRepo) FindByName(context.Context, string) (*Agent, error)   { return nil, nil }
func (f fakeRepo) Update(context.Context, *Agent) error                 { return nil }
func (f fakeRepo) Delete(context.Context, string) error                 { return nil }
func (f fakeRepo) List(context.Context, ListFilter) ([]*Agent, error)   { return f.items, f.err }
func (f fakeRepo) FindBySkill(context.Context, Skill) ([]*Agent, error) { return f.items, f.err }

type fakeAudit struct{ called bool }

func (f *fakeAudit) LogTriadSelection(context.Context, Skill, Triad, bool) error {
	f.called = true
	return nil
}

func mkAgent(name string, provider Provider) *Agent {
	return &Agent{ID: bson.NewObjectID(), Name: name, Provider: provider, Enabled: true}
}

func TestSelectorDiversidadeProviders(t *testing.T) {
	repo := fakeRepo{items: []*Agent{
		mkAgent("a", ProviderOpenAI),
		mkAgent("b", ProviderAnthropic),
		mkAgent("c", ProviderGoogle),
		mkAgent("d", ProviderOllama),
	}}
	audit := &fakeAudit{}
	svc := NewSelectorService(repo, audit, true, nil)

	triad, err := svc.SelectTriad(context.Background(), SkillArchitecture)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !audit.called {
		t.Fatal("expected audit to be called")
	}
	providers := map[Provider]struct{}{
		triad.Producer.Provider: {},
		triad.Reviewer.Provider: {},
		triad.Refiner.Provider:  {},
	}
	if len(providers) != 3 {
		t.Fatalf("expected 3 unique providers, got %d", len(providers))
	}
}

func TestSelectorComDoisProviders(t *testing.T) {
	repo := fakeRepo{items: []*Agent{
		mkAgent("a", ProviderOpenAI),
		mkAgent("b", ProviderAnthropic),
	}}
	svc := NewSelectorService(repo, nil, true, nil)
	triad, err := svc.SelectTriad(context.Background(), SkillArchitecture)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	providers := map[Provider]struct{}{
		triad.Producer.Provider: {},
		triad.Reviewer.Provider: {},
		triad.Refiner.Provider:  {},
	}
	if len(providers) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(providers))
	}
}

func TestSelectorComUmAgente(t *testing.T) {
	repo := fakeRepo{items: []*Agent{mkAgent("a", ProviderOpenAI)}}
	svc := NewSelectorService(repo, nil, true, nil)
	triad, err := svc.SelectTriad(context.Background(), SkillArchitecture)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if triad.Producer.ID != triad.Reviewer.ID || triad.Reviewer.ID != triad.Refiner.ID {
		t.Fatal("expected same agent in all positions")
	}
}

func TestSelectorSemAgentes(t *testing.T) {
	repo := fakeRepo{items: []*Agent{}}
	svc := NewSelectorService(repo, nil, true, nil)
	_, err := svc.SelectTriad(context.Background(), SkillArchitecture)
	if err == nil {
		t.Fatal("expected error when no agents available")
	}
}

func TestSelectorModoFixo(t *testing.T) {
	fixed := map[Skill]Triad{
		SkillArchitecture: {
			Producer: *mkAgent("p", ProviderOpenAI),
			Reviewer: *mkAgent("r", ProviderAnthropic),
			Refiner:  *mkAgent("f", ProviderGoogle),
		},
	}
	svc := NewSelectorService(fakeRepo{}, nil, false, fixed)
	triad, err := svc.SelectTriad(context.Background(), SkillArchitecture)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if triad.Producer.Name != "p" || triad.Reviewer.Name != "r" || triad.Refiner.Name != "f" {
		t.Fatal("expected fixed triad to be returned")
	}
}
