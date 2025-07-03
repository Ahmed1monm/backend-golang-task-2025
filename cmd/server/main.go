package main

import (
	"context"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/api/middleware"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/api/routes"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/config"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/utils"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/redis"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	env := utils.GetEnv("ENV", "development")
	logger.Init(env)
	defer logger.Sync()

	// Create background context for startup operations
	ctx := context.Background()

	// Initialize Redis
	if err := redis.InitRedis(ctx); err != nil {
		logger.Fatal(ctx, "Failed to initialize Redis", zap.Error(err))
	}
	redisClient := redis.GetClient()
	redisRepo := redis.NewRepository(redisClient)
	redisService := redis.NewService(redisRepo)
	logger.Info(ctx, "Successfully connected to Redis")

	// Initialize Echo instance
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
	e.Use(middleware.TraceMiddleware())
	e.Use(middleware.ErrorHandler)
	e.Use(middleware.RateLimit(middleware.DefaultRateLimitConfig))

	// Initialize database
	dbConfig := config.NewDBConfig()
	db, err := dbConfig.Connect()
	if err != nil {
		logger.Fatal(ctx, "Failed to connect to database", zap.Error(err))
	}

	// Verify database connection
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal(ctx, "Failed to get database instance", zap.Error(err))
	}
	defer sqlDB.Close()

	// Auto-migrate database schemas
	if err := models.AutoMigrate(db); err != nil {
		logger.Fatal(ctx, "Failed to migrate database", zap.Error(err))
	}

	logger.Info(ctx, "Successfully connected to database and migrated schemas")

	// Setup routes
	routes.SetupRoutes(e, db, redisService)

	// Health check route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "healthy"})
	})

	// Start server
	port := ":" + utils.GetEnv("PORT", "8080")
	logger.Info(ctx, "Starting server", zap.String("port", port))

	if err := e.Start(port); err != nil {
		logger.Fatal(ctx, "Server failed to start", zap.Error(err))
	}
}
