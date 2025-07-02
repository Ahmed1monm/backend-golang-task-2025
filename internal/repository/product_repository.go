package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product, inventory *models.Inventory) error
	FindByID(ctx context.Context, id uint) (*models.Product, error)
	List(ctx context.Context, offset, limit int) ([]models.Product, int64, error)
	GetInventory(ctx context.Context, productID uint) (*models.Inventory, error)
	Update(ctx context.Context, product *models.Product, inventory *models.Inventory) error
	GetTopProducts(ctx context.Context, tx *gorm.DB, date time.Time, limit int) ([]models.TopProduct, error)
	GetLowStockProducts(ctx context.Context, tx *gorm.DB) ([]models.LowStockAlert, error)
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

func (r *productRepository) GetInventory(ctx context.Context, productID uint) (*models.Inventory, error) {
	var inventory models.Inventory
	result := r.db.WithContext(ctx).Where("product_id = ?", productID).First(&inventory)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &inventory, nil
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

func (r *productRepository) Update(ctx context.Context, product *models.Product, inventory *models.Inventory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update product if provided
		if product != nil {
			if err := tx.Model(product).Updates(product).Error; err != nil {
				return err
			}
		}

		// Update inventory if provided
		if inventory != nil {
			if err := tx.Model(&models.Inventory{}).Where("product_id = ?", product.ID).Updates(map[string]interface{}{
				"quantity": inventory.Quantity,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *productRepository) GetTopProducts(ctx context.Context, tx *gorm.DB, date time.Time, limit int) ([]models.TopProduct, error) {
	var products []models.TopProduct

	err := tx.WithContext(ctx).
		Table("order_items").
		Select(
			"order_items.product_id,"+
				"products.name,"+
				"SUM(order_items.quantity) as total_quantity,"+
				"SUM(order_items.quantity * order_items.price) as total_revenue",
		).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("DATE(orders.created_at) = DATE(?)", date).
		Group("order_items.product_id, products.name").
		Order("total_quantity DESC").
		Limit(limit).
		Scan(&products).Error

	return products, err
}

func (r *productRepository) GetLowStockProducts(ctx context.Context, tx *gorm.DB) ([]models.LowStockAlert, error) {
	var alerts []models.LowStockAlert

	err := tx.WithContext(ctx).
		Table("inventories").
		Select(
			"inventories.product_id,"+
				"products.name,"+
				"inventories.quantity as current_stock,"+
				"products.reorder_threshold",
		).
		Joins("JOIN products ON products.id = inventories.product_id").
		Where("inventories.quantity <= products.reorder_threshold").
		Scan(&alerts).Error

	return alerts, err
}
