package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/contextkey"
)
const (
	TraceIDKey contextkey.Key = "traceID"
)
// TraceMiddleware injects a trace ID into the request context
func TraceMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get trace ID from header or generate new one
			traceID := c.Request().Header.Get("X-Trace-ID")
			if traceID == "" {
				traceID = uuid.New().String()
			}

			// Add trace ID to response header
			c.Response().Header().Set("X-Trace-ID", traceID)
			// Create new context with trace ID using custom key type
			ctx := context.WithValue(c.Request().Context(), TraceIDKey, traceID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
