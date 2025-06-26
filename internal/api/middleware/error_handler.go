package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"backend-golang-task-2025/pkg/errors"
	"backend-golang-task-2025/pkg/logger"
)

type ErrorResponse struct {
	Message    string            `json:"message"`
	ErrorCode  string            `json:"error_code,omitempty"`
	Fields     map[string]string `json:"fields,omitempty"`
	StatusCode int              `json:"status_code"`
}

func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		ctx := c.Request().Context()

		switch e := err.(type) {
		case *errors.ValidationError:
			logger.Warn(ctx, "Validation error", zap.Error(err))
			return c.JSON(e.StatusCode, ErrorResponse{
				Message:    e.Message,
				ErrorCode:  e.ErrorCode,
				Fields:     e.Fields,
				StatusCode: e.StatusCode,
			})

		case *errors.BusinessError:
			logger.Warn(ctx, "Business error", zap.Error(err))
			return c.JSON(e.StatusCode, ErrorResponse{
				Message:    e.Message,
				ErrorCode:  e.ErrorCode,
				StatusCode: e.StatusCode,
			})

		case *errors.ServerError:
			logger.Error(ctx, "Server error", zap.Error(e.Internal))
			return c.JSON(e.StatusCode, ErrorResponse{
				Message:    e.Message,
				ErrorCode:  e.ErrorCode,
				StatusCode: e.StatusCode,
			})

		default:
			logger.Error(ctx, "Unhandled error", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message:    "Internal server error",
				ErrorCode:  "INTERNAL_SERVER_ERROR",
				StatusCode: http.StatusInternalServerError,
			})
		}
	}
}
