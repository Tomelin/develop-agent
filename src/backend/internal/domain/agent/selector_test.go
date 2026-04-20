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

func (f *fakeAudit) LogTriadSelection(context.Context, Skill, TriadSelection, bool) error {
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
	svc := NewSelectorService(repo, audit, true, nil).WithSeed(42)

	s, err := svc.SelectTriadDetailed(context.Background(), SkillArchitecture, "Fase 3")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !audit.called {
		t.Fatal("expected audit to be called")
	}
	providers := map[Provider]struct{}{
		s.Producer.Provider: {},
		s.Reviewer.Provider: {},
		s.Refiner.Provider:  {},
	}
	if len(providers) != 3 {
		t.Fatalf("expected 3 unique providers, got %d", len(providers))
	}
}

func TestSelectorComDoisProvidersSemRepetirAgenteQuandoPossivel(t *testing.T) {
	repo := fakeRepo{items: []*Agent{
		mkAgent("a", ProviderOpenAI),
		mkAgent("b", ProviderAnthropic),
		mkAgent("c", ProviderOpenAI),
	}}
	svc := NewSelectorService(repo, nil, true, nil).WithSeed(7)
	s, err := svc.SelectTriadDetailed(context.Background(), SkillArchitecture, "")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if s.Producer.Name == s.Reviewer.Name || s.Producer.Name == s.Refiner.Name || s.Reviewer.Name == s.Refiner.Name {
		t.Fatal("expected distinct agents when enough candidates exist")
	}
}

func TestSelectorComUmAgente(t *testing.T) {
	repo := fakeRepo{items: []*Agent{mkAgent("a", ProviderOpenAI)}}
	svc := NewSelectorService(repo, nil, true, nil)
	s, err := svc.SelectTriadDetailed(context.Background(), SkillArchitecture, "")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if s.Producer.ID != s.Reviewer.ID || s.Reviewer.ID != s.Refiner.ID {
		t.Fatal("expected same agent in all positions")
	}
	if len(s.Warnings) == 0 {
		t.Fatal("expected warning when only one agent is available")
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

func TestSelectorDeterministicoComSeedFixo(t *testing.T) {
	agents := []*Agent{
		mkAgent("a", ProviderOpenAI),
		mkAgent("b", ProviderAnthropic),
		mkAgent("c", ProviderGoogle),
		mkAgent("d", ProviderOllama),
	}
	svcA := NewSelectorService(fakeRepo{items: agents}, nil, true, nil).WithSeed(99)
	svcB := NewSelectorService(fakeRepo{items: agents}, nil, true, nil).WithSeed(99)

	selA, err := svcA.SelectTriadDetailed(context.Background(), SkillArchitecture, "")
	if err != nil {
		t.Fatalf("unexpected err A: %v", err)
	}
	selB, err := svcB.SelectTriadDetailed(context.Background(), SkillArchitecture, "")
	if err != nil {
		t.Fatalf("unexpected err B: %v", err)
	}

	if selA.Producer.Name != selB.Producer.Name || selA.Reviewer.Name != selB.Reviewer.Name || selA.Refiner.Name != selB.Refiner.Name {
		t.Fatal("expected deterministic selection with fixed seed")
	}
}
