package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/develop-agent/backend/internal/domain/billing"
)

type BillingRepository struct {
	col *mongo.Collection
}

func NewBillingRepository(adapter *Adapter, dbName string) *BillingRepository {
	return &BillingRepository{col: adapter.Client.Database(dbName).Collection("billing_records")}
}

func (r *BillingRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "organization_id", Value: 1}, {Key: "user_id", Value: 1}, {Key: "project_id", Value: 1}}},
		{Keys: bson.D{{Key: "organization_id", Value: 1}, {Key: "user_id", Value: 1}, {Key: "timestamp", Value: -1}}},
		{Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "phase_number", Value: 1}}},
	})
	return err
}

func (r *BillingRepository) Create(ctx context.Context, record *billing.BillingRecord) error {
	record.ComputeTotals()
	_, err := r.col.InsertOne(ctx, record)
	return err
}

func (r *BillingRepository) Summary(ctx context.Context, filter billing.QueryFilter) (*billing.Summary, error) {
	query, err := buildBillingFilter(filter)
	if err != nil {
		return nil, err
	}
	summary := &billing.Summary{}

	totalPipeline := mongo.Pipeline{{{Key: "$match", Value: query}}, {{Key: "$group", Value: bson.M{"_id": nil, "cost_usd": bson.M{"$sum": "$estimated_cost_usd"}, "tokens": bson.M{"$sum": "$total_tokens"}}}}}
	if totalRows, err := r.aggregateGrouped(ctx, totalPipeline); err == nil && len(totalRows) > 0 {
		summary.TotalCostUSD = totalRows[0].CostUSD
		summary.TotalTokens = totalRows[0].Tokens
	}

	byProject, err := r.aggregateGrouped(ctx, mongo.Pipeline{{{Key: "$match", Value: query}}, {{Key: "$group", Value: bson.M{"_id": "$project_id", "cost_usd": bson.M{"$sum": "$estimated_cost_usd"}, "tokens": bson.M{"$sum": "$total_tokens"}, "executions": bson.M{"$sum": 1}}}}, {{Key: "$sort", Value: bson.M{"cost_usd": -1}}}})
	if err != nil {
		return nil, err
	}
	summary.ByProject = byProject

	byModel, err := r.aggregateGrouped(ctx, mongo.Pipeline{{{Key: "$match", Value: query}}, {{Key: "$group", Value: bson.M{"_id": "$model", "cost_usd": bson.M{"$sum": "$estimated_cost_usd"}, "tokens": bson.M{"$sum": "$total_tokens"}, "executions": bson.M{"$sum": 1}}}}, {{Key: "$sort", Value: bson.M{"cost_usd": -1}}}})
	if err != nil {
		return nil, err
	}
	summary.ByModel = byModel

	return summary, nil
}

func (r *BillingRepository) ProjectDetails(ctx context.Context, filter billing.QueryFilter) (*billing.ProjectDetails, error) {
	query, err := buildBillingFilter(filter)
	if err != nil {
		return nil, err
	}
	out := &billing.ProjectDetails{ProjectID: filter.ProjectID}
	if out.ByPhase, err = r.aggregateGrouped(ctx, mongo.Pipeline{{{Key: "$match", Value: query}}, {{Key: "$group", Value: bson.M{"_id": "$phase_name", "cost_usd": bson.M{"$sum": "$estimated_cost_usd"}, "tokens": bson.M{"$sum": "$total_tokens"}, "executions": bson.M{"$sum": 1}}}}, {{Key: "$sort", Value: bson.M{"cost_usd": -1}}}}); err != nil {
		return nil, err
	}
	if out.ByAgent, err = r.aggregateGrouped(ctx, mongo.Pipeline{{{Key: "$match", Value: query}}, {{Key: "$group", Value: bson.M{"_id": "$agent_name", "cost_usd": bson.M{"$sum": "$estimated_cost_usd"}, "tokens": bson.M{"$sum": "$total_tokens"}, "executions": bson.M{"$sum": 1}}}}, {{Key: "$sort", Value: bson.M{"cost_usd": -1}}}}); err != nil {
		return nil, err
	}
	if out.ByModel, err = r.aggregateGrouped(ctx, mongo.Pipeline{{{Key: "$match", Value: query}}, {{Key: "$group", Value: bson.M{"_id": "$model", "cost_usd": bson.M{"$sum": "$estimated_cost_usd"}, "tokens": bson.M{"$sum": "$total_tokens"}, "executions": bson.M{"$sum": 1}}}}, {{Key: "$sort", Value: bson.M{"cost_usd": -1}}}}); err != nil {
		return nil, err
	}
	for _, i := range out.ByPhase {
		out.TotalUSD += i.CostUSD
	}
	return out, nil
}

func (r *BillingRepository) ByModel(ctx context.Context, filter billing.QueryFilter) ([]billing.GroupedCostItem, error) {
	query, err := buildBillingFilter(filter)
	if err != nil {
		return nil, err
	}
	return r.aggregateGrouped(ctx, mongo.Pipeline{{{Key: "$match", Value: query}}, {{Key: "$group", Value: bson.M{"_id": "$model", "cost_usd": bson.M{"$sum": "$estimated_cost_usd"}, "tokens": bson.M{"$sum": "$total_tokens"}, "executions": bson.M{"$sum": 1}}}}, {{Key: "$sort", Value: bson.M{"cost_usd": -1}}}})
}

func (r *BillingRepository) ByPhase(ctx context.Context, filter billing.QueryFilter) ([]billing.GroupedCostItem, error) {
	query, err := buildBillingFilter(filter)
	if err != nil {
		return nil, err
	}
	return r.aggregateGrouped(ctx, mongo.Pipeline{{{Key: "$match", Value: query}}, {{Key: "$group", Value: bson.M{"_id": "$phase_name", "cost_usd": bson.M{"$sum": "$estimated_cost_usd"}, "tokens": bson.M{"$sum": "$total_tokens"}, "executions": bson.M{"$sum": 1}}}}, {{Key: "$sort", Value: bson.M{"cost_usd": -1}}}})
}

func (r *BillingRepository) TopProjects(ctx context.Context, filter billing.QueryFilter) ([]billing.GroupedCostItem, error) {
	query, err := buildBillingFilter(filter)
	if err != nil {
		return nil, err
	}
	pipeline := mongo.Pipeline{{{Key: "$match", Value: query}}, {{Key: "$group", Value: bson.M{"_id": "$project_id", "cost_usd": bson.M{"$sum": "$estimated_cost_usd"}, "tokens": bson.M{"$sum": "$total_tokens"}, "executions": bson.M{"$sum": 1}}}}, {{Key: "$sort", Value: bson.M{"cost_usd": -1}}}, {{Key: "$limit", Value: 10}}}
	return r.aggregateGrouped(ctx, pipeline)
}

func (r *BillingRepository) List(ctx context.Context, filter billing.QueryFilter) ([]billing.BillingRecord, int64, error) {
	query, err := buildBillingFilter(filter)
	if err != nil {
		return nil, 0, err
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	count, err := r.col.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	cur, err := r.col.Find(ctx, query, options.Find().SetSort(bson.M{"timestamp": -1}).SetSkip((page-1)*limit).SetLimit(limit))
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = cur.Close(ctx) }()
	items := make([]billing.BillingRecord, 0)
	for cur.Next(ctx) {
		var rec billing.BillingRecord
		if err := cur.Decode(&rec); err != nil {
			return nil, 0, err
		}
		items = append(items, rec)
	}
	return items, count, cur.Err()
}

func (r *BillingRepository) aggregateGrouped(ctx context.Context, pipeline mongo.Pipeline) ([]billing.GroupedCostItem, error) {
	cur, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer func() { _ = cur.Close(ctx) }()
	out := make([]billing.GroupedCostItem, 0)
	for cur.Next(ctx) {
		row := struct {
			ID         any     `bson:"_id"`
			CostUSD    float64 `bson:"cost_usd"`
			Tokens     int64   `bson:"tokens"`
			Executions int64   `bson:"executions"`
		}{}
		if err := cur.Decode(&row); err != nil {
			return nil, err
		}
		out = append(out, billing.GroupedCostItem{Key: fmt.Sprint(row.ID), CostUSD: row.CostUSD, Tokens: row.Tokens, Executions: row.Executions})
	}
	return out, cur.Err()
}

func buildBillingFilter(filter billing.QueryFilter) (bson.M, error) {
	uid, err := bson.ObjectIDFromHex(filter.UserID)
	if err != nil {
		return nil, err
	}
	orgID, err := bson.ObjectIDFromHex(filter.OrganizationID)
	if err != nil {
		return nil, err
	}
	query := bson.M{"user_id": uid, "organization_id": orgID}
	if filter.ProjectID != "" {
		pid, err := bson.ObjectIDFromHex(filter.ProjectID)
		if err != nil {
			return nil, err
		}
		query["project_id"] = pid
	}
	if filter.Provider != "" {
		query["provider"] = filter.Provider
	}
	if filter.From != nil || filter.To != nil {
		rangeQuery := bson.M{}
		if filter.From != nil {
			rangeQuery["$gte"] = filter.From.UTC()
		}
		if filter.To != nil {
			rangeQuery["$lte"] = filter.To.UTC()
		}
		query["timestamp"] = rangeQuery
	}
	return query, nil
}
