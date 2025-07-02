package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, tx *gorm.DB, order *models.Order) error
	GetProductByID(ctx context.Context, tx *gorm.DB, productID uint) (*models.Product, error)
	GetOrderByID(ctx context.Context, tx *gorm.DB, orderID uint) (*models.Order, error)
	ListOrdersByUserID(ctx context.Context, tx *gorm.DB, userID uint) ([]models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, tx *gorm.DB, order *models.Order) error {
	return tx.WithContext(ctx).Create(order).Error
}



func (r *orderRepository) GetProductByID(ctx context.Context, tx *gorm.DB, productID uint) (*models.Product, error) {
	var product models.Product
	err := tx.WithContext(ctx).First(&product, productID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *orderRepository) GetOrderByID(ctx context.Context, tx *gorm.DB, orderID uint) (*models.Order, error) {
	var order models.Order
	err := tx.WithContext(ctx).Preload("OrderItems").First(&order, orderID).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) ListOrdersByUserID(ctx context.Context, tx *gorm.DB, userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := tx.WithContext(ctx).Where("user_id = ?", userID).Preload("OrderItems").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}
