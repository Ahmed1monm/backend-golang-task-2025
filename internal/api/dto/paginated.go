package dto

// PaginatedOrdersResponse represents a paginated list of orders
type PaginatedOrdersResponse struct {
	Orders     []OrderResponse `json:"orders"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
}
