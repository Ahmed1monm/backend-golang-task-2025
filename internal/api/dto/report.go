package dto

import "time"

// DailySalesReportResponse represents a daily sales report response
type DailySalesReportResponse struct {
	ID                  uint              `json:"id"`
	Date                time.Time         `json:"date"`
	TotalOrders         int               `json:"total_orders"`
	PendingOrders       int               `json:"pending_orders"`
	ProcessingOrders    int               `json:"processing_orders"`
	ShippedOrders       int               `json:"shipped_orders"`
	DeliveredOrders     int               `json:"delivered_orders"`
	CancelledOrders     int               `json:"cancelled_orders"`
	TotalRevenue        float64           `json:"total_revenue"`
	AverageOrderValue   float64           `json:"average_order_value"`
	UniqueCustomers     int               `json:"unique_customers"`
	NewCustomers        int               `json:"new_customers"`
	TopProducts         []TopProductDTO   `json:"top_products"`
	LowStockProducts    []LowStockAlert   `json:"low_stock_products"`
	OrderFulfillmentRate float64          `json:"order_fulfillment_rate"`
	CancellationRate    float64          `json:"cancellation_rate"`
}

// TopProductDTO represents a top product in the sales report
type TopProductDTO struct {
	ProductID     uint    `json:"product_id"`
	ProductName   string  `json:"product_name"`
	QuantitySold  int     `json:"quantity_sold"`
	Revenue       float64 `json:"revenue"`
	StockTurnover float64 `json:"stock_turnover"`
}

// LowStockAlert represents a low stock alert in the sales report
type LowStockAlert struct {
	ProductID      uint    `json:"product_id"`
	ProductName    string  `json:"product_name"`
	CurrentStock   int     `json:"current_stock"`
	ReservedStock  int     `json:"reserved_stock"`
	ReorderPoint   int     `json:"reorder_point"`
}
