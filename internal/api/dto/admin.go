package dto

// UpdateOrderStatusRequest represents a request to update an order's status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending processing shipped delivered cancelled"`
}

// AdminOrderResponse represents an order response with admin-specific fields
type AdminOrderResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	Status      string              `json:"status"`
	Items       []OrderItemResponse `json:"items"`
	TotalAmount float64            `json:"total_amount"`
	CreatedAt   string             `json:"created_at,omitempty"`
	UpdatedAt   string             `json:"updated_at,omitempty"`
}
