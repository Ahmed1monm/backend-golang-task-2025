package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"go.uber.org/zap"
)

// Service defines the interface for Redis cache operations
type Service interface {
	GetCached(ctx context.Context, key string, dest interface{}) error
	Cache(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Invalidate(ctx context.Context, keys ...string) error
	InvalidatePattern(ctx context.Context, pattern string) error
}

type service struct {
	repo Repository
}

// NewService creates a new Redis service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// GetCached retrieves and unmarshals cached data
func (s *service) GetCached(ctx context.Context, key string, dest interface{}) error {
	data, err := s.repo.Get(ctx, key)
	if err != nil {
		if err == ErrNil {
			return err
		}
		logger.Error(ctx, "Failed to get cached data", 
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		logger.Error(ctx, "Failed to unmarshal cached data", 
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	return nil
}

// Cache marshals and stores data in Redis.
// If expiration is 0, the key will persist until manually invalidated.
// Otherwise, the key will be automatically removed after the expiration duration.
func (s *service) Cache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		logger.Error(ctx, "Failed to marshal data for caching", 
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	if err := s.repo.Set(ctx, key, data, expiration); err != nil {
		logger.Error(ctx, "Failed to cache data", 
			zap.String("key", key),
			zap.Duration("expiration", expiration),
			zap.Error(err))
		return err
	}

	logger.Debug(ctx, "Successfully cached data",
		zap.String("key", key),
		zap.Duration("expiration", expiration))
	return nil
}

// Invalidate removes one or more cached items
func (s *service) Invalidate(ctx context.Context, keys ...string) error {
	if err := s.repo.Del(ctx, keys...); err != nil {
		logger.Error(ctx, "Failed to invalidate cache", 
			zap.Strings("keys", keys),
			zap.Error(err))
		return err
	}
	return nil
}

// InvalidatePattern removes all cached items matching a pattern
func (s *service) InvalidatePattern(ctx context.Context, pattern string) error {
	if err := s.repo.DeleteByPattern(ctx, pattern); err != nil {
		logger.Error(ctx, "Failed to invalidate cache pattern", 
			zap.String("pattern", pattern),
			zap.Error(err))
		return err
	}
	return nil
}
