package routes

import (
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/api/handlers"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/api/middleware"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/repository"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupRoutes configures all API routes
func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// Initialize handlers
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo, orderRepo, inventoryRepo, db)
	orderService := service.NewOrderService(db, orderRepo, inventoryRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)
	adminHandler := handlers.NewAdminHandler()

	// API v1 group
	v1 := e.Group("/api/v1")

	// User routes
	users := v1.Group("/users")
	users.POST("", userHandler.CreateUser)
	users.GET("/:id", userHandler.GetUserProfile)
	users.PUT("/:id", userHandler.UpdateUserProfile, middleware.JWTAuthentication())

	// Product routes
	products := v1.Group("/products")
	products.GET("", productHandler.ListProducts)
	products.GET("/:id", productHandler.GetProduct)
	products.POST("", productHandler.CreateProduct, middleware.JWTAuthentication())
	products.PUT("/:id", productHandler.UpdateProduct, middleware.JWTAuthentication())
	products.GET("/:id/inventory", productHandler.CheckInventory, middleware.JWTAuthentication())

	// Order routes
	orders := v1.Group("/orders")
	orders.POST("", orderHandler.CreateOrder, middleware.JWTAuthentication())
	orders.GET("", orderHandler.ListOrders, middleware.JWTAuthentication())
	orders.GET("/:id", orderHandler.GetOrder, middleware.JWTAuthentication())
	orders.PUT("/:id/cancel", orderHandler.CancelOrder, middleware.JWTAuthentication())
	orders.GET("/:id/status", orderHandler.GetOrderStatus, middleware.JWTAuthentication())

	// Admin routes - protected with JWT auth and admin role requirement
	admin := v1.Group("/admin", middleware.JWTAuthentication(), middleware.RequireRoles(models.RoleAdmin))
	admin.GET("/orders", adminHandler.ListAllOrders)
	admin.PUT("/orders/:id/status", adminHandler.UpdateOrderStatus)
	admin.GET("/reports/daily", adminHandler.GetDailySalesReport)
	admin.GET("/inventory/low-stock", adminHandler.GetLowStockAlerts)
}
