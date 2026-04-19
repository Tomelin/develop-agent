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
	defer cur.Close(ctx)

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
