package redis

import (
	"context"
	"encoding/json"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Adapter struct {
	Client *redis.Client
}

// NewAdapter creates a new Redis adapter.
func NewAdapter(addr, password string) (*Adapter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Adapter{Client: client}, nil
}

// Close disconnects the Redis client.
func (a *Adapter) Close() error {
	return a.Client.Close()
}

// Ping checks the connection to Redis and returns latency in milliseconds.
func (a *Adapter) Ping(ctx context.Context) (int64, error) {
	start := time.Now()
	err := a.Client.Ping(ctx).Err()
	return time.Since(start).Milliseconds(), err
}

// Set stores a struct as JSON string with a given TTL.
func (a *Adapter) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return a.Client.Set(ctx, key, bytes, ttl).Err()
}

// Get retrieves a JSON string and unmarshals it into a struct.
func (a *Adapter) Get(ctx context.Context, key string, dest interface{}) error {
	bytes, err := a.Client.Get(ctx, key).Bytes()
	if err != nil {
		return err // Could be redis.Nil
	}
	return json.Unmarshal(bytes, dest)
}

// Delete removes a key.
func (a *Adapter) Delete(ctx context.Context, key string) error {
	return a.Client.Del(ctx, key).Err()
}

// Exists checks if a key exists.
func (a *Adapter) Exists(ctx context.Context, key string) (bool, error) {
	res, err := a.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

// SetNX sets a key if it does not exist, useful for locks.
func (a *Adapter) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return a.Client.SetArgs(ctx, key, bytes, redis.SetArgs{Mode: "NX", TTL: ttl}).Err() == nil, nil // basic implementation
}
