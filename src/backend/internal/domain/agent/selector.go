package agent

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Triad struct {
	Producer Agent `json:"producer"`
	Reviewer Agent `json:"reviewer"`
	Refiner  Agent `json:"refiner"`
}

type AuditLogger interface {
	LogTriadSelection(ctx context.Context, skill Skill, triad Triad, dynamic bool) error
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
		rand:          rand.New(rand.NewSource(time.Now().UnixNano())),
		fixedBySkill:  fixed,
	}
}

func (s *Service) SelectTriad(ctx context.Context, skill Skill) (Triad, error) {
	if !skill.IsValid() {
		return Triad{}, errors.New("invalid skill")
	}

	if !s.dynamicModeOn {
		triad, ok := s.fixedBySkill[skill]
		if !ok {
			return Triad{}, fmt.Errorf("no fixed triad configured for skill %s", skill)
		}
		if s.audit != nil {
			_ = s.audit.LogTriadSelection(ctx, skill, triad, false)
		}
		return triad, nil
	}

	agents, err := s.repo.FindBySkill(ctx, skill)
	if err != nil {
		return Triad{}, err
	}
	if len(agents) == 0 {
		return Triad{}, fmt.Errorf("no enabled agents found for skill %s", skill)
	}

	selected := s.pickThreeWithProviderDiversity(agents)
	triad := Triad{Producer: *selected[0], Reviewer: *selected[1], Refiner: *selected[2]}
	if s.audit != nil {
		_ = s.audit.LogTriadSelection(ctx, skill, triad, true)
	}
	return triad, nil
}

func (s *Service) pickThreeWithProviderDiversity(agents []*Agent) []*Agent {
	if len(agents) == 1 {
		return []*Agent{agents[0], agents[0], agents[0]}
	}

	pool := append([]*Agent(nil), agents...)
	s.rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })

	selected := make([]*Agent, 0, 3)
	usedProviders := map[Provider]struct{}{}

	for _, a := range pool {
		if len(selected) == 3 {
			break
		}
		if _, exists := usedProviders[a.Provider]; exists {
			continue
		}
		selected = append(selected, a)
		usedProviders[a.Provider] = struct{}{}
	}

	for len(selected) < 3 {
		selected = append(selected, pool[s.rand.Intn(len(pool))])
	}

	return selected
}
