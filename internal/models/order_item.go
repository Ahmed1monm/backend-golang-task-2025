package models

import (
	"gorm.io/gorm"
)

type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"not null"`
	Order     Order   `gorm:"foreignKey:OrderID"`
	ProductID uint    `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"type:decimal(10,2);not null"`
}
