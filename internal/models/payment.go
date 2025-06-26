package models

import (
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
)

type Payment struct {
	gorm.Model
	OrderID     uint          `gorm:"uniqueIndex;not null"`
	Order       Order         `gorm:"foreignKey:OrderID"`
	Amount      float64       `gorm:"type:decimal(10,2);not null"`
	Status      PaymentStatus `gorm:"type:varchar(20);default:'pending'"`
	PaymentMethod string      `gorm:"size:50;not null"`
	TransactionID string      `gorm:"size:100"`
}
