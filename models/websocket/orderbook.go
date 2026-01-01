// Package websocket provides data models for Bithumb WebSocket API messages.
package websocket

// OrderBookMessage represents an orderbook WebSocket message.
type OrderBookMessage struct {
	// Type is the message type ("orderbook").
	Type string `json:"type"`
	// Code is the market code (e.g., "KRW-BTC").
	Code string `json:"code"`
	// TotalAskSize is the total ask size.
	TotalAskSize float64 `json:"total_ask_size"`
	// TotalBidSize is the total bid size.
	TotalBidSize float64 `json:"total_bid_size"`
	// OrderBookUnits contains the orderbook units.
	OrderBookUnits []OrderBookUnit `json:"orderbook_units"`
	// Level is the orderbook level.
	Level int `json:"level"`
	// StreamType is the stream type ("SNAPSHOT" or "REALTIME").
	StreamType string `json:"stream_type"`
	// Timestamp is the message timestamp (milliseconds).
	Timestamp int64 `json:"timestamp"`
}

// OrderBookUnit represents a single orderbook unit.
type OrderBookUnit struct {
	// AskPrice is the ask price.
	AskPrice float64 `json:"ask_price"`
	// BidPrice is the bid price.
	BidPrice float64 `json:"bid_price"`
	// AskSize is the ask size.
	AskSize float64 `json:"ask_size"`
	// BidSize is the bid size.
	BidSize float64 `json:"bid_size"`
}
