package seed

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/develop-agent/backend/internal/domain/agent"
)

type AgentsSeeder struct {
	repo agent.Repository
}

func NewAgentsSeeder(repo agent.Repository) *AgentsSeeder {
	return &AgentsSeeder{repo: repo}
}

func (s *AgentsSeeder) Run(ctx context.Context) error {
	defaults := []*agent.Agent{}
	for _, spec := range defaultAgentSpecs() {
		a, err := agent.New(spec.name, spec.description, spec.provider, spec.model, spec.prompts, spec.skills, true, spec.apiKeyRef)
		if err != nil {
			return err
		}
		defaults = append(defaults, a)
	}

	for _, candidate := range defaults {
		_, err := s.repo.FindByName(ctx, candidate.Name)
		if err == nil {
			continue
		}
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
		if err := s.repo.Create(ctx, candidate); err != nil {
			return err
		}
	}
	return nil
}

type seedSpec struct {
	name        string
	description string
	provider    agent.Provider
	model       string
	prompts     []string
	skills      []agent.Skill
	apiKeyRef   string
}

func defaultAgentSpecs() []seedSpec {
	return []seedSpec{
		{"Entrevistador", "Especialista em descoberta de produto", agent.ProviderOpenAI, "gpt-4o", []string{"Faça perguntas para extrair requisitos claros e validáveis."}, []agent.Skill{agent.SkillProjectCreation}, "provider-openai-default"},
		{"Engenheiro de Requisitos", "Converte ideias em requisitos funcionais e não funcionais", agent.ProviderAnthropic, "claude-3-5-sonnet", []string{"Especifique regras de negócio com rastreabilidade."}, []agent.Skill{agent.SkillEngineering}, "provider-anthropic-default"},
		{"Arquiteto de Software", "Define arquitetura, dados e integrações", agent.ProviderGoogle, "gemini-2.5-pro", []string{"Projete arquitetura evolutiva e orientada a interfaces."}, []agent.Skill{agent.SkillArchitecture}, "provider-google-default"},
		{"Planejador de Roadmap", "Quebra escopo em épicos e tarefas", agent.ProviderOpenAI, "gpt-4o", []string{"Gere exclusivamente JSON válido (sem markdown) no schema de roadmap com phases/epics/tasks, estimando complexity e estimated_hours por task."}, []agent.Skill{agent.SkillPlanning}, "provider-openai-default"},
		{"Dev Frontend", "Implementa interfaces ricas e acessíveis", agent.ProviderAnthropic, "claude-3-5-sonnet", []string{"Escreva UI performática, acessível e consistente."}, []agent.Skill{agent.SkillDevelopmentFrontend}, "provider-anthropic-default"},
		{"Dev Backend Golang", "Constrói APIs resilientes em Go", agent.ProviderOpenAI, "gpt-4o", []string{"Aplique clean architecture, observabilidade e testes."}, []agent.Skill{agent.SkillDevelopmentBackend}, "provider-openai-default"},
		{"QA Engineer", "Especialista em testes e qualidade", agent.ProviderGoogle, "gemini-2.5-pro", []string{"Cubra cenários críticos com testes determinísticos."}, []agent.Skill{agent.SkillTesting}, "provider-google-default"},
		{"Security Engineer", "Focado em segurança de aplicação", agent.ProviderAnthropic, "claude-3-5-sonnet", []string{"Aplique OWASP, hardening e validações defensivas."}, []agent.Skill{agent.SkillSecurity}, "provider-anthropic-default"},
		{"Tech Writer", "Documenta APIs, guias e manuais", agent.ProviderOpenAI, "gpt-4o-mini", []string{"Documente de forma objetiva, prática e atualizada."}, []agent.Skill{agent.SkillDocumentation}, "provider-openai-default"},
		{"DevOps Engineer", "Automatiza deploy, CI/CD e infraestrutura", agent.ProviderGoogle, "gemini-2.5-flash", []string{"Implemente pipelines e IaC com segurança."}, []agent.Skill{agent.SkillDevOps}, "provider-google-default"},
		{"Landing Page Designer", "Cria páginas de alta conversão", agent.ProviderOpenAI, "gpt-4o", []string{"Priorize proposta de valor, copy e CTA."}, []agent.Skill{agent.SkillLandingPage}, "provider-openai-default"},
		{"Marketing Strategist", "Planeja campanhas de crescimento", agent.ProviderAnthropic, "claude-3-5-sonnet", []string{"Estruture campanhas orientadas a métricas."}, []agent.Skill{agent.SkillMarketing}, "provider-anthropic-default"},
	}
}
