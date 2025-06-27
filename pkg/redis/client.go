package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/utils"
	"github.com/redis/go-redis/v9"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"go.uber.org/zap"
)

var redisClient *redis.Client

// InitRedis initializes the Redis client
func InitRedis(ctx context.Context) error {
	host := utils.GetEnv("REDIS_HOST", "localhost")
	port := utils.GetEnv("REDIS_PORT", "6379")
	password := utils.GetEnv("REDIS_PASSWORD", "")
	db, _ := strconv.Atoi(utils.GetEnv("REDIS_DB", "0"))

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	// Test the connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error(ctx, "Failed to connect to Redis", zap.Error(err))
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info(ctx, "Successfully connected to Redis")
	return nil
}

// GetClient returns the Redis client instance
func GetClient() *redis.Client {
	return redisClient
}
