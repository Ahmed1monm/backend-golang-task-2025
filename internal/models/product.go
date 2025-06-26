package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `gorm:"size:255;not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"type:decimal(10,2);not null"`
	SKU         string  `gorm:"uniqueIndex;size:50;not null"`
	Inventory   *Inventory
	OrderItems  []OrderItem `gorm:"foreignKey:ProductID"`
}
