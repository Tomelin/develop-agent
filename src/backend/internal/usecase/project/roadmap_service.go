package project

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	domain "github.com/develop-agent/backend/internal/domain/project"
)

func (s *Service) RoadmapSummary(ctx context.Context, projectID string) (*domain.RoadmapSummary, error) {
	return s.tasks.RoadmapSummary(ctx, projectID)
}

func (s *Service) ExportRoadmap(ctx context.Context, projectID, format string) (contentType, filename string, payload []byte, err error) {
	p, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return "", "", nil, err
	}
	if strings.TrimSpace(p.RoadmapJSON) == "" {
		return "", "", nil, errors.New("roadmap is not available")
	}
	items, err := s.tasks.ListByProject(ctx, domain.TaskListFilter{ProjectID: projectID})
	if err != nil {
		return "", "", nil, err
	}

	switch format {
	case "json":
		return "application/json", "roadmap.json", []byte(p.RoadmapJSON), nil
	case "csv", "jira":
		buf := &bytes.Buffer{}
		w := csv.NewWriter(buf)
		headers := []string{"task_id", "title", "description", "type", "complexity", "estimated_hours", "track", "phase_id", "epic_id", "status", "dependencies"}
		if err := w.Write(headers); err != nil {
			return "", "", nil, err
		}
		for _, t := range items {
			if err := w.Write([]string{t.ID.Hex(), t.Title, t.Description, string(t.Type), string(t.Complexity), fmt.Sprintf("%.0f", t.EstimatedHours), string(t.Track), t.PhaseID, t.EpicID, string(t.Status), strings.Join(t.Dependencies, "|")}); err != nil {
				return "", "", nil, err
			}
		}
		w.Flush()
		if err := w.Error(); err != nil {
			return "", "", nil, err
		}
		name := "roadmap.csv"
		if format == "jira" {
			name = "roadmap-jira.csv"
		}
		return "text/csv", name, buf.Bytes(), nil
	case "markdown":
		var doc domain.RoadmapDocument
		if err := json.Unmarshal([]byte(p.RoadmapJSON), &doc); err != nil {
			return "", "", nil, err
		}
		var b strings.Builder
		b.WriteString("# Roadmap\n\n")
		for _, phase := range doc.Phases {
			b.WriteString(fmt.Sprintf("## %s — %s\n\n", phase.ID, phase.Name))
			for _, epic := range phase.Epics {
				b.WriteString(fmt.Sprintf("### %s — %s\n\n", epic.ID, epic.Title))
				for _, task := range epic.Tasks {
					b.WriteString(fmt.Sprintf("- [%s] %s (%s/%s, %dh)\n", task.ID, task.Title, task.Type, task.Complexity, task.EstimatedHours))
				}
				b.WriteString("\n")
			}
		}
		return "text/markdown", "roadmap.md", []byte(b.String()), nil
	default:
		return "", "", nil, errors.New("invalid format (use json, csv, markdown or jira)")
	}
}
