package websocket

// Event types
const (
	EventOrderCreated      EventType = "order_created"
	EventOrderCancelled    EventType = "order_cancelled"
	EventInventoryUpdated  EventType = "inventory_updated"
)

// OrderEventPayload represents the payload for order-related events
type OrderEventPayload struct {
	OrderID     uint    `json:"order_id"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_amount"`
}

// InventoryEventPayload represents the payload for inventory-related events
type InventoryEventPayload struct {
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Name      string  `json:"name"`
}
