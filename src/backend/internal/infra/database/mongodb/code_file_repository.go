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

type CodeFileRepository struct {
	col        *mongo.Collection
	metricsCol *mongo.Collection
}

func NewCodeFileRepository(adapter *Adapter, dbName string) *CodeFileRepository {
	db := adapter.Client.Database(dbName)
	return &CodeFileRepository{col: db.Collection("code_files"), metricsCol: db.Collection("phase_5_metrics")}
}

func (r *CodeFileRepository) EnsureIndexes(ctx context.Context) error {
	if _, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "path", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "task_id", Value: 1}}},
	}); err != nil {
		return err
	}
	_, err := r.metricsCol.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{Key: "project_id", Value: 1}}, Options: options.Index().SetUnique(true)})
	return err
}

func (r *CodeFileRepository) Upsert(ctx context.Context, file *project.CodeFile) error {
	now := time.Now().UTC()
	if file.ID.IsZero() {
		file.ID = bson.NewObjectID()
	}
	if file.CreatedAt.IsZero() {
		file.CreatedAt = now
	}
	file.UpdatedAt = now
	file.Version = now

	_, err := r.col.UpdateOne(ctx,
		bson.M{"project_id": file.ProjectID, "path": file.Path},
		bson.M{"$set": file, "$setOnInsert": bson.M{"created_at": file.CreatedAt}},
		options.UpdateOne().SetUpsert(true),
	)
	return err
}

func (r *CodeFileRepository) ListByProject(ctx context.Context, projectID string) ([]*project.CodeFile, error) {
	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, err
	}
	cur, err := r.col.Find(ctx, bson.M{"project_id": pid}, options.Find().SetSort(bson.M{"path": 1}))
	if err != nil {
		return nil, err
	}
	defer func() { _ = cur.Close(ctx) }()

	out := make([]*project.CodeFile, 0)
	for cur.Next(ctx) {
		var item project.CodeFile
		if err := cur.Decode(&item); err != nil {
			return nil, err
		}
		out = append(out, &item)
	}
	return out, cur.Err()
}

func (r *CodeFileRepository) FindByID(ctx context.Context, projectID, fileID string) (*project.CodeFile, error) {
	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, err
	}
	fid, err := bson.ObjectIDFromHex(fileID)
	if err != nil {
		return nil, err
	}
	var item project.CodeFile
	if err := r.col.FindOne(ctx, bson.M{"_id": fid, "project_id": pid}).Decode(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CodeFileRepository) CountAutoRejections(ctx context.Context, projectID string) (int64, error) {
	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return 0, err
	}
	var row struct {
		Count int64 `bson:"auto_rejections"`
	}
	err = r.metricsCol.FindOne(ctx, bson.M{"project_id": pid}).Decode(&row)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return row.Count, nil
}

func (r *CodeFileRepository) IncrementAutoRejections(ctx context.Context, projectID string) error {
	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return err
	}
	_, err = r.metricsCol.UpdateOne(ctx,
		bson.M{"project_id": pid},
		bson.M{"$inc": bson.M{"auto_rejections": 1}, "$set": bson.M{"updated_at": time.Now().UTC()}, "$setOnInsert": bson.M{"project_id": pid, "created_at": time.Now().UTC()}},
		options.UpdateOne().SetUpsert(true),
	)
	return err
}
