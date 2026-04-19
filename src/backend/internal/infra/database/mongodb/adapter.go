package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Adapter struct {
	Client *mongo.Client
}

// NewAdapter creates a new MongoDB adapter.
func NewAdapter(uri string) (*Adapter, error) {
	opts := options.Client().
		ApplyURI(uri).
		SetMinPoolSize(10).
		SetMaxPoolSize(100).
		SetConnectTimeout(10 * time.Second).
		SetTimeout(30 * time.Second)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil, err
	}

	return &Adapter{Client: client}, nil
}

// Close disconnects the MongoDB client.
func (a *Adapter) Close(ctx context.Context) error {
	return a.Client.Disconnect(ctx)
}

// Ping checks the connection to MongoDB and returns latency in milliseconds.
func (a *Adapter) Ping(ctx context.Context) (int64, error) {
	start := time.Now()
	err := a.Client.Ping(ctx, readpref.Primary())
	return time.Since(start).Milliseconds(), err
}

// Repository is a generic interface for MongoDB CRUD operations.
type Repository[T any] interface {
	FindOne(ctx context.Context, filter interface{}) (*T, error)
	FindMany(ctx context.Context, filter interface{}) ([]*T, error)
	InsertOne(ctx context.Context, doc *T) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	Aggregate(ctx context.Context, pipeline mongo.Pipeline) ([]*T, error)
}
