package mongodb

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/develop-agent/backend/internal/domain/agent"
)

type AgentRepository struct {
	col *mongo.Collection
}

func NewAgentRepository(adapter *Adapter, dbName string) *AgentRepository {
	return &AgentRepository{col: adapter.Client.Database(dbName).Collection("agents")}
}

func (r *AgentRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "name", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "skills", Value: 1}, {Key: "enabled", Value: 1}}},
		{Keys: bson.D{{Key: "provider", Value: 1}}},
		{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
	})
	return err
}

func (r *AgentRepository) Create(ctx context.Context, a *agent.Agent) error {
	_, err := r.col.InsertOne(ctx, a)
	return err
}

func (r *AgentRepository) FindByID(ctx context.Context, id string) (*agent.Agent, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID, "deleted_at": bson.M{"$exists": false}}
	var out agent.Agent
	if err := r.col.FindOne(ctx, filter).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *AgentRepository) FindByName(ctx context.Context, name string) (*agent.Agent, error) {
	filter := bson.M{"name": strings.TrimSpace(name), "deleted_at": bson.M{"$exists": false}}
	var out agent.Agent
	if err := r.col.FindOne(ctx, filter).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *AgentRepository) Update(ctx context.Context, a *agent.Agent) error {
	a.UpdatedAt = time.Now().UTC()
	res, err := r.col.UpdateOne(ctx,
		bson.M{"_id": a.ID, "deleted_at": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{
			"name":           a.Name,
			"description":    a.Description,
			"provider":       a.Provider,
			"model":          a.Model,
			"system_prompts": a.SystemPrompts,
			"skills":         a.Skills,
			"enabled":        a.Enabled,
			"api_key_ref":    a.ApiKeyRef,
			"status":         a.Status,
			"updated_at":     a.UpdatedAt,
		}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("agent not found")
	}
	return nil
}

func (r *AgentRepository) Delete(ctx context.Context, id string) error {
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
		return errors.New("agent not found")
	}
	return nil
}

func (r *AgentRepository) List(ctx context.Context, filter agent.ListFilter) ([]*agent.Agent, error) {
	query := bson.M{"deleted_at": bson.M{"$exists": false}}
	if filter.Enabled != nil {
		query["enabled"] = *filter.Enabled
	}
	if filter.Provider != "" {
		query["provider"] = filter.Provider
	}
	if filter.Skill != "" {
		query["skills"] = filter.Skill
	}

	cur, err := r.col.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = cur.Close(ctx) }()

	var items []*agent.Agent
	for cur.Next(ctx) {
		var a agent.Agent
		if err := cur.Decode(&a); err != nil {
			return nil, err
		}
		items = append(items, &a)
	}
	return items, cur.Err()
}

func (r *AgentRepository) FindBySkill(ctx context.Context, skill agent.Skill) ([]*agent.Agent, error) {
	enabled := true
	return r.List(ctx, agent.ListFilter{Enabled: &enabled, Skill: skill})
}
