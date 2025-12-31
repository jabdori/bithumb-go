// Package private provides JWT authentication for Bithumb Private API.
package private

import (
	"strings"
	"testing"

	"github.com/bithumb-go/bithumb-go/client"
)

func TestGenerateToken(t *testing.T) {
	base, _ := client.NewClient(
		client.WithAPIKey("test-key", "test-secret"),
	)
	c := NewClient(base)

	token, err := c.GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if token == "" {
		t.Fatal("GenerateToken() returned empty token")
	}

	// Token should be JWT format (3 parts separated by dots)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("Token should have 3 parts, got %d", len(parts))
	}
}

func TestGenerateToken_NoAPIKey(t *testing.T) {
	base, _ := client.NewClient()
	c := NewClient(base)

	_, err := c.GenerateToken()
	if err == nil {
		t.Error("GenerateToken() should return error without API key")
	}
}
