package models

import (
	"time"

	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTypeOrder    NotificationType = "order"
	NotificationTypePayment  NotificationType = "payment"
	NotificationTypeInventory NotificationType = "inventory"
)

type Notification struct {
	gorm.Model
	UserID    uint             `gorm:"not null"`
	User      User             `gorm:"foreignKey:UserID"`
	Type      NotificationType `gorm:"type:varchar(20);not null"`
	Title     string          `gorm:"size:255;not null"`
	Message   string          `gorm:"type:text;not null"`
	Read      bool            `gorm:"default:false"`
	ReadAt    *time.Time
}
