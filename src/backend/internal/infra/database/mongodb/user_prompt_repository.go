package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/develop-agent/backend/internal/domain/prompt"
)

type UserPromptRepository struct {
	col    *mongo.Collection
	client *mongo.Client
}

func NewUserPromptRepository(adapter *Adapter, dbName string) *UserPromptRepository {
	return &UserPromptRepository{col: adapter.Client.Database(dbName).Collection("user_prompts"), client: adapter.Client}
}

func (r *UserPromptRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "group", Value: 1}, {Key: "priority", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "enabled", Value: 1}}},
	})
	return err
}

func (r *UserPromptRepository) Create(ctx context.Context, p *prompt.UserPrompt) error {
	_, err := r.col.InsertOne(ctx, p)
	return err
}

func (r *UserPromptRepository) FindByID(ctx context.Context, id string) (*prompt.UserPrompt, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var out prompt.UserPrompt
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *UserPromptRepository) FindByUserAndGroup(ctx context.Context, userID string, group prompt.Group) ([]*prompt.UserPrompt, error) {
	uid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	cur, err := r.col.Find(ctx, bson.M{"user_id": uid, "group": group, "enabled": true}, options.Find().SetSort(bson.D{{Key: "priority", Value: 1}, {Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer func() { _ = cur.Close(ctx) }()
	out := make([]*prompt.UserPrompt, 0)
	for cur.Next(ctx) {
		var p prompt.UserPrompt
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		out = append(out, &p)
	}
	return out, cur.Err()
}

func (r *UserPromptRepository) FindAllByUser(ctx context.Context, filter prompt.ListFilter) ([]*prompt.UserPrompt, error) {
	uid, err := bson.ObjectIDFromHex(filter.UserID)
	if err != nil {
		return nil, err
	}
	query := bson.M{"user_id": uid}
	if filter.Group != "" {
		query["group"] = filter.Group
	}
	if filter.Enabled != nil {
		query["enabled"] = *filter.Enabled
	}
	cur, err := r.col.Find(ctx, query, options.Find().SetSort(bson.D{{Key: "group", Value: 1}, {Key: "priority", Value: 1}, {Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer func() { _ = cur.Close(ctx) }()
	out := make([]*prompt.UserPrompt, 0)
	for cur.Next(ctx) {
		var p prompt.UserPrompt
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		out = append(out, &p)
	}
	return out, cur.Err()
}

func (r *UserPromptRepository) CountByUserAndGroup(ctx context.Context, userID string, group prompt.Group) (int64, error) {
	uid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return 0, err
	}
	return r.col.CountDocuments(ctx, bson.M{"user_id": uid, "group": group})
}

func (r *UserPromptRepository) Update(ctx context.Context, p *prompt.UserPrompt) error {
	p.UpdatedAt = time.Now().UTC()
	res, err := r.col.UpdateOne(ctx, bson.M{"_id": p.ID, "user_id": p.UserID}, bson.M{"$set": bson.M{
		"title":      p.Title,
		"content":    p.Content,
		"group":      p.Group,
		"priority":   p.Priority,
		"enabled":    p.Enabled,
		"tags":       p.Tags,
		"updated_at": p.UpdatedAt,
	}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("prompt not found")
	}
	return nil
}

func (r *UserPromptRepository) Delete(ctx context.Context, userID, id string) error {
	uid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": oid, "user_id": uid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("prompt not found")
	}
	return nil
}

func (r *UserPromptRepository) Reorder(ctx context.Context, userID string, items []prompt.ReorderItem) error {
	if len(items) == 0 {
		return nil
	}
	uid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc context.Context) (interface{}, error) {
		now := time.Now().UTC()
		for _, item := range items {
			oid, err := bson.ObjectIDFromHex(item.ID)
			if err != nil {
				return nil, err
			}
			res, err := r.col.UpdateOne(sc, bson.M{"_id": oid, "user_id": uid}, bson.M{"$set": bson.M{"priority": item.Priority, "updated_at": now}})
			if err != nil {
				return nil, err
			}
			if res.MatchedCount == 0 {
				return nil, errors.New("prompt not found")
			}
		}
		return nil, nil
	})
	return err
}
