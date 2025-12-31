// Package websocket provides data models for Bithumb WebSocket API messages.
package websocket

// TickerMessage represents a ticker WebSocket message.
type TickerMessage struct {
	// Type is the message type.
	Type string `json:"type"`
	// Content contains the ticker data.
	Content TickerContent `json:"content"`
}

// TickerContent contains ticker data.
// TODO: Add remaining fields based on Bithumb WebSocket API documentation.
// Current fields are placeholders for initial structure.
type TickerContent struct {
	// StreamType is the type of stream.
	StreamType string `json:"stream_type"`
	// MarketCode is the market code (e.g., "KRW-BTC").
	MarketCode string `json:"market_code"`
}

// OrderBookMessage represents an orderbook WebSocket message.
type OrderBookMessage struct {
	// Type is the message type.
	Type string `json:"type"`
	// Content contains the orderbook data.
	Content OrderBookContent `json:"content"`
}

// OrderBookContent contains orderbook data.
// TODO: Add remaining fields based on Bithumb WebSocket API documentation.
// Current fields are placeholders for initial structure.
type OrderBookContent struct {
	// StreamType is the type of stream.
	StreamType string `json:"stream_type"`
	// MarketCode is the market code (e.g., "KRW-BTC").
	MarketCode string `json:"market_code"`
}

// TradeMessage represents a trade WebSocket message.
type TradeMessage struct {
	// Type is the message type.
	Type string `json:"type"`
	// Content contains the trade data.
	Content TradeContent `json:"content"`
}

// TradeContent contains trade data.
// TODO: Add remaining fields based on Bithumb WebSocket API documentation.
// Current fields are placeholders for initial structure.
type TradeContent struct {
	// StreamType is the type of stream.
	StreamType string `json:"stream_type"`
	// MarketCode is the market code (e.g., "KRW-BTC").
	MarketCode string `json:"market_code"`
}
