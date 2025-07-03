package service

import (
	"context"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/api/dto"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/repository"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProduct(ctx context.Context, id uint) (*dto.ProductResponse, error)
	ListProducts(ctx context.Context, page, limit int) (*dto.PaginatedProductsResponse, error)
	GetInventory(ctx context.Context, productID uint) (*dto.InventoryResponse, error)
	UpdateProduct(ctx context.Context, id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
}

type productService struct {
	productRepo   repository.ProductRepository
	orderRepo     repository.OrderRepository
	inventoryRepo repository.InventoryRepository
	db            *gorm.DB
}

func (s *productService) GetProduct(ctx context.Context, id uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error(ctx, "Failed to get product", zap.Error(err))
		return nil, err
	}

	if product == nil {
		return nil, nil
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		StockLevel:  product.Quantity,
	}, nil
}

func (s *productService) ListProducts(ctx context.Context, page, limit int) (*dto.PaginatedProductsResponse, error) {
	offset := (page - 1) * limit

	products, total, err := s.productRepo.List(ctx, offset, limit)
	if err != nil {
		logger.Error(ctx, "Failed to list products", zap.Error(err))
		return nil, err
	}

	responseProducts := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		responseProducts[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			SKU:         product.SKU,
			StockLevel:  product.Quantity,
		}
	}

	return &dto.PaginatedProductsResponse{
		Products: responseProducts,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

func NewProductService(repo repository.ProductRepository, orderRepo repository.OrderRepository, inventoryRepo repository.InventoryRepository, db *gorm.DB) ProductService {
	return &productService{
		productRepo:   repo,
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
		db:            db,
	}
}

func (s *productService) GetInventory(ctx context.Context, productID uint) (*dto.InventoryResponse, error) {
	// Check if product exists first
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		logger.Error(ctx, "Failed to get product", zap.Error(err))
		return nil, err
	}
	if product == nil {
		return nil, nil
	}

	// Get inventory
	inventory, err := s.productRepo.GetInventory(ctx, productID)
	if err != nil {
		logger.Error(ctx, "Failed to get inventory", zap.Error(err))
		return nil, err
	}
	if inventory == nil {
		return nil, nil
	}

	productInfo, err := s.productRepo.FindByID(ctx, inventory.ProductID)
	if err != nil {
		logger.Error(ctx, "Failed to get product for inventory", zap.Error(err))
		return nil, err
}

return &dto.InventoryResponse{
		ProductID:    inventory.ProductID,
		SKU:          productInfo.SKU,
		StockLevel:   inventory.Quantity,
		MinimumStock: inventory.MinimumStock,
	}, nil
}

func (s *productService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// Create product model
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	// Create inventory model
	inventory := &models.Inventory{
		Quantity: req.Quantity,
	}

	// Create product with inventory
	if err := s.productRepo.Create(ctx, product, inventory); err != nil {
		logger.Error(ctx, "Failed to create product", zap.Error(err))
		return nil, err
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		StockLevel:  inventory.Quantity,
	}, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	// Get existing product
	existingProduct, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error(ctx, "Failed to get product", zap.Error(err))
		return nil, err
	}
	if existingProduct == nil {
		return nil, nil
	}

	// Update product fields if provided
	if req.Name != nil {
		existingProduct.Name = *req.Name
	}
	if req.Description != nil {
		existingProduct.Description = *req.Description
	}
	if req.Price != nil {
		existingProduct.Price = *req.Price
	}

	// Get and update inventory if quantity provided
	if req.Quantity != nil {
		// Start a transaction
		err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// Get current inventory with lock
			inventory, err := s.inventoryRepo.GetForUpdate(ctx, tx, id)
			if err != nil {
				return err
			}

			// Update inventory quantity
			inventory.Quantity = *req.Quantity

			// Update with transaction-safe method
			return s.inventoryRepo.Update(ctx, tx, inventory)
		})
		if err != nil {
			logger.Error(ctx, "Failed to update inventory", zap.Error(err))
			return nil, err
		}
	}

	// Update product
	if err := s.productRepo.Update(ctx, existingProduct, nil); err != nil {
		logger.Error(ctx, "Failed to update product", zap.Error(err))
		return nil, err
	}

	// Get updated inventory for response
	updatedInventory, err := s.productRepo.GetInventory(ctx, id)
	if err != nil {
		logger.Error(ctx, "Failed to get updated inventory", zap.Error(err))
		return nil, err
	}

	// Return response
	return &dto.ProductResponse{
		ID:          existingProduct.ID,
		Name:        existingProduct.Name,
		Description: existingProduct.Description,
		Price:       existingProduct.Price,
		SKU:         existingProduct.SKU,
		StockLevel:  updatedInventory.Quantity,
	}, nil
}
