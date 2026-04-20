package project

import (
	"encoding/json"
	"fmt"
	"strings"

	domain "github.com/develop-agent/backend/internal/domain/project"
)

type RoadmapSchemaValidator struct{}

func (v *RoadmapSchemaValidator) Validate(raw []byte) (*domain.RoadmapDocument, *domain.RoadmapValidationError) {
	var doc domain.RoadmapDocument
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, &domain.RoadmapValidationError{Issues: []domain.ValidationIssue{{Path: "$", Code: "INVALID_JSON", Message: err.Error()}}}
	}

	issues := make([]domain.ValidationIssue, 0)
	add := func(path, code, msg string) {
		issues = append(issues, domain.ValidationIssue{Path: path, Code: code, Message: msg})
	}

	if strings.TrimSpace(doc.ProjectID) == "" {
		add("$.project_id", "REQUIRED", "project_id is required")
	}
	if len(doc.Phases) == 0 {
		add("$.phases", "REQUIRED", "at least one phase is required")
	}

	taskIDs := make(map[string]struct{})
	deps := make(map[string][]string)
	for i, p := range doc.Phases {
		phasePath := fmt.Sprintf("$.phases[%d]", i)
		if strings.TrimSpace(p.ID) == "" {
			add(phasePath+".id", "REQUIRED", "phase id is required")
		}
		if strings.TrimSpace(p.Name) == "" {
			add(phasePath+".name", "REQUIRED", "phase name is required")
		}
		if p.Order <= 0 {
			add(phasePath+".order", "INVALID", "phase order must be >= 1")
		}
		if len(p.Epics) == 0 {
			add(phasePath+".epics", "REQUIRED", "phase must have at least one epic")
		}
		for j, e := range p.Epics {
			epicPath := fmt.Sprintf("%s.epics[%d]", phasePath, j)
			if strings.TrimSpace(e.ID) == "" {
				add(epicPath+".id", "REQUIRED", "epic id is required")
			}
			if strings.TrimSpace(e.Title) == "" {
				add(epicPath+".title", "REQUIRED", "epic title is required")
			}
			if len(e.Tasks) == 0 {
				add(epicPath+".tasks", "REQUIRED", "epic must have at least one task")
			}
			for k, t := range e.Tasks {
				taskPath := fmt.Sprintf("%s.tasks[%d]", epicPath, k)
				if strings.TrimSpace(t.ID) == "" {
					add(taskPath+".id", "REQUIRED", "task id is required")
				} else {
					if _, ok := taskIDs[t.ID]; ok {
						add(taskPath+".id", "DUPLICATE", "task id must be unique")
					}
					taskIDs[t.ID] = struct{}{}
					deps[t.ID] = append([]string(nil), t.Dependencies...)
				}
				if strings.TrimSpace(t.Title) == "" {
					add(taskPath+".title", "REQUIRED", "task title is required")
				}
				if !t.Type.IsValid() {
					add(taskPath+".type", "INVALID_ENUM", "invalid task type")
				}
				if !t.Complexity.IsValid() {
					add(taskPath+".complexity", "INVALID_ENUM", "invalid task complexity")
				}
				if t.Track != domain.TrackFrontend && t.Track != domain.TrackBackend && t.Track != domain.TrackFull {
					add(taskPath+".track", "INVALID_ENUM", "invalid task track")
				}
				if t.EstimatedHours < 1 || t.EstimatedHours > 200 {
					add(taskPath+".estimated_hours", "OUT_OF_RANGE", "estimated_hours must be between 1 and 200")
				}
			}
		}
	}

	for taskID, taskDeps := range deps {
		for _, dep := range taskDeps {
			if _, ok := taskIDs[dep]; !ok {
				add("$.dependencies", "UNKNOWN_DEPENDENCY", fmt.Sprintf("task %s references unknown dependency %s", taskID, dep))
			}
		}
	}

	if len(issues) > 0 {
		return nil, &domain.RoadmapValidationError{Issues: issues}
	}
	return &doc, nil
}
