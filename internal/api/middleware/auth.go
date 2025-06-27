package middleware

import (
	"net/http"
	"strings"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/jwt"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// UserContext is the key used to store user information in the context
const UserContext = "user"

// JWTAuthentication middleware authenticates requests using JWT tokens
func JWTAuthentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			authHeader := c.Request().Header.Get("Authorization")

			// Check if Authorization header exists and has the correct format
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				logger.Debug(ctx, "Missing or invalid Authorization header")
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid Authorization header")
			}

			// Extract the token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate the token
			claims, err := jwt.ValidateToken(tokenString)
			if err != nil {
				logger.Error(ctx, "Invalid JWT token", zap.Error(err))
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
			}

			// Store user information in context
			c.Set(UserContext, claims)
			return next(c)
		}
	}
}

// RequireRoles middleware checks if the authenticated user has one of the required roles
func RequireRoles(roles ...models.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			// Get claims from context (set by JWTAuthentication middleware)
			claims, ok := c.Get(UserContext).(*jwt.Claims)
			if !ok {
				logger.Error(ctx, "User context not found")
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}

			// Check if user has one of the required roles
			hasRole := false
			for _, role := range roles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				logger.Warn(ctx, "Insufficient permissions", 
					zap.String("user_email", claims.Email),
					zap.String("user_role", string(claims.Role)),
					zap.Strings("required_roles", convertRolesToStrings(roles)))
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
			}

			return next(c)
		}
	}
}

// GetAuthenticatedUser returns the authenticated user's claims from the context
func GetAuthenticatedUser(c echo.Context) (*jwt.Claims, error) {
	claims, ok := c.Get(UserContext).(*jwt.Claims)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}
	return claims, nil
}

// Helper function to convert UserRole slice to string slice for logging
func convertRolesToStrings(roles []models.UserRole) []string {
	strRoles := make([]string, len(roles))
	for i, role := range roles {
		strRoles[i] = string(role)
	}
	return strRoles
}
