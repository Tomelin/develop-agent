package interview

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	domaininterview "github.com/develop-agent/backend/internal/domain/interview"
	domainproject "github.com/develop-agent/backend/internal/domain/project"
	"github.com/develop-agent/backend/pkg/agentsdk"
)

type Service struct {
	repo         domaininterview.Repository
	projects     domainproject.ProjectRepository
	provider     agentsdk.Provider
	stateMachine *domainproject.ProjectStateMachine
	broker       *Broker
}

func NewService(repo domaininterview.Repository, projects domainproject.ProjectRepository, provider agentsdk.Provider, broker *Broker) *Service {
	if broker == nil {
		broker = NewBroker()
	}
	return &Service{repo: repo, projects: projects, provider: provider, stateMachine: domainproject.NewProjectStateMachine(), broker: broker}
}

func (s *Service) Broker() *Broker { return s.broker }

func (s *Service) GetSession(ctx context.Context, projectID, ownerID string) (*domaininterview.InterviewSession, error) {
	if _, err := s.validateProjectOwner(ctx, projectID, ownerID); err != nil {
		return nil, err
	}
	session, err := s.repo.FindByProjectID(ctx, projectID)
	if err == nil {
		return session, nil
	}
	if err != mongo.ErrNoDocuments {
		return nil, err
	}
	return s.createSession(ctx, projectID)
}

func (s *Service) StreamMessage(ctx context.Context, projectID, ownerID, content string, onDelta func(string) error) (*domaininterview.InterviewSession, string, error) {
	if strings.TrimSpace(content) == "" {
		return nil, "", errors.New("content is required")
	}
	if _, err := s.validateProjectOwner(ctx, projectID, ownerID); err != nil {
		return nil, "", err
	}
	session, err := s.GetSession(ctx, projectID, ownerID)
	if err != nil {
		return nil, "", err
	}
	if !session.CanIterate() {
		return nil, "", errors.New("iteration limit reached (10)")
	}
	if err := session.AddMessage(domaininterview.RoleUser, content); err != nil {
		return nil, "", err
	}
	if err := s.repo.UpsertByProjectID(ctx, projectID, session); err != nil {
		return nil, "", err
	}

	req := agentsdk.CompletionRequest{Messages: buildPromptMessages(session)}
	stream, err := s.provider.Stream(ctx, req)
	if err != nil {
		return nil, "", err
	}

	var acc strings.Builder
	for ev := range stream {
		if ev.Err != nil {
			return nil, "", ev.Err
		}
		if ev.Delta != "" {
			acc.WriteString(ev.Delta)
			if onDelta != nil {
				if cbErr := onDelta(ev.Delta); cbErr != nil {
					return nil, "", cbErr
				}
			}
		}
	}

	assistantText := strings.TrimSpace(acc.String())
	if assistantText == "" {
		assistantText = "Preciso de mais contexto para continuar a entrevista."
	}
	if session.IterationCount+1 >= session.MaxIterations {
		assistantText += "\n\n⚠️ Você atingiu o limite de 10 iterações desta entrevista."
	}
	if err := session.AddMessage(domaininterview.RoleAssistant, assistantText); err != nil {
		return nil, "", err
	}
	session.IterationCount++
	if session.IterationCount >= 3 && session.Status == domaininterview.StatusActive {
		session.Status = domaininterview.StatusAwaitingConfirmation
	}
	if err := s.repo.UpsertByProjectID(ctx, projectID, session); err != nil {
		return nil, "", err
	}
	return session, assistantText, nil
}

func (s *Service) RegenerateVision(ctx context.Context, projectID, ownerID string) (*domaininterview.InterviewSession, error) {
	if _, err := s.validateProjectOwner(ctx, projectID, ownerID); err != nil {
		return nil, err
	}
	session, err := s.GetSession(ctx, projectID, ownerID)
	if err != nil {
		return nil, err
	}
	vision, err := s.generateVision(ctx, session)
	if err != nil {
		return nil, err
	}
	session.VisionMD = vision
	if err := s.repo.UpsertByProjectID(ctx, projectID, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) Confirm(ctx context.Context, projectID, ownerID string) (*domaininterview.InterviewSession, error) {
	p, err := s.validateProjectOwner(ctx, projectID, ownerID)
	if err != nil {
		return nil, err
	}
	session, err := s.GetSession(ctx, projectID, ownerID)
	if err != nil {
		return nil, err
	}
	if session.IterationCount < 3 {
		return nil, errors.New("minimum of 3 interactions is required before confirmation")
	}
	if strings.TrimSpace(session.VisionMD) == "" {
		vision, genErr := s.generateVision(ctx, session)
		if genErr != nil {
			return nil, genErr
		}
		session.VisionMD = vision
	}
	now := time.Now().UTC()
	session.Status = domaininterview.StatusCompleted
	session.CompletedAt = &now
	if err := s.repo.UpsertByProjectID(ctx, projectID, session); err != nil {
		return nil, err
	}

	p.SpecMD = session.VisionMD
	_ = s.stateMachine.TransitionPhaseStatus(p, 1, domainproject.PhaseInProgress, "interview started", ownerID)
	_ = s.stateMachine.TransitionPhaseStatus(p, 1, domainproject.PhaseReview, "vision generated", ownerID)
	_ = s.stateMachine.TransitionPhaseStatus(p, 1, domainproject.PhaseCompleted, "user confirmed interview", ownerID)
	if err := s.projects.Update(ctx, p); err != nil {
		return nil, err
	}

	s.broker.Emit(Event{ProjectID: projectID, Type: "PHASE_1_COMPLETED", Message: "Fase 1 concluída com sucesso."})
	return session, nil
}

func (s *Service) createSession(ctx context.Context, projectID string) (*domaininterview.InterviewSession, error) {
	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, err
	}
	session := domaininterview.NewSession(pid)
	if err := s.repo.UpsertByProjectID(ctx, projectID, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) validateProjectOwner(ctx context.Context, projectID, ownerID string) (*domainproject.Project, error) {
	p, err := s.projects.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if p.OwnerUserID.Hex() != ownerID {
		return nil, errors.New("project not found")
	}
	return p, nil
}

func buildPromptMessages(session *domaininterview.InterviewSession) []agentsdk.Message {
	messages := []agentsdk.Message{{Role: agentsdk.RoleSystem, Content: interviewerSystemPrompt()}}
	for _, m := range session.Messages {
		role := agentsdk.RoleAssistant
		if m.Role == domaininterview.RoleUser {
			role = agentsdk.RoleUser
		}
		messages = append(messages, agentsdk.Message{Role: role, Content: m.Content})
	}
	return messages
}

func interviewerSystemPrompt() string {
	return `Você é um Agente Entrevistador de produto.
Conduza uma entrevista orgânica e iterativa para entender o projeto.
Cubra obrigatoriamente: problema, público-alvo, diferencial competitivo, MVP (3-5 funcionalidades), tecnologia/restrições, prazo/urgência, integrações e monetização.
Sempre reformule o entendimento e peça confirmação antes de avançar.`
}

func (s *Service) generateVision(ctx context.Context, session *domaininterview.InterviewSession) (string, error) {
	history := make([]string, 0, len(session.Messages))
	for _, m := range session.Messages {
		history = append(history, fmt.Sprintf("- %s: %s", m.Role, m.Content))
	}
	prompt := "Gere um VISION.md em markdown com seções: Sumário Executivo, Problema a ser Resolvido, Público-Alvo, Proposta de Valor Única, Escopo do MVP (in/out), Requisitos Não-Funcionais, Tecnologias/restrições, Integrações, Modelo de Negócio, Próximos Passos. Quando faltar contexto, use 'A definir'.\n\nHistórico:\n" + strings.Join(history, "\n")
	resp, err := s.provider.Complete(ctx, agentsdk.CompletionRequest{Messages: []agentsdk.Message{{Role: agentsdk.RoleUser, Content: prompt}}})
	if err != nil {
		return "", err
	}
	content := strings.TrimSpace(resp.Message.Content)
	if content == "" {
		return defaultVisionDocument(), nil
	}
	return content, nil
}

func defaultVisionDocument() string {
	return `# VISION.md

## Sumário Executivo
A definir.

## Problema a ser Resolvido
A definir.

## Público-Alvo
A definir.

## Proposta de Valor Única
A definir.

## Escopo do MVP
- In scope: A definir.
- Out of scope: A definir.

## Requisitos Não-Funcionais
A definir.

## Tecnologias e Restrições
A definir.

## Integrações Necessárias
A definir.

## Modelo de Negócio
A definir.

## Próximos Passos
A definir.`
}
