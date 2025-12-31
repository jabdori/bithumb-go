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

	// Check for error response
	if _, hasError := raw["error"]; hasError {
		log.Printf("[WebSocket] Error response: %s", string(data))
		return nil // Error responses are logged but don't stop processing
	}

	msgType, ok := raw["type"].(string)
	if !ok {
		return fmt.Errorf("missing type field")
	}

	// Find and call the appropriate handler
	c.mu.RLock()
	handler, exists := c.handlers[SubscriptionType(msgType)]
	c.mu.RUnlock()

	if !exists {
		return nil // No handler registered for this type
	}

	return handler(data)
}
