package dto

import (
	"time"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
)

type CreateOrderItemRequest struct {
	ProductID uint    `json:"product_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type OrderResponse struct {
	ID          uint              `json:"id"`
	UserID      uint              `json:"user_id"`
	TotalAmount float64           `json:"total_amount"`
	Status      string            `json:"status"`
	Items       []OrderItemResponse `json:"items"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type OrderItemResponse struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// OrderToResponse converts an Order model to an OrderResponse DTO
func OrderToResponse(order *models.Order) *OrderResponse {
	items := make([]OrderItemResponse, len(order.OrderItems))
	for i, item := range order.OrderItems {
		items[i] = OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:    item.Price,
		}
	}

	return &OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		Items:       items,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}
