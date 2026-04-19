package project

import (
	"context"
	"fmt"
	"strings"
)

type InheritanceService struct {
	repo ProjectRepository
}

func NewInheritanceService(repo ProjectRepository) *InheritanceService {
	return &InheritanceService{repo: repo}
}

func (s *InheritanceService) BuildInitialContext(ctx context.Context, linkedProjectID string) (string, error) {
	base, err := s.repo.FindByID(ctx, linkedProjectID)
	if err != nil {
		return "", err
	}
	if base.FlowType != FlowSoftware {
		return "", fmt.Errorf("linked project must be SOFTWARE flow")
	}
	if base.Status != ProjectCompleted && base.Status != ProjectInProgress {
		return "", fmt.Errorf("linked project must be IN_PROGRESS or COMPLETED")
	}

	var sb strings.Builder
	sb.WriteString("# Contexto herdado\n")
	sb.WriteString("Projeto base: " + base.Name + "\n")
	sb.WriteString("Descrição: " + base.Description + "\n\n")
	if strings.TrimSpace(base.SpecMD) != "" {
		sb.WriteString("## SPEC acumulado\n")
		sb.WriteString(strings.TrimSpace(base.SpecMD) + "\n")
	}
	return sb.String(), nil
}
