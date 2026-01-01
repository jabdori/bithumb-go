// Package websocket provides message handling for Bithumb WebSocket API.
package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/hysuki/bithumb-go/logger"
)

// handleMessage processes a WebSocket message and dispatches to the appropriate handler.
func (c *Client) handleMessage(data []byte) error {
	// Parse message to determine type
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parse message: %w", err)
	}

	// Check for error response
	if errData, hasError := raw["error"]; hasError {
		c.logger.Error("WebSocket error response", logger.F("data", string(data)))

		// Parse error and call Error() callback on all handlers
		var wsErrResp WSErrorResponse
		if err := json.Unmarshal(data, &wsErrResp); err == nil {
			c.mu.RLock()
			for _, handler := range c.handlers {
				handler.Error(wsErrResp.Error)
			}
			c.mu.RUnlock()
		} else {
			// If we can't parse the error properly, log the raw data
			c.logger.Error("Failed to parse WebSocket error", logger.F("error", err), logger.F("raw_error", errData))
		}
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

	return handler.Handle(data)
}
