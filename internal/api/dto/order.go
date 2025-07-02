package dto

type CreateOrderItemRequest struct {
	ProductID uint    `json:"product_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type OrderResponse struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
	Items       []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
