package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/api/dto"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/service"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/validator"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) ListProducts(c echo.Context) error {
	// Parse pagination query
	var query dto.PaginationQuery
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid pagination parameters")
	}

	// Set defaults if not provided
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	// Validate query
	if errs := validator.Validate(query); len(errs) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid pagination parameters")
	}

	// Get products from service
	resp, err := h.productService.ListProducts(c.Request().Context(), query.Page, query.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list products")
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) GetProduct(c echo.Context) error {
	// Parse product ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	// Get product from service
	resp, err := h.productService.GetProduct(c.Request().Context(), uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get product")
	}

	if resp == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse request body
	req := new(dto.CreateProductRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if errs := validator.Validate(req); len(errs) > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": errs})
	}

	// Create product
	resp, err := h.productService.CreateProduct(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create product")
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	// Parse product ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	// Parse request body
	req := new(dto.UpdateProductRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if errs := validator.Validate(req); len(errs) > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": errs})
	}

	// Update product
	resp, err := h.productService.UpdateProduct(c.Request().Context(), uint(id), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update product")
	}

	if resp == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) CheckInventory(c echo.Context) error {
	// Parse product ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	// Get inventory from service
	resp, err := h.productService.GetInventory(c.Request().Context(), uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get product inventory")
	}

	if resp == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product or inventory not found")
	}

	return c.JSON(http.StatusOK, resp)
}
