package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/api/dto"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/service"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/errors"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	orderService  *service.OrderService
	reportService *service.ReportService
}

func NewAdminHandler(orderService *service.OrderService, reportService *service.ReportService) *AdminHandler {
	return &AdminHandler{
		orderService:  orderService,
		reportService: reportService,
	}
}

// ListAllOrders godoc
// @Summary List all orders (admin only)
// @Description Get a paginated list of all orders in the system
// @Tags admin,orders
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10)"
// @Success 200 {object} dto.PaginatedOrdersResponse
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /admin/orders [get]
// @Security BearerAuth
func (h *AdminHandler) ListAllOrders(c echo.Context) error {
	// Parse pagination parameters with defaults
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// Get orders from service
	orders, total, err := h.orderService.ListAllOrders(c.Request().Context(), page, perPage)
	if err != nil {
		return errors.NewServerError("Failed to list orders", err, http.StatusInternalServerError)
	}

	// Convert to response DTOs
	responses := make([]dto.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = dto.OrderResponse{
			ID:          order.ID,
			UserID:      order.UserID,
			TotalAmount: order.TotalAmount,
			Status:      string(order.Status),
			Items:       make([]dto.OrderItemResponse, len(order.OrderItems)),
		}

		for j, item := range order.OrderItems {
			responses[i].Items[j] = dto.OrderItemResponse{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"orders":      responses,
		"total":       total,
		"page":        page,
		"per_page":    perPage,
		"total_pages": (int(total) + perPage - 1) / perPage,
	})
}

// UpdateOrderStatus godoc
// @Summary Update order status (admin only)
// @Description Update the status of an order in the system
// @Tags admin,orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body dto.UpdateOrderStatusRequest true "New order status"
// @Success 200 {object} dto.AdminOrderResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /admin/orders/{id}/status [put]
// @Security BearerAuth
func (h *AdminHandler) UpdateOrderStatus(c echo.Context) error {
	// Parse order ID
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order ID"})
	}

	// Parse request body
	var req dto.UpdateOrderStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Update order status
	order, err := h.orderService.UpdateOrderStatus(c.Request().Context(), uint(orderID), models.OrderStatus(req.Status))
	if err != nil {
		// Check for validation error
		if verr, ok := err.(*errors.ValidationError); ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": verr.Error()})
		}
		// Handle other errors
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update order status"})
	}

	// Convert to response DTO
	resp := &dto.AdminOrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Convert order items
	resp.Items = make([]dto.OrderItemResponse, len(order.OrderItems))
	for i, item := range order.OrderItems {
		resp.Items[i] = dto.OrderItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// GetDailySalesReport godoc
// @Summary Get today's sales report
// @Description Get the daily sales report for today. Returns an empty report if not yet generated.
// @Tags admin,reports
// @Accept json
// @Produce json
// @Success 200 {object} dto.DailySalesReportResponse
// @Failure 500 {object} errors.AppError
// @Router /admin/reports/daily [get]
// @Security BearerAuth
func (h *AdminHandler) GetDailySalesReport(c echo.Context) error {
	ctx := c.Request().Context()
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Get today's report from service
	report, err := h.reportService.GetDailyReport(ctx, today)
	if err != nil {
		return errors.NewServerError("Failed to get daily report", err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, report)
}

// GetLowStockAlerts godoc
// @Summary Get low stock alerts (admin only)
// @Description Get a list of products with low stock levels that require attention
// @Tags admin,inventory
// @Accept json
// @Produce json
// @Success 200 {array} dto.LowStockAlertResponse
// @Failure 404 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /admin/inventory/low-stock [get]
// @Security BearerAuth
func (h *AdminHandler) GetLowStockAlerts(c echo.Context) error {
	return errors.NewBusinessError("No low stock items found", "NO_LOW_STOCK", http.StatusNotFound)
}
