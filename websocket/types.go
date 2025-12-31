// Package websocket provides types and handlers for Bithumb WebSocket API.
//
// The package supports both public and private WebSocket subscriptions for
// real-time market data and account updates.
//
// Basic usage:
//
//	manager := websocket.NewSubscriptionManager()
//	params := []*websocket.SubscriptionParam{
//	    {Type: websocket.SubscriptionTypeTicker, Codes: []string{"KRW-BTC"}},
//	}
//	msg, ticket, err := manager.CreateSubscriptionMessage(params)
//
// Thread Safety: The SubscriptionManager is safe for concurrent use.
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
