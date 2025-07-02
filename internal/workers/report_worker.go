package workers

import (
	"context"
	"time"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/repository"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReportWorker struct {
	db            *gorm.DB
	reportRepo    repository.ReportRepository
	orderRepo     repository.OrderRepository
	userRepo      repository.UserRepository
	productRepo   repository.ProductRepository
	cron         *cron.Cron
}

func NewReportWorker(
	db *gorm.DB,
	reportRepo repository.ReportRepository,
	orderRepo repository.OrderRepository,
	userRepo repository.UserRepository,
	productRepo repository.ProductRepository,
) *ReportWorker {
	worker := &ReportWorker{
		db:         db,
		reportRepo: reportRepo,
		orderRepo: orderRepo,
		userRepo:  userRepo,
		productRepo: productRepo,
		cron:      cron.New(cron.WithSeconds()),
	}

	return worker
}

func (w *ReportWorker) Start() error {
	// Schedule report generation for 23:59:59 every day
	_, err := w.cron.AddFunc("59 59 23 * * *", w.generateDailyReport)
	if err != nil {
		return err
	}

	w.cron.Start()
	return nil
}

func (w *ReportWorker) Stop() {
	w.cron.Stop()
}

func (w *ReportWorker) generateDailyReport() {
	ctx := context.Background()
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Start transaction
	tx := w.db.Begin()
	if tx.Error != nil {
		logger.Error(ctx, "Failed to start transaction", zap.Error(tx.Error))
		return
	}
	defer tx.Rollback()

	// Get order statistics
	orderStats, totalRevenue, err := w.orderRepo.GetOrderStatsByDate(ctx, tx, today)
	if err != nil {
		logger.Error(ctx, "Failed to get order stats", zap.Error(err))
		return
	}

	// Get customer statistics
	totalCustomers, newCustomers, err := w.userRepo.GetUniqueCustomerStats(ctx, tx, today)
	if err != nil {
		logger.Error(ctx, "Failed to get customer stats", zap.Error(err))
		return
	}

	// Get top products
	topProducts, err := w.productRepo.GetTopProducts(ctx, tx, today, 10)
	if err != nil {
		logger.Error(ctx, "Failed to get top products", zap.Error(err))
		return
	}

	// Get low stock alerts
	lowStockAlerts, err := w.productRepo.GetLowStockProducts(ctx, tx)
	if err != nil {
		logger.Error(ctx, "Failed to get low stock alerts", zap.Error(err))
		return
	}

	// Calculate rates
	totalOrders := 0
	for _, count := range orderStats {
		totalOrders += count
	}

	var fulfillmentRate, cancellationRate float64
	if totalOrders > 0 {
		fulfillmentRate = float64(orderStats[models.OrderStatusDelivered]) / float64(totalOrders) * 100
		cancellationRate = float64(orderStats[models.OrderStatusCancelled]) / float64(totalOrders) * 100
	}

	// Create report
	report := &models.DailySalesReport{
		Date:                today,
		TotalOrders:        totalOrders,
		PendingOrders:      orderStats[models.OrderStatusPending],
		ProcessingOrders:   orderStats[models.OrderStatusProcessing],
		ShippedOrders:      orderStats[models.OrderStatusShipped],
		DeliveredOrders:    orderStats[models.OrderStatusDelivered],
		CancelledOrders:    orderStats[models.OrderStatusCancelled],
		TotalRevenue:       totalRevenue,
		AverageOrderValue:  totalRevenue / float64(totalOrders),
		UniqueCustomers:    totalCustomers,
		NewCustomers:       newCustomers,
		TopProducts:        topProducts,
		LowStockProducts:   lowStockAlerts,
		OrderFulfillmentRate: fulfillmentRate,
		CancellationRate:   cancellationRate,
	}

	// Save report
	if err := w.reportRepo.CreateReport(ctx, tx, report); err != nil {
		logger.Error(ctx, "Failed to create report", zap.Error(err))
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.Error(ctx, "Failed to commit transaction", zap.Error(err))
		return
	}

	logger.Info(ctx, "Successfully generated daily report", zap.String("date", today.Format("2006-01-02")))
}
