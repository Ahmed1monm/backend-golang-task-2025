package dto

// PaginationQuery represents query parameters for pagination
type PaginationQuery struct {
	Page     int `query:"page" validate:"gte=1"`
	Limit    int `query:"limit" validate:"gte=1,lte=100"`
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description string  `json:"description" validate:"required,min=10,max=1000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Quantity    int     `json:"quantity" validate:"required,gte=0"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string  `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Quantity    *int     `json:"quantity,omitempty" validate:"omitempty,gte=0"`
}

// ProductResponse represents a product in responses
type ProductResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type InventoryResponse struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
	Reserved  int  `json:"reserved"`
}

// ListProductsResponse represents the response for listing products
type ListProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}

// CreateProductResponse represents the response body for creating a product
type CreateProductResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CreatedAt   string  `json:"created_at"`
}
