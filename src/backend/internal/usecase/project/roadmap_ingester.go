package project

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"

	domain "github.com/develop-agent/backend/internal/domain/project"
)

type RoadmapIngestResult struct {
	TaskCount  int `json:"task_count"`
	PhaseCount int `json:"phase_count"`
	EpicCount  int `json:"epic_count"`
}

type RoadmapIngester struct {
	validator *RoadmapSchemaValidator
	tasks     domain.TaskRepository
}

func NewRoadmapIngester(tasks domain.TaskRepository) *RoadmapIngester {
	return &RoadmapIngester{validator: &RoadmapSchemaValidator{}, tasks: tasks}
}

func (ri *RoadmapIngester) Ingest(ctx context.Context, projectID string, raw []byte) (*RoadmapIngestResult, string, error) {
	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, "", errors.New("invalid project id")
	}

	doc, validationErr := ri.validator.Validate(raw)
	if validationErr != nil {
		return nil, "", validationErr
	}

	tasks := make([]*domain.Task, 0)
	epicCount := 0
	for _, phase := range doc.Phases {
		epicCount += len(phase.Epics)
		for _, epic := range phase.Epics {
			for _, t := range epic.Tasks {
				tasks = append(tasks, &domain.Task{
					ProjectID:      pid,
					PhaseID:        strings.TrimSpace(phase.ID),
					EpicID:         strings.TrimSpace(epic.ID),
					Title:          strings.TrimSpace(t.Title),
					Description:    strings.TrimSpace(t.Description),
					Type:           t.Type,
					Complexity:     t.Complexity,
					EstimatedHours: float64(t.EstimatedHours),
					Track:          t.Track,
					Dependencies:   t.Dependencies,
					Status:         domain.TaskTodo,
				})
			}
		}
	}

	if err := ri.tasks.BulkCreate(ctx, tasks); err != nil {
		return nil, "", err
	}
	canonical, _ := json.Marshal(doc)
	return &RoadmapIngestResult{TaskCount: len(tasks), PhaseCount: len(doc.Phases), EpicCount: epicCount}, string(canonical), nil
}
