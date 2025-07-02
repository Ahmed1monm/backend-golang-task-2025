package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	redisClient *redis.Client
	// ErrNil is returned when a key doesn't exist in Redis
	ErrNil = redis.Nil
)

// Config holds Redis configuration
type Config struct {
	Host           string
	Port           string
	Password       string
	DB             int
	MaxRetries     int
	MinIdleConns   int
	PoolSize       int
	ConnectTimeout time.Duration
}

// GetDefaultConfig returns a default Redis configuration
func GetDefaultConfig() Config {
	return Config{
		Host:           utils.GetEnv("REDIS_HOST", "localhost"),
		Port:           utils.GetEnv("REDIS_PORT", "6379"),
		Password:       utils.GetEnv("REDIS_PASSWORD", ""),
		DB:            func() int { db, _ := strconv.Atoi(utils.GetEnv("REDIS_DB", "0")); return db }(),
		MaxRetries:     3,
		MinIdleConns:   10,
		PoolSize:       20,
		ConnectTimeout: time.Second * 5,
	}
}

// InitRedis initializes the Redis client
func InitRedis(ctx context.Context) error {
	return InitRedisWithConfig(ctx, GetDefaultConfig())
}

// InitRedisWithConfig initializes the Redis client with custom configuration
func InitRedisWithConfig(ctx context.Context, cfg Config) error {
	// Create new client with configuration
	redisClient = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              cfg.DB,
		MaxRetries:      cfg.MaxRetries,
		MinIdleConns:    cfg.MinIdleConns,
		PoolSize:        cfg.PoolSize,
		ConnMaxIdleTime: time.Minute * 5,
		DialTimeout:     cfg.ConnectTimeout,
		ReadTimeout:     cfg.ConnectTimeout,
		WriteTimeout:    cfg.ConnectTimeout,
	})

	// Test the connection with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	if err := redisClient.Ping(timeoutCtx).Err(); err != nil {
		logger.Error(ctx, "Failed to connect to Redis", 
			zap.String("host", cfg.Host),
			zap.String("port", cfg.Port),
			zap.Error(err))
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info(ctx, "Successfully connected to Redis",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port))
	return nil
}

// GetClient returns the Redis client instance
func GetClient() *redis.Client {
	return redisClient
}

// Close gracefully shuts down the Redis client
func Close() error {
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			return fmt.Errorf("error closing Redis connection: %w", err)
		}
		redisClient = nil
	}
	return nil
}
