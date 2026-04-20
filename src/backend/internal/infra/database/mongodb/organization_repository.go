package mongodb

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/develop-agent/backend/internal/domain/organization"
)

type OrganizationRepository struct {
	col *mongo.Collection
}

func NewOrganizationRepository(adapter *Adapter, dbName string) *OrganizationRepository {
	return &OrganizationRepository{col: adapter.Client.Database(dbName).Collection("organizations")}
}

func (r *OrganizationRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "slug", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "plan", Value: 1}}},
	})
	return err
}

func (r *OrganizationRepository) Create(ctx context.Context, org *organization.Organization) error {
	_, err := r.col.InsertOne(ctx, org)
	return err
}

func (r *OrganizationRepository) FindByID(ctx context.Context, id string) (*organization.Organization, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var out organization.Organization
	if err := r.col.FindOne(ctx, bson.M{"_id": objectID}).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *OrganizationRepository) FindBySlug(ctx context.Context, slug string) (*organization.Organization, error) {
	var out organization.Organization
	if err := r.col.FindOne(ctx, bson.M{"slug": strings.ToLower(strings.TrimSpace(slug))}).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *OrganizationRepository) Update(ctx context.Context, org *organization.Organization) error {
	org.UpdatedAt = time.Now().UTC()
	res, err := r.col.UpdateOne(ctx, bson.M{"_id": org.ID}, bson.M{"$set": bson.M{
		"name":                   org.Name,
		"slug":                   org.Slug,
		"plan":                   org.Plan,
		"max_users":              org.MaxUsers,
		"max_projects_per_month": org.MaxProjectsPerMonth,
		"max_tokens_per_month":   org.MaxTokensPerMonth,
		"billing_email":          org.BillingEmail,
		"updated_at":             org.UpdatedAt,
	}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("organization not found")
	}
	return nil
}
