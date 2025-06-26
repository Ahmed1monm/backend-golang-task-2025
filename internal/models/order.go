package models

import "gorm.io/gorm"

type OrderStatus string

const (
	OrderStatusPending     OrderStatus = "pending"
	OrderStatusProcessing  OrderStatus = "processing"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusFulfilled  OrderStatus = "fulfilled"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	gorm.Model
	UserID      uint        `gorm:"not null"`
	User        User        `gorm:"foreignKey:UserID"`
	OrderItems  []OrderItem
	TotalAmount float64     `gorm:"type:decimal(10,2);not null"`
	Status      OrderStatus `gorm:"type:varchar(20);default:'pending'"`
	PaymentID   *uint
	Payment     *Payment
}
