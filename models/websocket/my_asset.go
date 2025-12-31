// Package websocket provides data models for Bithumb WebSocket messages.
package websocket

// MyAssetMessage represents a user's asset WebSocket message.
type MyAssetMessage struct {
	// Type is the message type.
	Type string `json:"type"`
	// Currency is the currency code.
	Currency string `json:"currency"`
	// Balance is the available balance.
	Balance string `json:"balance"`
	// Locked is the locked amount.
	Locked string `json:"locked"`
	// AvgBuyPrice is the average buy price.
	AvgBuyPrice string `json:"avg_buy_price"`
	// AvgBuyPriceModified is whether avg buy price was modified.
	AvgBuyPriceModified bool `json:"avg_buy_price_modified"`
	// UnitCurrency is the unit currency.
	UnitCurrency string `json:"unit_currency"`
	// Timestamp is the message timestamp.
	Timestamp int64 `json:"timestamp"`
}
