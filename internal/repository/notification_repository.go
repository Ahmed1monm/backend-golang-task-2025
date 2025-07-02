package repository

import (
	"context"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	CreateInTx(ctx context.Context, tx *gorm.DB, notification *models.Notification) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	if err := r.db.WithContext(ctx).Create(notification).Error; err != nil {
		logger.Error(ctx, "Failed to create notification", zap.Error(err))
		return err
	}
	return nil
}

func (r *notificationRepository) CreateInTx(ctx context.Context, tx *gorm.DB, notification *models.Notification) error {
	if err := tx.WithContext(ctx).Create(notification).Error; err != nil {
		logger.Error(ctx, "Failed to create notification in transaction", zap.Error(err))
		return err
	}
	return nil
}
