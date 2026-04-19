package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	redislib "github.com/redis/go-redis/v9"

	"github.com/develop-agent/backend/internal/infra/cache/redis"
)

type RefreshStore interface {
	Save(ctx context.Context, token, userID string, expiresAt time.Time) error
	GetUserID(ctx context.Context, token string) (string, error)
	Delete(ctx context.Context, token string) error
}

type RedisRefreshStore struct {
	client *redis.Adapter
}

func NewRedisRefreshStore(client *redis.Adapter) *RedisRefreshStore {
	return &RedisRefreshStore{client: client}
}

func (s *RedisRefreshStore) Save(ctx context.Context, token, userID string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return errors.New("refresh token expiration must be in the future")
	}
	return s.client.Client.Set(ctx, key(token), userID, ttl).Err()
}

func (s *RedisRefreshStore) GetUserID(ctx context.Context, token string) (string, error) {
	res, err := s.client.Client.Get(ctx, key(token)).Result()
	if err != nil {
		if errors.Is(err, redislib.Nil) {
			return "", errors.New("refresh token not found")
		}
		return "", fmt.Errorf("get refresh token: %w", err)
	}
	return res, nil
}

func (s *RedisRefreshStore) Delete(ctx context.Context, token string) error {
	return s.client.Client.Del(ctx, key(token)).Err()
}

func key(token string) string {
	return "auth:refresh:" + token
}
