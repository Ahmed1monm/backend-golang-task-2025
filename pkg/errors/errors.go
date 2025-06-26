package errors

import (
	"fmt"
	"net/http"
)

// AppError is the base error type for the application
type AppError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	ErrorCode  string `json:"error_code"`
	Internal   error  `json:"-"` // Internal error details (not exposed to API)
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

// ValidationError represents validation related errors
type ValidationError struct {
	AppError
	Fields map[string]string `json:"fields,omitempty"` // Field-specific validation errors
}

// NewValidationError creates a new validation error
func NewValidationError(message string, fields map[string]string, statusCode int) *ValidationError {
	return &ValidationError{
		AppError: AppError{
			Message:    message,
			StatusCode: statusCode,
			ErrorCode:  "VALIDATION_ERROR",
		},
		Fields: fields,
	}
}

// BusinessError represents business logic related errors
type BusinessError struct {
	AppError
}

// NewBusinessError creates a new business error
func NewBusinessError(message string, errorCode string, statusCode int) *BusinessError {
	return &BusinessError{
		AppError: AppError{
			Message:    message,
			StatusCode: statusCode,
			ErrorCode:  errorCode,
		},
	}
}

// ServerError represents internal server errors
type ServerError struct {
	AppError
}

// NewServerError creates a new server error
func NewServerError(message string, internal error, statusCode int) *ServerError {
	return &ServerError{
		AppError: AppError{
			Message:    message,
			StatusCode: statusCode,
			ErrorCode:  "INTERNAL_SERVER_ERROR",
			Internal:   internal,
		},
	}
}

// Common business error codes
const (
	ErrCodeInsufficientStock    = "INSUFFICIENT_STOCK"
	ErrCodeInvalidOrderStatus   = "INVALID_ORDER_STATUS"
	ErrCodePaymentFailed       = "PAYMENT_FAILED"
	ErrCodeResourceNotFound    = "RESOURCE_NOT_FOUND"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
)

// IsValidationError checks if the error is a ValidationError
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// IsBusinessError checks if the error is a BusinessError
func IsBusinessError(err error) bool {
	_, ok := err.(*BusinessError)
	return ok
}

// IsServerError checks if the error is a ServerError
func IsServerError(err error) bool {
	_, ok := err.(*ServerError)
	return ok
}

// GetStatusCode returns the HTTP status code for the error
func GetStatusCode(err error) int {
	switch e := err.(type) {
	case *ValidationError:
		return e.StatusCode
	case *BusinessError:
		return e.StatusCode
	case *ServerError:
		return e.StatusCode
	default:
		return http.StatusInternalServerError
	}
}
