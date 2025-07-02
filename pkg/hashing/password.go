package hashing

import (
	"context"

	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Service defines the interface for password hashing operations
type Service interface {
	HashPassword(ctx context.Context, password string) (string, error)
	ComparePasswords(ctx context.Context, hashedPassword, plainPassword string) bool
}

type service struct {
	cost int
}

// NewService creates a new hashing service
func NewService() Service {
	return &service{
		cost: bcrypt.DefaultCost,
	}
}

// HashPassword hashes a plain text password using bcrypt
func (s *service) HashPassword(ctx context.Context, password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		logger.Error(ctx, "Failed to hash password", zap.Error(err))
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePasswords compares a hashed password with a plain text password
func (s *service) ComparePasswords(ctx context.Context, hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		logger.Debug(ctx, "Password comparison failed", zap.Error(err))
		return false
	}
	return true
}
