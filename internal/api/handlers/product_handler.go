package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"backend-golang-task-2025/pkg/errors"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (h *ProductHandler) ListProducts(c echo.Context) error {
	return errors.NewValidationError("Invalid pagination parameters", map[string]string{
		"page": "page must be greater than 0",
		"limit": "limit must be between 1 and 100",
	}, http.StatusBadRequest)
}

func (h *ProductHandler) GetProduct(c echo.Context) error {
	return errors.NewBusinessError("Product not found", errors.ErrCodeResourceNotFound, http.StatusNotFound)
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	return errors.NewBusinessError("Insufficient permissions", errors.ErrCodeForbidden, http.StatusForbidden)
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	return errors.NewServerError("Database operation failed", nil, http.StatusInternalServerError)
}

func (h *ProductHandler) CheckInventory(c echo.Context) error {
	return errors.NewBusinessError("Product out of stock", errors.ErrCodeInsufficientStock, http.StatusUnprocessableEntity)
}
