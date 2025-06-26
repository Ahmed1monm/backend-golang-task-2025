package models

import (
	"gorm.io/gorm"
)

type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
)

type AuditLog struct {
	gorm.Model
	UserID    uint       `gorm:"not null"`
	User      User       `gorm:"foreignKey:UserID"`
	Action    ActionType `gorm:"type:varchar(20);not null"`
	EntityType string    `gorm:"size:50;not null"`
	EntityID   uint      `gorm:"not null"`
	OldValue   string    `gorm:"type:text"`
	NewValue   string    `gorm:"type:text"`
	IPAddress  string    `gorm:"size:45"`
}
