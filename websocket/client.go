// Package websocket provides a WebSocket client for Bithumb API.
package websocket

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/coder/websocket"

	"github.com/hysuki/bithumb-go/internal/base"
)

const (
	// DefaultPublicURL is the default WebSocket URL for public data.
	DefaultPublicURL = "wss://ws-api.bithumb.com/websocket/v1"
	// DefaultPrivateURL is the default WebSocket URL for private data.
	DefaultPrivateURL = "wss://ws-api.bithumb.com/websocket/v1/private"
	// DefaultReconnectDelay is the default delay between reconnection attempts.
	DefaultReconnectDelay = 5 * time.Second
	// DefaultReconnectTimeout is the default timeout for reconnection attempts.
	DefaultReconnectTimeout = 10 * time.Second
)

// Client represents a WebSocket client.
type Client struct {
	base             base.Client
	conn             *websocket.Conn
	url              string
	subs             *SubscriptionManager
	handlers         map[SubscriptionType]MessageHandler
	done             chan struct{}
	mu               sync.RWMutex
	reconnect        bool
	reconnectDelay   time.Duration
	reconnectTimeout time.Duration
	isConnected      bool
}

// NewClient creates a new WebSocket client.
func NewClient(base base.Client) *Client {
	return &Client{
		base:             base,
		url:              DefaultPublicURL,
		subs:             NewSubscriptionManager(),
		handlers:         make(map[SubscriptionType]MessageHandler),
		done:             make(chan struct{}),
		reconnect:        true,
		reconnectDelay:   DefaultReconnectDelay,
		reconnectTimeout: DefaultReconnectTimeout,
	}
}

// SetPrivateURL sets the URL for private WebSocket connections.
func (c *Client) SetPrivateURL() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.url = DefaultPrivateURL
}

// SetReconnect enables or disables automatic reconnection.
func (c *Client) SetReconnect(enabled bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.reconnect = enabled
}

// SetReconnectDelay sets the delay between reconnection attempts.
func (c *Client) SetReconnectDelay(delay time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.reconnectDelay = delay
}

// SetReconnectTimeout sets the timeout for reconnection attempts.
func (c *Client) SetReconnectTimeout(timeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.reconnectTimeout = timeout
}

// Connect establishes a WebSocket connection.
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected {
		return fmt.Errorf("already connected")
	}

	conn, _, err := websocket.Dial(ctx, c.url, nil)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	c.conn = conn
	c.isConnected = true

	// Start read and reconnect loops
	go c.readLoop()
	go c.reconnectLoop()

	return nil
}

// readLoop reads messages from the WebSocket connection.
func (c *Client) readLoop() {
	// Create context that respects c.done
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-c.done
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		default:
		}

		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()

		if conn == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Read message from WebSocket
		messageType, message, err := conn.Read(ctx)
		if err != nil {
			log.Printf("[WebSocket] Read error: %v", err)
			c.mu.Lock()
			c.isConnected = false
			c.mu.Unlock()
			return
		}

		// Process text messages (Bithumb sends JSON as binary type 2)
		if messageType == websocket.MessageText || messageType == websocket.MessageBinary {
			// Log received messages for debugging
			log.Printf("[WebSocket] Received: %s", string(message))

			if err := c.handleMessage(message); err != nil {
				// Log handler error but continue processing other messages
				log.Printf("[WebSocket] handler error: %v", err)
			}
		} else {
			log.Printf("[WebSocket] Received unhandled message type: %d", messageType)
		}
	}
}

// reconnectLoop handles automatic reconnection.
func (c *Client) reconnectLoop() {
	ticker := time.NewTicker(c.reconnectDelay)
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			c.mu.RLock()
			connected := c.isConnected
			shouldReconnect := c.reconnect
			timeout := c.reconnectTimeout
			c.mu.RUnlock()

			if !connected && shouldReconnect {
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				if err := c.Connect(ctx); err == nil {
					// Restore subscriptions after reconnection
					c.RestoreSubscriptions()
				} else {
					log.Printf("[WebSocket] reconnect failed: %v", err)
				}
				cancel()
			}
		}
	}
}

// Subscribe subscribes to WebSocket data with the given parameters and handlers.
func (c *Client) Subscribe(params []*SubscriptionParam, handlers MessageHandlers) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Register handlers
	if handlers.Ticker != nil {
		c.handlers[SubscriptionTypeTicker] = handlers.Ticker
	}
	if handlers.OrderBook != nil {
		c.handlers[SubscriptionTypeOrderBook] = handlers.OrderBook
	}
	if handlers.Trade != nil {
		c.handlers[SubscriptionTypeTrade] = handlers.Trade
	}
	if handlers.MyOrder != nil {
		c.handlers[SubscriptionTypeMyOrder] = handlers.MyOrder
	}
	if handlers.MyAsset != nil {
		c.handlers[SubscriptionTypeMyAsset] = handlers.MyAsset
	}

	// Create and send subscription message
	body, _, err := c.subs.CreateSubscriptionMessage(params)
	if err != nil {
		return fmt.Errorf("create subscription message: %w", err)
	}

	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	// Log the subscription message for debugging
	log.Printf("[WebSocket] Sending subscription: %s", string(body))

	err = c.conn.Write(context.Background(), websocket.MessageText, body)
	if err != nil {
		return fmt.Errorf("send subscription: %w", err)
	}

	return nil
}

// Unsubscribe unsubscribes from all subscriptions.
func (c *Client) Unsubscribe() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	body, _, err := c.subs.CreateUnsubscribeMessage()
	if err != nil {
		return fmt.Errorf("create unsubscribe message: %w", err)
	}

	if c.conn != nil {
		err = c.conn.Write(context.Background(), websocket.MessageText, body)
		if err != nil {
			return fmt.Errorf("send unsubscribe: %w", err)
		}
	}

	return nil
}

// RestoreSubscriptions restores all active subscriptions.
func (c *Client) RestoreSubscriptions() error {
	body, _, err := c.subs.RestoreSubscriptions()
	if err != nil {
		return fmt.Errorf("create restore message: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	return c.conn.Write(context.Background(), websocket.MessageText, body)
}

// Close closes the WebSocket connection.
func (c *Client) Close() error {
	close(c.done)

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		err := c.conn.Close(websocket.StatusNormalClosure, "client closing")
		c.conn = nil
		c.isConnected = false
		return err
	}

	return nil
}
