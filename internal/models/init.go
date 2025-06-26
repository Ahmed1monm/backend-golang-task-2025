package models

import (
	"gorm.io/gorm"
)

// AutoMigrate automatically migrates all models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Product{},
		&Order{},
		&OrderItem{},
		&Payment{},
		&Inventory{},
		&Notification{},
		&AuditLog{},
	)
}
