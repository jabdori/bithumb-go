package client

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.HasAPIKey() {
		t.Error("NewClient() should not have API key by default")
	}
}

func TestNewClientWithOptions(t *testing.T) {
	client, err := NewClient(
		WithAPIKey("test-key", "test-secret"),
		WithTimeout(10*time.Second),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if !client.HasAPIKey() {
		t.Error("NewClient(WithAPIKey()) should have API key")
	}
}
