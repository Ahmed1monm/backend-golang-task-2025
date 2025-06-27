package repository

import (
	"context"
	"errors"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product, inventory *models.Inventory) error
	FindByID(ctx context.Context, id uint) (*models.Product, error)
	List(ctx context.Context, offset, limit int) ([]models.Product, int64, error)
}

type productRepository struct {
	db *gorm.DB
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	result := r.db.WithContext(ctx).First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &product, nil
}

func (r *productRepository) List(ctx context.Context, offset, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated products
	result := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&products)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return products, total, nil
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product, inventory *models.Inventory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create product
		if err := tx.Create(product).Error; err != nil {
			return err
		}

		// Create inventory
		inventory.ProductID = product.ID
		if err := tx.Create(inventory).Error; err != nil {
			return err
		}

		return nil
	})
}
