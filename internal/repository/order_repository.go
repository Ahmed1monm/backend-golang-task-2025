package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, tx *gorm.DB, order *models.Order) error
	GetProductByID(ctx context.Context, tx *gorm.DB, productID uint) (*models.Product, error)
	GetOrderByID(ctx context.Context, tx *gorm.DB, orderID uint) (*models.Order, error)
	ListOrdersByUserID(ctx context.Context, tx *gorm.DB, userID uint) ([]models.Order, error)
	ListOrders(ctx context.Context, tx *gorm.DB, offset, limit int) ([]models.Order, error)
	CountOrders(ctx context.Context, tx *gorm.DB) (int64, error)
	Update(ctx context.Context, tx *gorm.DB, order *models.Order) error
	GetOrderStatsByDate(ctx context.Context, tx *gorm.DB, date time.Time) (stats map[models.OrderStatus]int, revenue float64, err error)
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
	err := tx.WithContext(ctx).Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) ListOrders(ctx context.Context, tx *gorm.DB, offset, limit int) ([]models.Order, error) {
	var orders []models.Order
	err := tx.WithContext(ctx).
		Preload("OrderItems").
		Preload("User").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) CountOrders(ctx context.Context, tx *gorm.DB) (int64, error) {
	var total int64
	err := tx.WithContext(ctx).Model(&models.Order{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *orderRepository) Update(ctx context.Context, tx *gorm.DB, order *models.Order) error {
	return tx.WithContext(ctx).Save(order).Error
}

func (r *orderRepository) GetOrderStatsByDate(ctx context.Context, tx *gorm.DB, date time.Time) (map[models.OrderStatus]int, float64, error) {
	type Result struct {
		Status      models.OrderStatus
		Count       int
		TotalAmount float64
	}
	var results []Result

	err := tx.WithContext(ctx).
		Model(&models.Order{}).
		Select("status, COUNT(*) as count, SUM(total_amount) as total_amount").
		Where("DATE(created_at) = DATE(?)", date).
		Group("status").
		Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	stats := make(map[models.OrderStatus]int)
	var totalRevenue float64
	for _, result := range results {
		stats[result.Status] = result.Count
		if result.Status == models.OrderStatusDelivered {
			totalRevenue = result.TotalAmount
		}
	}

	return stats, totalRevenue, nil
}
