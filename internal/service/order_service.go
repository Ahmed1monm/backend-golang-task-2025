package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/repository"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/errors"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"go.uber.org/zap"
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

// ListAllOrders returns a paginated list of all orders in the system
func (s *OrderService) ListAllOrders(ctx context.Context, page, perPage int) ([]models.Order, int64, error) {
	// Calculate offset
	offset := (page - 1) * perPage

	// Get total count
	total, err := s.orderRepo.CountOrders(ctx, s.db)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Get paginated orders with their items and user
	orders, err := s.orderRepo.ListOrders(ctx, s.db, offset, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch orders: %w", err)
	}

	return orders, total, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uint, status models.OrderStatus) (*models.Order, error) {
	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// Get order
	order, err := s.orderRepo.GetOrderByID(ctx, tx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Validate status transition
	if !isValidStatusTransition(order.Status, status) {
		return nil, errors.NewValidationError(
			"Invalid status transition",
			map[string]string{"status": fmt.Sprintf("cannot transition from %s to %s", order.Status, status)},
			http.StatusBadRequest,
		)
	}

	// Update status
	order.Status = status
	if err := s.orderRepo.Update(ctx, tx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

// isValidStatusTransition checks if the status transition is valid
func isValidStatusTransition(current, new models.OrderStatus) bool {
	// Define valid transitions
	validTransitions := map[models.OrderStatus][]models.OrderStatus{
		models.OrderStatusPending:    {models.OrderStatusProcessing, models.OrderStatusCancelled},
		models.OrderStatusProcessing: {models.OrderStatusShipped, models.OrderStatusCancelled},
		models.OrderStatusShipped:    {models.OrderStatusDelivered},
		models.OrderStatusDelivered:  {},
		models.OrderStatusCancelled:  {},
	}

	// Check if transition is valid
	for _, validStatus := range validTransitions[current] {
		if validStatus == new {
			return true
		}
	}

	return false
}

type orderResult struct {
	Order   *models.Order
	Error   error
	Success bool
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uint, items []struct {
	ProductID uint
	Quantity  int
}) (*models.Order, error) {
	// Create result channel with buffer to avoid goroutine leak
	resultChan := make(chan orderResult, 1)

	// Process order in goroutine
	go func() {
		defer close(resultChan)

		// Start transaction
		tx := s.db.Begin()
		if tx.Error != nil {
			resultChan <- orderResult{Error: tx.Error}
			return
		}
		defer tx.Rollback()

		// Create order with pending status
		order := &models.Order{
			UserID: userID,
			Status: models.OrderStatusPending,
		}

		// Calculate total and validate inventory for each item
		var totalAmount float64
		orderItems := make([]models.OrderItem, 0, len(items))

		for _, item := range items {
			// Get product with lock
			product, err := s.orderRepo.GetProductByID(ctx, tx, item.ProductID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					resultChan <- orderResult{Error: errors.NewValidationError(
						fmt.Sprintf("Product with ID %d not found", item.ProductID),
						map[string]string{"product_id": "product not found"},
						400,
					)}
					return
				}
				resultChan <- orderResult{Error: err}
				return
			}

			// Get and lock inventory
			inventory, err := s.inventoryRepo.GetForUpdate(ctx, tx, item.ProductID)
			if err != nil {
				resultChan <- orderResult{Error: err}
				return
			}

			// Check if enough stock
			if inventory.Quantity < item.Quantity {
				resultChan <- orderResult{Error: errors.NewValidationError(
					fmt.Sprintf("Insufficient stock for product %d", item.ProductID),
					map[string]string{"quantity": "insufficient stock"},
					400,
				)}
				return
			}

			// Update inventory
			inventory.Quantity -= item.Quantity
			inventory.Reserved += item.Quantity
			if err := s.inventoryRepo.Update(ctx, tx, inventory); err != nil {
				resultChan <- orderResult{Error: err}
				return
			}

			// Create order item
			orderItems = append(orderItems, models.OrderItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     product.Price,
			})
			totalAmount += float64(item.Quantity) * product.Price
		}

		order.TotalAmount = totalAmount
		order.OrderItems = orderItems

		// Create the order
		if err := s.orderRepo.CreateOrder(ctx, tx, order); err != nil {
			resultChan <- orderResult{Error: err}
			return
		}

		// Process mock payment (simulated success)
		order.Status = models.OrderStatusProcessing
		if err := s.orderRepo.Update(ctx, tx, order); err != nil {
			resultChan <- orderResult{Error: err}
			return
		}

		// Create notification asynchronously (non-blocking)
		go func(orderID, userID uint) {
			notification := &models.Notification{
				UserID:  userID,
				Type:    models.NotificationTypeOrder,
				Title:   "Order Placed Successfully",
				Message: fmt.Sprintf("Your order #%d has been placed and is being processed.", orderID),
			}
			// Use a separate transaction for notification
			ntx := s.db.Begin()
			if ntx.Error != nil {
				logger.Error(ctx, "Failed to start notification transaction",
					zap.Error(ntx.Error),
					zap.Uint("order_id", orderID),
					zap.Uint("user_id", userID))
				return
			}
			defer ntx.Rollback()

			if err := ntx.Create(notification).Error; err != nil {
				logger.Error(ctx, "Failed to create notification",
					zap.Error(err),
					zap.Uint("order_id", orderID),
					zap.Uint("user_id", userID))
				return
			}

			if err := ntx.Commit().Error; err != nil {
				logger.Error(ctx, "Failed to commit notification transaction",
					zap.Error(err),
					zap.Uint("order_id", orderID),
					zap.Uint("user_id", userID))
			}
		}(order.ID, userID)

		// Commit main transaction
		if err := tx.Commit().Error; err != nil {
			resultChan <- orderResult{Error: err}
			return
		}

		resultChan <- orderResult{Order: order, Success: true}
	}()

	// Wait for result or context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-resultChan:
		if result.Error != nil {
			return nil, result.Error
		}
		return result.Order, nil
	}
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

// GetOrderStatus returns the current status of an order and verifies the user has access to it
func (s *OrderService) GetOrderStatus(ctx context.Context, orderID, userID uint) (models.OrderStatus, error) {
	var status models.OrderStatus
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Get order with minimal fields needed
		order, err := s.orderRepo.GetOrderByID(ctx, tx, orderID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.NewBusinessError(
					"Order not found",
					errors.ErrCodeResourceNotFound,
					http.StatusNotFound,
				)
			}
			return fmt.Errorf("failed to get order: %w", err)
		}

		// Verify order belongs to user
		if order.UserID != userID {
			return errors.NewBusinessError(
				"Order does not belong to user",
				"UNAUTHORIZED_ACCESS",
				http.StatusForbidden,
			)
		}

		status = order.Status
		return nil
	})

	if err != nil {
		return "", err
	}

	return status, nil
}

// CancelOrder cancels an order if it's in a cancellable state and belongs to the given user
func (s *OrderService) CancelOrder(ctx context.Context, orderID, userID uint) (*models.Order, error) {
	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// Get order
	order, err := s.orderRepo.GetOrderByID(ctx, tx, orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(
				"Order not found",
				errors.ErrCodeResourceNotFound,
				http.StatusNotFound,
			)
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Verify order belongs to user
	if order.UserID != userID {
		return nil, errors.NewBusinessError(
			"Order does not belong to user",
			"UNAUTHORIZED_ACCESS",
			http.StatusForbidden,
		)
	}

	// Check if order can be cancelled
	if !isValidStatusTransition(order.Status, models.OrderStatusCancelled) {
		return nil, errors.NewBusinessError(
			"Order cannot be cancelled",
			"INVALID_STATUS_TRANSITION",
			http.StatusBadRequest,
		)
	}

	// Update order status to cancelled
	order.Status = models.OrderStatusCancelled
	if err := s.orderRepo.Update(ctx, tx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Release reserved inventory for all order items
	for _, item := range order.OrderItems {
		// Get and lock inventory
		inventory, err := s.inventoryRepo.GetForUpdate(ctx, tx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get inventory: %w", err)
		}

		// Release reserved quantity
		inventory.Reserved -= item.Quantity
		if err := s.inventoryRepo.Update(ctx, tx, inventory); err != nil {
			return nil, fmt.Errorf("failed to update inventory: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}
