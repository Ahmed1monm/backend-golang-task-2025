package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/redis"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RateLimitConfig holds the configuration for rate limiting
type RateLimitConfig struct {
	// Requests per window
	Limit int
	// Window duration in seconds
	WindowSeconds int
}

// DefaultRateLimitConfig provides default rate limit settings from environment variables
var DefaultRateLimitConfig = RateLimitConfig{
	Limit:         utils.GetEnvAsInt("RATE_LIMIT_REQUESTS", 100),         // default: 100 requests
	WindowSeconds: utils.GetEnvAsInt("RATE_LIMIT_WINDOW_SECONDS", 3600), // default: 1 hour
}

// RateLimit middleware limits the number of requests per IP address
func RateLimit(config RateLimitConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			key := fmt.Sprintf("rate_limit:%s", ip)
			
			client := redis.GetClient()
			ctx := c.Request().Context()

			// Get the current count for this IP
			val, err := client.Get(ctx, key).Int()
			if err != nil && err.Error() != "redis: nil" {
				logger.Error(ctx, "Rate limit check failed", zap.Error(err), zap.String("ip", ip))
				return echo.NewHTTPError(http.StatusInternalServerError, "Rate limit check failed")
			}

			// If this is the first request, initialize the counter
			if err != nil && err.Error() == "redis: nil" {
				err = client.Set(ctx, key, 1, time.Duration(config.WindowSeconds)*time.Second).Err()
				if err != nil {
					logger.Error(ctx, "Rate limit initialization failed", zap.Error(err), zap.String("ip", ip))
					return echo.NewHTTPError(http.StatusInternalServerError, "Rate limit initialization failed")
				}
				logger.Debug(ctx, "Rate limit initialized", zap.String("ip", ip))
				return next(c)
			}

			// If we've exceeded the limit, return an error
			if val >= config.Limit {
				logger.Warn(ctx, "Rate limit exceeded", zap.String("ip", ip), zap.Int("count", val))
				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			}

			// Increment the counter
			err = client.Incr(ctx, key).Err()
			if err != nil {
				logger.Error(ctx, "Rate limit increment failed", zap.Error(err), zap.String("ip", ip))
				return echo.NewHTTPError(http.StatusInternalServerError, "Rate limit increment failed")
			}

			logger.Debug(ctx, "Rate limit request processed", zap.String("ip", ip), zap.Int("count", val+1))
			return next(c)
		}
	}
}
