// Package websocket provides message handling tests for Bithumb WebSocket API.
package websocket

import (
	"sync"
	"testing"
)

// mockHandler tracks whether Error was called
type mockHandler struct {
	mu          sync.Mutex
	errors      []WSError
	handleCalls int
}

func (m *mockHandler) Handle(data []byte) error {
	m.mu.Lock()
	m.handleCalls++
	m.mu.Unlock()
	return nil
}

func (m *mockHandler) Error(err WSError) {
	m.mu.Lock()
	m.errors = append(m.errors, err)
	m.mu.Unlock()
}

func TestHandleMessage_WithErrorResponse(t *testing.T) {
	client := NewClient(nil)

	errorData := []byte(`{
        "error": {
            "name": "NO_TICKET",
            "message": "티켓이 존재하지 않거나, 유효하지 않습니다."
        }
    }`)

	mock := &mockHandler{}
	client.handlers[SubscriptionTypeTicker] = mock

	err := client.handleMessage(errorData)
	if err != nil {
		t.Fatalf("handleMessage returned error: %v", err)
	}

	if len(mock.errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(mock.errors))
	}

	if mock.errors[0].Name != WSErrNoTicket {
		t.Errorf("Expected error name %s, got %s", WSErrNoTicket, mock.errors[0].Name)
	}
}

func TestHandleMessage_WithNormalMessage(t *testing.T) {
	client := NewClient(nil)

	normalData := []byte(`{"type":"ticker","content":"test"}`)

	mock := &mockHandler{}
	client.handlers[SubscriptionTypeTicker] = mock

	err := client.handleMessage(normalData)
	if err != nil {
		t.Fatalf("handleMessage returned error: %v", err)
	}

	if len(mock.errors) != 0 {
		t.Fatalf("Expected 0 errors, got %d", len(mock.errors))
	}

	if mock.handleCalls != 1 {
		t.Errorf("Expected 1 handle call, got %d", mock.handleCalls)
	}
}
