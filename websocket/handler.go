// Package websocket provides message handling for Bithumb WebSocket API.
package websocket

import (
	"encoding/json"
	"fmt"
	"log"
)

// handleMessage processes a WebSocket message and dispatches to the appropriate handler.
func (c *Client) handleMessage(data []byte) error {
	// Parse message to determine type
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parse message: %w", err)
	}

	msgType, ok := raw["type"].(string)
	if !ok {
		return fmt.Errorf("missing type field")
	}

	// Log for debugging
	log.Printf("[WebSocket] Message type: %s, looking for handler", msgType)

	// Find and call the appropriate handler
	c.mu.RLock()
	handler, exists := c.handlers[SubscriptionType(msgType)]
	c.mu.RUnlock()

	if !exists {
		log.Printf("[WebSocket] No handler for type: %s, available handlers: %v", msgType, c.getHandlerTypes())
		return nil // No handler registered for this type
	}

	return handler(data)
}

// getHandlerTypes returns the list of registered handler types for debugging.
func (c *Client) getHandlerTypes() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var types []string
	for t := range c.handlers {
		types = append(types, string(t))
	}
	return types
}
