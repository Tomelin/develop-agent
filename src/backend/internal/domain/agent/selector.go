package agent

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Triad struct {
	Producer Agent `json:"producer"`
	Reviewer Agent `json:"reviewer"`
	Refiner  Agent `json:"refiner"`
}

type TriadSelection struct {
	PhaseName          string            `json:"phase_name"`
	Producer           Agent             `json:"producer"`
	Reviewer           Agent             `json:"reviewer"`
	Refiner            Agent             `json:"refiner"`
	CandidateAgents    []Agent           `json:"candidate_agents"`
	SelectionReason    map[string]string `json:"selection_reason"`
	SelectionTimestamp time.Time         `json:"selection_timestamp"`
	Warnings           []string          `json:"warnings,omitempty"`
}

type AuditLogger interface {
	LogTriadSelection(ctx context.Context, skill Skill, selection TriadSelection, dynamic bool) error
}

type Service struct {
	repo          Repository
	audit         AuditLogger
	dynamicModeOn bool
	rand          *rand.Rand
	fixedBySkill  map[Skill]Triad
}

func NewSelectorService(repo Repository, audit AuditLogger, dynamicModeOn bool, fixed map[Skill]Triad) *Service {
	return &Service{
		repo:          repo,
		audit:         audit,
		dynamicModeOn: dynamicModeOn,
		rand:          rand.New(rand.NewSource(time.Now().UnixNano())), // #nosec G404 -- non-cryptographic selection for load balancing
		fixedBySkill:  fixed,
	}
}

func (s *Service) WithSeed(seed int64) *Service {
	s.rand = rand.New(rand.NewSource(seed)) // #nosec G404 -- deterministic pseudo-random for testability
	return s
}

func (s *Service) SelectTriad(ctx context.Context, skill Skill) (Triad, error) {
	selection, err := s.SelectTriadDetailed(ctx, skill, "")
	if err != nil {
		return Triad{}, err
	}
	return Triad{Producer: selection.Producer, Reviewer: selection.Reviewer, Refiner: selection.Refiner}, nil
}

func (s *Service) SelectTriadDetailed(ctx context.Context, skill Skill, phaseName string) (TriadSelection, error) {
	if !skill.IsValid() {
		return TriadSelection{}, errors.New("invalid skill")
	}

	if !s.dynamicModeOn {
		triad, ok := s.fixedBySkill[skill]
		if !ok {
			return TriadSelection{}, fmt.Errorf("no fixed triad configured for skill %s", skill)
		}
		selection := TriadSelection{
			PhaseName:          phaseName,
			Producer:           triad.Producer,
			Reviewer:           triad.Reviewer,
			Refiner:            triad.Refiner,
			CandidateAgents:    []Agent{triad.Producer, triad.Reviewer, triad.Refiner},
			SelectionTimestamp: time.Now().UTC(),
			SelectionReason: map[string]string{
				"producer": "configuração fixa do projeto",
				"reviewer": "configuração fixa do projeto",
				"refiner":  "configuração fixa do projeto",
			},
		}
		if s.audit != nil {
			_ = s.audit.LogTriadSelection(ctx, skill, selection, false)
		}
		return selection, nil
	}

	agents, err := s.repo.FindBySkill(ctx, skill)
	if err != nil {
		return TriadSelection{}, err
	}
	if len(agents) == 0 {
		return TriadSelection{}, fmt.Errorf("no enabled agents found for skill %s", skill)
	}

	selection := s.pickTriadSelection(agents)
	selection.PhaseName = phaseName
	if s.audit != nil {
		_ = s.audit.LogTriadSelection(ctx, skill, selection, true)
	}
	return selection, nil
}

func (s *Service) pickTriadSelection(agents []*Agent) TriadSelection {
	selection := TriadSelection{
		CandidateAgents:    cloneAgents(agents),
		SelectionReason:    map[string]string{},
		SelectionTimestamp: time.Now().UTC(),
	}

	if len(agents) == 1 {
		a := *agents[0]
		selection.Producer = a
		selection.Reviewer = a
		selection.Refiner = a
		selection.SelectionReason["producer"] = "único agente disponível para a skill"
		selection.SelectionReason["reviewer"] = "único agente disponível para a skill"
		selection.SelectionReason["refiner"] = "único agente disponível para a skill"
		selection.Warnings = append(selection.Warnings, "apenas um agente habilitado: mesmo agente usado nos três papéis")
		return selection
	}

	byProvider := map[Provider][]*Agent{}
	for _, a := range agents {
		byProvider[a.Provider] = append(byProvider[a.Provider], a)
	}
	providers := make([]Provider, 0, len(byProvider))
	for p := range byProvider {
		providers = append(providers, p)
	}
	sort.Slice(providers, func(i, j int) bool { return providers[i] < providers[j] })
	s.rand.Shuffle(len(providers), func(i, j int) { providers[i], providers[j] = providers[j], providers[i] })

	selected := make([]*Agent, 0, 3)
	usedIDs := map[string]struct{}{}

	for _, p := range providers {
		if len(selected) == 3 {
			break
		}
		pool := append([]*Agent(nil), byProvider[p]...)
		s.rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
		pick := pool[0]
		selected = append(selected, pick)
		usedIDs[pick.ID.Hex()] = struct{}{}
	}

	for len(selected) < 3 {
		candidatePool := make([]*Agent, 0, len(agents))
		for _, a := range agents {
			if _, used := usedIDs[a.ID.Hex()]; !used {
				candidatePool = append(candidatePool, a)
			}
		}
		if len(candidatePool) == 0 {
			candidatePool = agents
		}
		pick := candidatePool[s.rand.Intn(len(candidatePool))]
		selected = append(selected, pick)
		usedIDs[pick.ID.Hex()] = struct{}{}
	}

	selection.Producer = *selected[0]
	selection.Reviewer = *selected[1]
	selection.Refiner = *selected[2]
	selection.SelectionReason["producer"] = s.reasonForSelection(selected[0], byProvider)
	selection.SelectionReason["reviewer"] = s.reasonForSelection(selected[1], byProvider)
	selection.SelectionReason["refiner"] = s.reasonForSelection(selected[2], byProvider)

	providerSet := map[Provider]struct{}{
		selection.Producer.Provider: {},
		selection.Reviewer.Provider: {},
		selection.Refiner.Provider:  {},
	}
	if len(providerSet) < 3 {
		selection.Warnings = append(selection.Warnings, "não há três providers distintos disponíveis; fallback com repetição de provider aplicado")
	}

	return selection
}

func (s *Service) reasonForSelection(a *Agent, byProvider map[Provider][]*Agent) string {
	candidates := len(byProvider[a.Provider])
	if candidates <= 1 {
		return fmt.Sprintf("provider %s possuía 1 candidato disponível", a.Provider)
	}
	return fmt.Sprintf("sorteado aleatoriamente de %d candidatos do provider %s", candidates, a.Provider)
}

func cloneAgents(in []*Agent) []Agent {
	out := make([]Agent, 0, len(in))
	for _, a := range in {
		out = append(out, *a)
	}
	return out
}
