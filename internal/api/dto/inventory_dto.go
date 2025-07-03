package dto

// LowStockAlertResponse represents a product with low stock levels
type LowStockAlertResponse struct {
	ProductID   uint    `json:"product_id"`
	Name        string  `json:"name"`
	SKU         string  `json:"sku"`
	StockLevel  int     `json:"stock_level"`
	MinimumStock int    `json:"minimum_stock"`
	Price       float64 `json:"price"`
}
