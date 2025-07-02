package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Repository defines the interface for Redis operations
type Repository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	DeleteByPattern(ctx context.Context, pattern string) error
}

type repository struct {
	client *redis.Client
}

// NewRepository creates a new Redis repository
func NewRepository(client *redis.Client) Repository {
	return &repository{
		client: client,
	}
}

// Get retrieves a value from Redis by key
func (r *repository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set stores a value in Redis with an optional expiration.
// If expiration is 0, the key will persist until manually deleted.
func (r *repository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Del removes one or more keys from Redis
func (r *repository) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// DeleteByPattern removes all keys matching the given pattern
func (r *repository) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	return nil
}
