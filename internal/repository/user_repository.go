package repository

import (
	"context"
	"time"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	GetUniqueCustomerStats(ctx context.Context, tx *gorm.DB, date time.Time) (total int, new int, err error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUniqueCustomerStats(ctx context.Context, tx *gorm.DB, date time.Time) (total int, new int, err error) {
	var totalCount, newCount int64

	// Get total unique customers for the day
	err = tx.WithContext(ctx).
		Model(&models.Order{}).
		Where("DATE(created_at) = DATE(?)", date).
		Distinct("user_id").
		Count(&totalCount).Error
	if err != nil {
		return 0, 0, err
	}

	// Get new customers (first order on this date)
	err = tx.WithContext(ctx).
		Model(&models.Order{}).
		Where("DATE(created_at) = DATE(?)", date).
		Where("user_id IN (?)",
			tx.Model(&models.Order{}).
				Select("user_id").
				Where("DATE(created_at) = DATE(?)", date).
				Group("user_id").
				Having("MIN(DATE(created_at)) = DATE(?)", date),
		).
		Distinct("user_id").
		Count(&newCount).Error

	return int(totalCount), int(newCount), err
}

// Update updates a user in the database
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
