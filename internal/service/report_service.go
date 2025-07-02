package service

import (
	"context"
	"time"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/repository"
)

type ReportService struct {
	reportRepo repository.ReportRepository
}

func NewReportService(reportRepo repository.ReportRepository) *ReportService {
	return &ReportService{
		reportRepo: reportRepo,
	}
}

// GetDailyReport returns the sales report for the specified date.
// If no report exists, returns an empty report structure.
func (s *ReportService) GetDailyReport(ctx context.Context, date time.Time) (*models.DailySalesReport, error) {
	// Get report from repository
	report, err := s.reportRepo.GetReportByDate(ctx, nil, date)
	if err != nil {
		return nil, err
	}

	// If no report exists yet [this will mean that the CRON job not finished yet], return empty report structure
	if report == nil {
		report = &models.DailySalesReport{
			Date:                 date,
			TotalOrders:          0,
			PendingOrders:        0,
			ProcessingOrders:     0,
			ShippedOrders:        0,
			DeliveredOrders:      0,
			CancelledOrders:      0,
			TotalRevenue:         0,
			AverageOrderValue:    0,
			UniqueCustomers:      0,
			NewCustomers:         0,
			TopProducts:          []models.TopProduct{},
			LowStockProducts:     []models.LowStockAlert{},
			OrderFulfillmentRate: 0,
			CancellationRate:     0,
		}
	}

	return report, nil
}
