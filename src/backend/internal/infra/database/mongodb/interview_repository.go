package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/develop-agent/backend/internal/domain/interview"
)

type InterviewRepository struct {
	col *mongo.Collection
}

func NewInterviewRepository(adapter *Adapter, dbName string) *InterviewRepository {
	return &InterviewRepository{col: adapter.Client.Database(dbName).Collection("interview_sessions")}
}

func (r *InterviewRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "project_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "updated_at", Value: -1}}},
	})
	return err
}

func (r *InterviewRepository) FindByProjectID(ctx context.Context, projectID string) (*interview.InterviewSession, error) {
	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, err
	}
	var out interview.InterviewSession
	if err := r.col.FindOne(ctx, bson.M{"project_id": pid}).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *InterviewRepository) UpsertByProjectID(ctx context.Context, projectID string, session *interview.InterviewSession) error {
	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("session is nil")
	}
	session.ProjectID = pid
	session.UpdatedAt = time.Now().UTC()

	_, err = r.col.ReplaceOne(ctx, bson.M{"project_id": pid}, session, options.Replace().SetUpsert(true))
	return err
}
