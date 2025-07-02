package service

import (
	"context"
	"errors"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/api/dto"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/repository"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/hashing"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/jwt"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	GetUserProfile(ctx context.Context, id uint) (*dto.UserProfileResponse, error)
	UpdateUserProfile(ctx context.Context, userID uint, req *dto.UpdateUserProfileRequest) (*dto.UserProfileResponse, error)
}

type userService struct {
	userRepo     repository.UserRepository
	hashService  hashing.Service
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		userRepo:    repo,
		hashService: hashing.NewService(),
	}
}

func (s *userService) GetUserProfile(ctx context.Context, id uint) (*dto.UserProfileResponse, error) {
	// Get user from repository
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		logger.Error(ctx, "Failed to get user", zap.Error(err))
		return nil, err
	}

	return &dto.UserProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Active:    user.Active,
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error(ctx, "Failed to check existing user", zap.Error(err))
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := s.hashService.HashPassword(ctx, req.Password)
	if err != nil {
		return nil, err // Error already logged by hash service
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      models.RoleCustomer,
		Active:    true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.Error(ctx, "Failed to create user", zap.Error(err))
		return nil, err
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user)
	if err != nil {
		logger.Error(ctx, "Failed to generate token", zap.Error(err))
		return nil, err
	}

	return &dto.CreateUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Token:     token,
	}, nil
}

// UpdateUserProfile updates a user's profile information
func (s *userService) UpdateUserProfile(ctx context.Context, userID uint, req *dto.UpdateUserProfileRequest) (*dto.UserProfileResponse, error) {
	// Get user from repository
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		logger.Error(ctx, "Failed to get user", zap.Error(err))
		return nil, err
	}

	// Check if email is being updated and verify it's not taken
	if req.Email != nil && *req.Email != user.Email {
		existingUser, err := s.userRepo.FindByEmail(ctx, *req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(ctx, "Failed to check existing user", zap.Error(err))
			return nil, err
		}
		if existingUser != nil {
			return nil, errors.New("email is already taken")
		}
		user.Email = *req.Email
	}

	// Update password if provided
	if req.Password != nil {
		hashedPassword, err := s.hashService.HashPassword(ctx, *req.Password)
		if err != nil {
			return nil, err // Error already logged by hash service
		}
		user.Password = hashedPassword
	}

	// Update other fields if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	// Save updates
	if err := s.userRepo.Update(ctx, user); err != nil {
		logger.Error(ctx, "Failed to update user", zap.Error(err))
		return nil, err
	}

	return &dto.UserProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Active:    user.Active,
	}, nil
}
