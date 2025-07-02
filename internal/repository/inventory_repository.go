package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
)

type InventoryRepository interface {
	GetForUpdate(ctx context.Context, tx *gorm.DB, productID uint) (*models.Inventory, error)
	Update(ctx context.Context, tx *gorm.DB, inventory *models.Inventory) error
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) GetForUpdate(ctx context.Context, tx *gorm.DB, productID uint) (*models.Inventory, error) {
	var inventory models.Inventory
	err := tx.WithContext(ctx).Where("product_id = ?", productID).First(&inventory).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) Update(ctx context.Context, tx *gorm.DB, inventory *models.Inventory) error {
	// Use FOR UPDATE clause to prevent race conditions
	return tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Save(inventory).Error
}
