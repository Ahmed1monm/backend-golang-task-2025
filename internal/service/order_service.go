package service

import (
	"context"
	"fmt"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/repository"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/errors"
	"gorm.io/gorm"
)

type OrderService struct {
	db            *gorm.DB
	orderRepo     repository.OrderRepository
	inventoryRepo repository.InventoryRepository
}

func NewOrderService(db *gorm.DB, orderRepo repository.OrderRepository, inventoryRepo repository.InventoryRepository) *OrderService {
	return &OrderService{
		db:            db,
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uint, items []struct {
	ProductID uint
	Quantity  int
}) (*models.Order, error) {
	var order models.Order
	err := s.db.Transaction(func(tx *gorm.DB) error {
		order = models.Order{
			UserID: userID,
			Status: models.OrderStatusPending,
		}

		// Calculate total and validate inventory for each item
		var totalAmount float64
		orderItems := make([]models.OrderItem, 0, len(items))

		for _, item := range items {
			// Get product
			product, err := s.orderRepo.GetProductByID(ctx, tx, item.ProductID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return errors.NewValidationError(
						fmt.Sprintf("Product with ID %d not found", item.ProductID),
						map[string]string{"product_id": "product not found"},
						400,
					)
				}
				return err
			}

			// Get and lock inventory
			inventory, err := s.inventoryRepo.GetForUpdate(ctx, tx, item.ProductID)
			if err != nil {
				return err
			}

			// Check if enough inventory is available
			if inventory.Quantity-inventory.Reserved < item.Quantity {
				return errors.NewValidationError(
					fmt.Sprintf("Insufficient inventory for product %d", item.ProductID),
					map[string]string{"quantity": "insufficient inventory"},
					400,
				)
			}

			// Reserve the inventory
			inventory.Reserved += item.Quantity
			if err := s.inventoryRepo.Update(ctx, tx, inventory); err != nil {
				return err
			}

			// Create order item
			orderItem := models.OrderItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:    product.Price,
			}
			orderItems = append(orderItems, orderItem)
			totalAmount += float64(item.Quantity) * product.Price
		}

		order.TotalAmount = totalAmount
		order.OrderItems = orderItems

		// Create the order
		if err := s.orderRepo.CreateOrder(ctx, tx, &order); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID uint) (*models.Order, error) {
	var order *models.Order
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		order, err = s.orderRepo.GetOrderByID(ctx, tx, orderID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.NewBusinessError(
					"Order not found",
					errors.ErrCodeResourceNotFound,
					404,
				)
			}
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) ListOrdersByUserID(ctx context.Context, userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		orders, err = s.orderRepo.ListOrdersByUserID(ctx, tx, userID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return orders, nil
}
