// Package websocket provides types and handlers for Bithumb WebSocket API.
package websocket

// SubscriptionType represents a WebSocket subscription type.
type SubscriptionType string

const (
	// SubscriptionTypeTicker subscribes to ticker updates.
	SubscriptionTypeTicker SubscriptionType = "ticker"
	// SubscriptionTypeOrderBook subscribes to orderbook updates.
	SubscriptionTypeOrderBook SubscriptionType = "orderbook"
	// SubscriptionTypeTrade subscribes to trade updates.
	SubscriptionTypeTrade SubscriptionType = "transaction"
	// SubscriptionTypeMyOrder subscribes to user's order updates (Private).
	SubscriptionTypeMyOrder SubscriptionType = "myOrder"
	// SubscriptionTypeMyAsset subscribes to user's asset updates (Private).
	SubscriptionTypeMyAsset SubscriptionType = "myAsset"
)

// MessageHandler handles WebSocket messages.
type MessageHandler func(msg []byte) error

// MessageHandlers holds handlers for different subscription types.
type MessageHandlers struct {
	// Ticker handles ticker messages.
	Ticker MessageHandler
	// OrderBook handles orderbook messages.
	OrderBook MessageHandler
	// Trade handles trade messages.
	Trade MessageHandler
	// MyOrder handles user's order messages (Private).
	MyOrder MessageHandler
	// MyAsset handles user's asset messages (Private).
	MyAsset MessageHandler
}

// TickerMessage represents a ticker WebSocket message.
type TickerMessage struct {
	Type    string        `json:"type"`
	Content TickerContent `json:"content"`
}

// TickerContent contains ticker data.
type TickerContent struct {
	// Add fields based on Bithumb WebSocket API documentation
	// This is a placeholder structure
	AType string `json:"type"`
	Code  string `json:"code"`
	// Additional fields would be added based on actual API response
}

// OrderBookMessage represents an orderbook WebSocket message.
type OrderBookMessage struct {
	Type    string            `json:"type"`
	Content OrderBookContent `json:"content"`
}

// OrderBookContent contains orderbook data.
type OrderBookContent struct {
	AType string `json:"type"`
	Code  string `json:"code"`
	// Additional fields would be added based on actual API response
}

// TradeMessage represents a trade WebSocket message.
type TradeMessage struct {
	Type    string       `json:"type"`
	Content TradeContent `json:"content"`
}

// TradeContent contains trade data.
type TradeContent struct {
	AType string `json:"type"`
	Code  string `json:"code"`
	// Additional fields would be added based on actual API response
}
