package mongodb

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/develop-agent/backend/internal/domain/user"
)

type UserRepository struct {
	col *mongo.Collection
}

func NewUserRepository(adapter *Adapter, dbName string) *UserRepository {
	return &UserRepository{col: adapter.Client.Database(dbName).Collection("users")}
}

func (r *UserRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "role", Value: 1}}},
		{Keys: bson.D{{Key: "organization_id", Value: 1}, {Key: "organization_role", Value: 1}}},
	})
	return err
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	_, err := r.col.InsertOne(ctx, u)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID, "deleted_at": bson.M{"$exists": false}}
	var out user.User
	if err := r.col.FindOne(ctx, filter).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	filter := bson.M{"email": strings.ToLower(strings.TrimSpace(email)), "deleted_at": bson.M{"$exists": false}}
	var out user.User
	if err := r.col.FindOne(ctx, filter).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	u.UpdatedAt = time.Now().UTC()
	filter := bson.M{"_id": u.ID, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{
		"organization_id":   u.OrganizationID,
		"organization_role": u.OrganizationRole,
		"name":              u.Name,
		"email":             u.Email,
		"password_hash":     u.PasswordHash,
		"role":              u.Role,
		"prompts":           u.Prompts,
		"enabled":           u.Enabled,
		"updated_at":        u.UpdatedAt,
	}}
	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *UserRepository) SoftDelete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	res, err := r.col.UpdateOne(ctx,
		bson.M{"_id": objectID, "deleted_at": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context) ([]*user.User, error) {
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}
	opts := options.Find().SetProjection(bson.M{"password_hash": 0})
	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []*user.User
	for cur.Next(ctx) {
		var u user.User
		if err := cur.Decode(&u); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, cur.Err()
}
