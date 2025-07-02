package repository

import (
	"context"
	"time"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"gorm.io/gorm"
)

type ReportRepository interface {
	CreateReport(ctx context.Context, tx *gorm.DB, report *models.DailySalesReport) error
	GetReportByDate(ctx context.Context, tx *gorm.DB, date time.Time) (*models.DailySalesReport, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) CreateReport(ctx context.Context, tx *gorm.DB, report *models.DailySalesReport) error {
	return tx.WithContext(ctx).Create(report).Error
}

func (r *reportRepository) GetReportByDate(ctx context.Context, tx *gorm.DB, date time.Time) (*models.DailySalesReport, error) {
	var report models.DailySalesReport
	err := tx.WithContext(ctx).
		Preload("TopProducts").
		Preload("LowStockProducts").
		Where("DATE(date) = DATE(?)", date).
		First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}


