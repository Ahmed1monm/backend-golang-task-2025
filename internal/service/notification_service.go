package service

import (
	"context"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/repository"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, userID uint, notificationType models.NotificationType, title, message string, wsEvent *websocket.Event) error
}

type notificationService struct {
	db         *gorm.DB
	notifyRepo repository.NotificationRepository
	wsManager  *websocket.Manager
}

func NewNotificationService(db *gorm.DB, notifyRepo repository.NotificationRepository, wsManager *websocket.Manager) NotificationService {
	return &notificationService{
		db:         db,
		notifyRepo: notifyRepo,
		wsManager:  wsManager,
	}
}

func (s *notificationService) CreateNotification(ctx context.Context, userID uint, notificationType models.NotificationType, title, message string, wsEvent *websocket.Event) error {
	notification := &models.Notification{
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Message: message,
	}

	// Start a new transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		logger.Error(ctx, "Failed to start notification transaction",
			zap.Error(tx.Error),
			zap.Uint("user_id", userID))
		return tx.Error
	}
	defer tx.Rollback()

	// Create notification using repository
	if err := s.notifyRepo.CreateInTx(ctx, tx, notification); err != nil {
		logger.Error(ctx, "Failed to create notification",
			zap.Error(err),
			zap.Uint("user_id", userID))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		logger.Error(ctx, "Failed to commit notification transaction",
			zap.Error(err),
			zap.Uint("user_id", userID))
		return err
	}

	// Send WebSocket notification if provided
	if wsEvent != nil {
		s.wsManager.SendToUser(userID, *wsEvent)
	}

	return nil
}

