package models

import (
	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	ProductID    uint    `gorm:"uniqueIndex;not null"`
	Product      *Product `gorm:"foreignKey:ProductID"`
	Quantity     int     `gorm:"not null"`
	ReorderPoint int     `gorm:"not null;default:10"` // Threshold for reordering
	Reserved     int     `gorm:"not null;default:0"`  // Quantity reserved for pending orders
}
