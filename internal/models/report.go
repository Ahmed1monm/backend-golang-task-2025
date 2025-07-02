package models

import (
	"time"

	"gorm.io/gorm"
)

type DailySalesReport struct {
	gorm.Model
	Date                 time.Time       `gorm:"uniqueIndex;not null"`
	TotalOrders         int             `gorm:"not null"`
	PendingOrders       int             `gorm:"not null"`
	ProcessingOrders    int             `gorm:"not null"`
	ShippedOrders       int             `gorm:"not null"`
	DeliveredOrders     int             `gorm:"not null"`
	CancelledOrders     int             `gorm:"not null"`
	TotalRevenue        float64         `gorm:"type:decimal(10,2);not null"`
	AverageOrderValue   float64         `gorm:"type:decimal(10,2);not null"`
	UniqueCustomers     int             `gorm:"not null"`
	NewCustomers        int             `gorm:"not null"`
	TopProducts         []TopProduct    `gorm:"foreignKey:ReportID"`
	LowStockProducts    []LowStockAlert `gorm:"foreignKey:ReportID"`
	OrderFulfillmentRate float64        `gorm:"type:decimal(5,2);not null"` // Percentage
	CancellationRate    float64        `gorm:"type:decimal(5,2);not null"` // Percentage
}

type TopProduct struct {
	gorm.Model
	ReportID      uint    `gorm:"not null"`
	ProductID     uint    `gorm:"not null"`
	ProductName   string  `gorm:"size:100;not null"`
	QuantitySold  int     `gorm:"not null"`
	Revenue       float64 `gorm:"type:decimal(10,2);not null"`
	StockTurnover float64 `gorm:"type:decimal(5,2);not null"` // Sales quantity / Average inventory
}

type LowStockAlert struct {
	gorm.Model
	ReportID       uint    `gorm:"not null"`
	ProductID      uint    `gorm:"not null"`
	ProductName    string  `gorm:"size:100;not null"`
	CurrentStock   int     `gorm:"not null"`
	ReservedStock  int     `gorm:"not null"`
	ReorderPoint   int     `gorm:"not null"`
}
