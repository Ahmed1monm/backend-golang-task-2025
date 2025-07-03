package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/redis"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type CachedResponse struct {
	Status int         `json:"status"`
	Header http.Header `json:"header"`
	Body   []byte      `json:"body"`
}

// WithCache combines CheckCache and CacheResponse middlewares
func WithCache(redisService redis.Service, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Only cache GET requests
		if c.Request().Method != http.MethodGet {
			return next(c)
		}

		// Try to get from cache first
		key := c.Request().URL.Path
		var cachedResponse CachedResponse
		err := redisService.GetCached(c.Request().Context(), key, &cachedResponse)
		if err == redis.ErrNil {
			// Cache miss, continue to handler and cache the response
			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			c.Response().Writer = &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}

			if err := next(c); err != nil {
				return err
			}

			// After handler execution, cache the response
			response := CachedResponse{
				Status: c.Response().Status,
				Header: c.Response().Header(),
				Body:   resBody.Bytes(),
			}

			if err := redisService.Cache(c.Request().Context(), key, response, 0); err != nil {
				logger.Error(c.Request().Context(), "Error caching response", zap.Error(err))
			}

			return nil
		} else if err != nil {
			logger.Error(c.Request().Context(), "Error retrieving from cache", zap.Error(err))
			return next(c)
		}

		// Return cached response
		for k, v := range cachedResponse.Header {
			for _, val := range v {
				c.Response().Header().Set(k, val)
			}
		}
		return c.Blob(cachedResponse.Status, "application/json", cachedResponse.Body)
	}
}

// bodyDumpResponseWriter captures the response for caching
type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
