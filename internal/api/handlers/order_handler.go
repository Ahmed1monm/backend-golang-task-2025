package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/errors"
)

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	return errors.NewValidationError("Invalid order data", map[string]string{
		"items": "at least one item is required",
		"quantity": "quantity must be greater than 0",
	}, http.StatusBadRequest)
}

func (h *OrderHandler) ListOrders(c echo.Context) error {
	return errors.NewBusinessError("No orders found", errors.ErrCodeResourceNotFound, http.StatusNotFound)
}

func (h *OrderHandler) GetOrder(c echo.Context) error {
	return errors.NewBusinessError("Order not found", errors.ErrCodeResourceNotFound, http.StatusNotFound)
}

func (h *OrderHandler) CancelOrder(c echo.Context) error {
	return errors.NewBusinessError("Order cannot be cancelled", errors.ErrCodeInvalidOrderStatus, http.StatusUnprocessableEntity)
}

func (h *OrderHandler) GetOrderStatus(c echo.Context) error {
	return errors.NewServerError("Failed to retrieve order status", nil, http.StatusInternalServerError)
}
