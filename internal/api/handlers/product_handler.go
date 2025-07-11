package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/api/dto"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/service"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/redis"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ProductHandler struct {
	productService service.ProductService
	redisService   redis.Service
}

func NewProductHandler(productService service.ProductService, redisService redis.Service) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		redisService:   redisService,
	}
}

// ListProducts godoc
// @Summary List all products
// @Description Get a paginated list of products
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {object} dto.PaginatedProductsResponse
// @Failure 400 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /products [get]
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

// GetProduct godoc
// @Summary Get product by ID
// @Description Get detailed information about a specific product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /products/{id} [get]
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

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product in the system
// @Tags products
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Product creation details"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /products [post]
// @Security BearerAuth
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

// UpdateProduct godoc
// @Summary Update product
// @Description Update an existing product's information
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body dto.UpdateProductRequest true "Product update details"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /products/{id} [put]
// @Security BearerAuth
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
	// Clear cache asynchronously
	go func() {
		ctx := context.Background()
		// Invalidate all product-related caches using pattern
		if err := h.redisService.InvalidatePattern(ctx, "/api/v1/products*"); err != nil {
			logger.Error(ctx, "Failed to invalidate products cache", zap.Error(err))
		}
	}()
	return c.JSON(http.StatusOK, resp)
}

// CheckInventory godoc
// @Summary Check product inventory
// @Description Get the current inventory level for a specific product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dto.InventoryResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /products/{id}/inventory [get]
// @Security BearerAuth
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
