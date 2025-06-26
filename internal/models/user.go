package models

import "gorm.io/gorm"

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleCustomer UserRole = "customer"
)

type User struct {
	gorm.Model
	Email     string   `gorm:"uniqueIndex;not null"`
	Password  string   `gorm:"not null"`
	FirstName string   `gorm:"size:100"`
	LastName  string   `gorm:"size:100"`
	Role      UserRole `gorm:"type:varchar(20);default:'customer'"`
	Active    bool     `gorm:"default:true"`
	Orders    []Order  `gorm:"foreignKey:UserID"`
}
