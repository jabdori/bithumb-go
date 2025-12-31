// Package websocket provides subscription management for Bithumb WebSocket API.
package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
)

// SubscriptionParam represents a subscription parameter.
type SubscriptionParam struct {
	// Type is the subscription type.
	Type SubscriptionType `json:"type"`
	// Symbols are the market codes to subscribe to.
	Symbols []string `json:"symbols,omitempty"`
}

// SubscriptionInfo represents subscription information.
type SubscriptionInfo struct {
	// Type is the subscription type.
	Type SubscriptionType
	// Codes are the subscribed market codes.
	Codes []string
	// Ticket is the unique ticket for this subscription.
	Ticket string
	// CreatedAt is when the subscription was created.
	CreatedAt time.Time
	// IsActive indicates if the subscription is active.
	IsActive bool
}

// SubscriptionManager manages WebSocket subscriptions.
type SubscriptionManager struct {
	subscriptions map[string]*SubscriptionInfo
	mu            sync.RWMutex
}

// NewSubscriptionManager creates a new subscription manager.
func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions: make(map[string]*SubscriptionInfo),
	}
}

// generateTicket generates a unique ticket using UUID.
func (sm *SubscriptionManager) generateTicket() string {
	return uuid.New().String()
}

// CreateSubscriptionMessage creates a subscription message with the given parameters.
func (sm *SubscriptionManager) CreateSubscriptionMessage(params []*SubscriptionParam) ([]byte, string, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	ticket := sm.generateTicket()
	messages := []interface{}{
		map[string]string{"ticket": ticket},
	}

	for _, p := range params {
		messages = append(messages, p)

		// Store subscription info with composite key for uniqueness
		if len(p.Symbols) == 0 {
			// No symbols specified, use type as key
			key := string(p.Type)
			sm.subscriptions[key] = &SubscriptionInfo{
				Type:      p.Type,
				Codes:     p.Symbols,
				Ticket:    ticket,
				CreatedAt: time.Now(),
				IsActive:  true,
			}
		} else {
			// Use type:code composite keys for each code
			for _, code := range p.Symbols {
				key := string(p.Type) + ":" + code
				sm.subscriptions[key] = &SubscriptionInfo{
					Type:      p.Type,
					Codes:     []string{code},
					Ticket:    ticket,
					CreatedAt: time.Now(),
					IsActive:  true,
				}
			}
		}
	}

	body, err := json.Marshal(messages)
	return body, ticket, err
}

// CreateUnsubscribeMessage creates a message to unsubscribe from all subscriptions.
func (sm *SubscriptionManager) CreateUnsubscribeMessage() ([]byte, string, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Deactivate all subscriptions
	for _, sub := range sm.subscriptions {
		sub.IsActive = false
	}

	ticket := sm.generateTicket()
	message := []interface{}{
		map[string]string{"ticket": ticket},
	}
	body, err := json.Marshal(message)
	return body, ticket, err
}

// RestoreSubscriptions creates a message to restore all active subscriptions.
func (sm *SubscriptionManager) RestoreSubscriptions() ([]byte, string, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	ticket := sm.generateTicket()
	messages := []interface{}{
		map[string]string{"ticket": ticket},
	}

	for _, sub := range sm.subscriptions {
		if sub.IsActive {
			param := &SubscriptionParam{
				Type:    sub.Type,
				Symbols: sub.Codes,
			}
			messages = append(messages, param)
		}
	}

	body, err := json.Marshal(messages)
	return body, ticket, err
}

// GetSubscriptions returns a copy of all subscriptions.
func (sm *SubscriptionManager) GetSubscriptions() []*SubscriptionInfo {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	result := make([]*SubscriptionInfo, 0, len(sm.subscriptions))
	for _, sub := range sm.subscriptions {
		// Return a copy to prevent modification of internal state
		copy := *sub
		result = append(result, &copy)
	}
	return result
}
