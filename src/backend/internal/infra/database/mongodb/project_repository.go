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

type ProjectRepository struct {
	col *mongo.Collection
}

func NewProjectRepository(adapter *Adapter, dbName string) *ProjectRepository {
	return &ProjectRepository{col: adapter.Client.Database(dbName).Collection("projects")}
}

func (r *ProjectRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "owner_user_id", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "owner_user_id", Value: 1}, {Key: "name", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	return err
}

func (r *ProjectRepository) Create(ctx context.Context, p *project.Project) error {
	_, err := r.col.InsertOne(ctx, p)
	return err
}

func (r *ProjectRepository) FindByID(ctx context.Context, id string) (*project.Project, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var out project.Project
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *ProjectRepository) FindByOwner(ctx context.Context, filter project.ProjectListFilter) ([]*project.Project, int64, error) {
	query, opts, err := buildProjectFilter(filter, false)
	if err != nil {
		return nil, 0, err
	}
	return r.findProjects(ctx, query, opts)
}

func (r *ProjectRepository) FindDashboardByOwner(ctx context.Context, filter project.ProjectListFilter) ([]*project.Project, int64, error) {
	query, opts, err := buildProjectFilter(filter, true)
	if err != nil {
		return nil, 0, err
	}
	return r.findProjects(ctx, query, opts)
}

func (r *ProjectRepository) Update(ctx context.Context, p *project.Project) error {
	p.UpdatedAt = time.Now().UTC()
	res, err := r.col.ReplaceOne(ctx, bson.M{"_id": p.ID, "owner_user_id": p.OwnerUserID}, p)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("project not found")
	}
	return nil
}

func (r *ProjectRepository) Archive(ctx context.Context, id string, ownerID string) error {
	pid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	oid, err := bson.ObjectIDFromHex(ownerID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	res, err := r.col.UpdateOne(ctx,
		bson.M{"_id": pid, "owner_user_id": oid},
		bson.M{"$set": bson.M{"status": project.ProjectArchived, "archived_at": now, "updated_at": now}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("project not found")
	}
	return nil
}

func (r *ProjectRepository) UpdatePhase(ctx context.Context, projectID, ownerID string, phase project.PhaseExecution) error {
	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return err
	}
	oid, err := bson.ObjectIDFromHex(ownerID)
	if err != nil {
		return err
	}
	res, err := r.col.UpdateOne(
		ctx,
		bson.M{"_id": pid, "owner_user_id": oid, "phases.phase_number": phase.PhaseNumber},
		bson.M{"$set": bson.M{"phases.$": phase, "updated_at": time.Now().UTC()}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("project or phase not found")
	}
	return nil
}

func (r *ProjectRepository) UpdateSpecMD(ctx context.Context, projectID, ownerID, specMD string) error {
	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return err
	}
	oid, err := bson.ObjectIDFromHex(ownerID)
	if err != nil {
		return err
	}
	res, err := r.col.UpdateOne(ctx,
		bson.M{"_id": pid, "owner_user_id": oid},
		bson.M{"$set": bson.M{"spec_md": specMD, "updated_at": time.Now().UTC()}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("project not found")
	}
	return nil
}

func buildProjectFilter(filter project.ProjectListFilter, dashboard bool) (bson.M, *options.FindOptionsBuilder, error) {
	ownerID, err := bson.ObjectIDFromHex(filter.OwnerID)
	if err != nil {
		return nil, nil, err
	}
	query := bson.M{"owner_user_id": ownerID}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	if filter.FlowType != "" {
		query["flow_type"] = filter.FlowType
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	opts := options.Find().SetSkip((page - 1) * limit).SetLimit(limit).SetSort(bson.M{"updated_at": -1})
	if dashboard {
		opts = opts.SetProjection(bson.M{"name": 1, "status": 1, "flow_type": 1, "current_phase_number": 1, "updated_at": 1})
	}
	return query, opts, nil
}

func (r *ProjectRepository) findProjects(ctx context.Context, query bson.M, opts *options.FindOptionsBuilder) ([]*project.Project, int64, error) {
	count, err := r.col.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	cur, err := r.col.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	items := make([]*project.Project, 0)
	for cur.Next(ctx) {
		var p project.Project
		if err := cur.Decode(&p); err != nil {
			return nil, 0, err
		}
		items = append(items, &p)
	}
	return items, count, cur.Err()
}
