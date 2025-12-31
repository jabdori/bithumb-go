// Package websocket provides data models for Bithumb WebSocket messages.
package websocket

// MyOrderMessage represents a user's order WebSocket message.
type MyOrderMessage struct {
	// Type is the message type.
	Type string `json:"type"`
	// OrderID is the unique order identifier.
	OrderID string `json:"uuid"`
	// Market is the market identifier.
	Market string `json:"market"`
	// Side is the order side ("bid" or "ask").
	Side string `json:"side"`
	// OrderType is the order type ("limit" or "market").
	OrderType string `json:"order_type"`
	// Price is the order price.
	Price string `json:"price"`
	// Volume is the order volume.
	Volume string `json:"volume"`
	// RemainingVolume is the remaining volume to be executed.
	RemainingVolume string `json:"remaining_volume"`
	// ReservedFee is the reserved fee.
	ReservedFee string `json:"reserved_fee"`
	// RemainingFee is the remaining fee.
	RemainingFee string `json:"remaining_fee"`
	// PaidFee is the paid fee.
	PaidFee string `json:"paid_fee"`
	// Locked is the locked amount.
	Locked string `json:"locked"`
	// ExecutedVolume is the executed volume.
	ExecutedVolume string `json:"executed_volume"`
	// TradesCount is the number of trades.
	TradesCount int64 `json:"trades_count"`
	// CreatedAt is the order creation timestamp.
	CreatedAt int64 `json:"created_at"`
	// Status is the order status.
	Status string `json:"status"`
}
