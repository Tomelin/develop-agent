package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/develop-agent/backend/internal/domain/project"
)

type TaskRepository struct {
	col *mongo.Collection
}

func NewTaskRepository(adapter *Adapter, dbName string) *TaskRepository {
	return &TaskRepository{col: adapter.Client.Database(dbName).Collection("tasks")}
}

func (r *TaskRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "type", Value: 1}, {Key: "complexity", Value: 1}}},
	})
	return err
}

func (r *TaskRepository) BulkCreate(ctx context.Context, tasks []*project.Task) error {
	if len(tasks) == 0 {
		return nil
	}
	docs := make([]interface{}, 0, len(tasks))
	now := time.Now().UTC()
	for _, t := range tasks {
		if t.ID.IsZero() {
			t.ID = bson.NewObjectID()
		}
		t.CreatedAt = now
		t.UpdatedAt = now
		if t.Status == "" {
			t.Status = project.TaskTodo
		}
		docs = append(docs, t)
	}
	_, err := r.col.InsertMany(ctx, docs)
	return err
}

func (r *TaskRepository) ListByProject(ctx context.Context, filter project.TaskListFilter) ([]*project.Task, error) {
	pid, err := bson.ObjectIDFromHex(filter.ProjectID)
	if err != nil {
		return nil, err
	}
	query := bson.M{"project_id": pid}
	if filter.Type != "" {
		query["type"] = filter.Type
	}
	if filter.Complexity != "" {
		query["complexity"] = filter.Complexity
	}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	cur, err := r.col.Find(ctx, query, options.Find().SetSort(bson.M{"created_at": 1}))
	if err != nil {
		return nil, err
	}
	defer func() { _ = cur.Close(ctx) }()

	out := make([]*project.Task, 0)
	for cur.Next(ctx) {
		var t project.Task
		if err := cur.Decode(&t); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	return out, cur.Err()
}

func (r *TaskRepository) RoadmapSummary(ctx context.Context, projectID string) (*project.RoadmapSummary, error) {
	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, err
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"project_id": pid}}},
		bson.D{{Key: "$facet", Value: bson.M{
			"totals":        []bson.M{{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}, "critical_hours": bson.M{"$sum": bson.M{"$cond": []any{bson.M{"$eq": []any{"$complexity", project.ComplexityCritical}}, "$estimated_hours", 0}}}}}},
			"by_type":       []bson.M{{"$group": bson.M{"_id": "$type", "count": bson.M{"$sum": 1}, "hours": bson.M{"$sum": "$estimated_hours"}}}},
			"by_complexity": []bson.M{{"$group": bson.M{"_id": "$complexity", "count": bson.M{"$sum": 1}}}},
			"by_phase":      []bson.M{{"$group": bson.M{"_id": "$phase_id", "hours": bson.M{"$sum": "$estimated_hours"}}}},
			"phase_count":   []bson.M{{"$group": bson.M{"_id": "$phase_id"}}, {"$count": "count"}},
			"epic_count":    []bson.M{{"$group": bson.M{"_id": "$epic_id"}}, {"$count": "count"}},
		}}},
	}

	cur, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer func() { _ = cur.Close(ctx) }()

	type grouped struct {
		ID    string  `bson:"_id"`
		Count int64   `bson:"count"`
		Hours float64 `bson:"hours"`
	}
	type total struct {
		Count         int64   `bson:"count"`
		CriticalHours float64 `bson:"critical_hours"`
	}
	type countOnly struct {
		Count int64 `bson:"count"`
	}
	var rows []struct {
		Totals       []total     `bson:"totals"`
		ByType       []grouped   `bson:"by_type"`
		ByComplexity []grouped   `bson:"by_complexity"`
		ByPhase      []grouped   `bson:"by_phase"`
		PhaseCount   []countOnly `bson:"phase_count"`
		EpicCount    []countOnly `bson:"epic_count"`
	}
	if err := cur.All(ctx, &rows); err != nil {
		return nil, err
	}
	summary := &project.RoadmapSummary{
		TotalByType:       map[project.TaskType]int64{},
		TotalByComplexity: map[project.TaskComplexity]int64{},
		HoursByType:       map[project.TaskType]float64{},
		HoursByPhase:      map[string]float64{},
	}
	if len(rows) == 0 {
		return summary, nil
	}
	row := rows[0]
	if len(row.Totals) > 0 {
		summary.TotalTasks = row.Totals[0].Count
		summary.EstimatedCriticalPathHR = row.Totals[0].CriticalHours
	}
	for _, item := range row.ByType {
		t := project.TaskType(item.ID)
		summary.TotalByType[t] = item.Count
		summary.HoursByType[t] = item.Hours
	}
	for _, item := range row.ByComplexity {
		summary.TotalByComplexity[project.TaskComplexity(item.ID)] = item.Count
	}
	for _, item := range row.ByPhase {
		summary.HoursByPhase[item.ID] = item.Hours
	}
	if len(row.PhaseCount) > 0 {
		summary.PhaseCount = row.PhaseCount[0].Count
	}
	if len(row.EpicCount) > 0 {
		summary.EpicCount = row.EpicCount[0].Count
	}
	return summary, nil
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, projectID, taskID string, status project.TaskStatus) error {
	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return err
	}
	tid, err := bson.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}
	res, err := r.col.UpdateOne(ctx,
		bson.M{"_id": tid, "project_id": pid},
		bson.M{"$set": bson.M{"status": status, "updated_at": time.Now().UTC()}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}
