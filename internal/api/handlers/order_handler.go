package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/api/dto"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/service"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/contextkey"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/errors"
)

type OrderHandler struct{
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	ctx := c.Request().Context()

	// Get user ID from context (assuming it's set by auth middleware)
	userID, ok := ctx.Value(contextkey.UserIDKey).(uint)
	if !ok {
		return errors.NewAuthorizationError("User not authenticated", nil, http.StatusUnauthorized)
	}

	// Validate request body
	var req dto.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewValidationError("Invalid request body", nil, http.StatusBadRequest)
	}

	if err := c.Validate(&req); err != nil {
		return errors.NewValidationError("Validation failed", map[string]string{
			"items": "at least one item is required",
			"quantity": "quantity must be greater than 0",
		}, http.StatusBadRequest)
	}

	// Convert request items to service format
	items := make([]struct {
		ProductID uint
		Quantity  int
	}, len(req.Items))
	for i, item := range req.Items {
		items[i] = struct {
			ProductID uint
			Quantity  int
		}{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	// Create order
	order, err := h.orderService.CreateOrder(ctx, userID, items)
	if err != nil {
		switch e := err.(type) {
		case *errors.ValidationError:
			return e
		default:
			return errors.NewServerError("Failed to create order", err, http.StatusInternalServerError)
		}
	}

	// Convert order to response
	response := dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		TotalAmount: order.TotalAmount,
		Status:      string(order.Status),
		Items:       make([]dto.OrderItemResponse, len(order.OrderItems)),
	}

	for i, item := range order.OrderItems {
		response.Items[i] = dto.OrderItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *OrderHandler) ListOrders(c echo.Context) error {
	ctx := c.Request().Context()

	// Get user ID from context (assuming it's set by auth middleware)
	userID, ok := ctx.Value(contextkey.UserIDKey).(uint)
	if !ok {
		return errors.NewAuthorizationError("User not authenticated", nil, http.StatusUnauthorized)
	}

	// Get orders for the user
	orders, err := h.orderService.ListOrdersByUserID(ctx, userID)
	if err != nil {
		return errors.NewServerError("Failed to list orders", err, http.StatusInternalServerError)
	}

	// If no orders found, return empty array instead of error
	if len(orders) == 0 {
		return c.JSON(http.StatusOK, []dto.OrderResponse{})
	}

	// Convert orders to response format
	response := make([]dto.OrderResponse, len(orders))
	for i, order := range orders {
		response[i] = dto.OrderResponse{
			ID:          order.ID,
			UserID:      order.UserID,
			TotalAmount: order.TotalAmount,
			Status:      string(order.Status),
			Items:       make([]dto.OrderItemResponse, len(order.OrderItems)),
		}

		for j, item := range order.OrderItems {
			response[i].Items[j] = dto.OrderItemResponse{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
		}
	}

	return c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) GetOrder(c echo.Context) error {
	ctx := c.Request().Context()

	// Get user ID from context (assuming it's set by auth middleware)
	userID, ok := ctx.Value(contextkey.UserIDKey).(uint)
	if !ok {
		return errors.NewAuthorizationError("User not authenticated", nil, http.StatusUnauthorized)
	}

	// Get order ID from path parameter
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return errors.NewValidationError("Invalid order ID", nil, http.StatusBadRequest)
	}

	// Get order
	order, err := h.orderService.GetOrderByID(ctx, uint(orderID))
	if err != nil {
		switch e := err.(type) {
		case *errors.BusinessError:
			return e
		default:
			return errors.NewServerError("Failed to get order", err, http.StatusInternalServerError)
		}
	}

	// Verify the order belongs to the user
	if order.UserID != userID {
		return errors.NewAuthorizationError("Not authorized to view this order", nil, http.StatusForbidden)
	}

	// Convert order to response format
	response := dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		TotalAmount: order.TotalAmount,
		Status:      string(order.Status),
		Items:       make([]dto.OrderItemResponse, len(order.OrderItems)),
	}

	for i, item := range order.OrderItems {
		response.Items[i] = dto.OrderItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) CancelOrder(c echo.Context) error {
	return errors.NewBusinessError("Order cannot be cancelled", errors.ErrCodeInvalidOrderStatus, http.StatusUnprocessableEntity)
}

func (h *OrderHandler) GetOrderStatus(c echo.Context) error {
	return errors.NewServerError("Failed to retrieve order status", nil, http.StatusInternalServerError)
}
