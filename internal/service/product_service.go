package service

import (
	"context"
	"time"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/api/dto"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/models"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/repository"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"go.uber.org/zap"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.CreateProductResponse, error)
	GetProduct(ctx context.Context, id uint) (*dto.ProductResponse, error)
	ListProducts(ctx context.Context, page, limit int) (*dto.ListProductsResponse, error)
}

type productService struct {
	productRepo repository.ProductRepository
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
		Quantity:    product.Quantity,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *productService) ListProducts(ctx context.Context, page, limit int) (*dto.ListProductsResponse, error) {
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
			Quantity:    product.Quantity,
			CreatedAt:   product.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &dto.ListProductsResponse{
		Products: responseProducts,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: repo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.CreateProductResponse, error) {
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

	return &dto.CreateProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    inventory.Quantity,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
