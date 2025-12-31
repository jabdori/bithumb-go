// Package private provides JWT authentication for Bithumb Private API.
package private

import (
	"strings"
	"testing"
	"time"

	"github.com/bithumb-go/bithumb-go"
	"github.com/bithumb-go/bithumb-go/client"
	"github.com/golang-jwt/jwt/v5"
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

	// Verify it's our custom error type
	if !bithumbgo.IsAPIError(err) {
		t.Errorf("Expected APIError type, got %T", err)
	}
}

func TestGenerateToken_Claims(t *testing.T) {
	base, _ := client.NewClient(
		client.WithAPIKey("test-access-key", "test-secret"),
	)
	c := NewClient(base)

	tokenString, _ := c.GenerateToken()

	// Parse and verify claims
	token, _ := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})

	if claims, ok := token.Claims.(*TokenClaims); ok {
		if claims.AccessKey != "test-access-key" {
			t.Errorf("AccessKey = %v, want test-access-key", claims.AccessKey)
		}
		if claims.Nonce == "" {
			t.Error("Nonce should not be empty")
		}
		if claims.Timestamp == 0 {
			t.Error("Timestamp should be set")
		}
	} else {
		t.Error("Failed to parse token claims")
	}
}

func TestGenerateToken_Expiration(t *testing.T) {
	base, _ := client.NewClient(
		client.WithAPIKey("test-key", "test-secret"),
	)
	c := NewClient(base)

	tokenString, _ := c.GenerateToken()

	// Parse and verify expiration
	token, _ := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})

	if claims, ok := token.Claims.(*TokenClaims); ok {
		expectedExpiry := time.Now().Add(DefaultTokenExpiration)
		actualExpiry := claims.ExpiresAt.Time

		// Allow 1 second tolerance
		if actualExpiry.Sub(expectedExpiry) > time.Second {
			t.Errorf("Expiration = %v, want ~%v", actualExpiry, expectedExpiry)
		}
	}
}

func TestGenerateToken_NonceUniqueness(t *testing.T) {
	base, _ := client.NewClient(
		client.WithAPIKey("test-key", "test-secret"),
	)
	c := NewClient(base)

	token1, _ := c.GenerateToken()
	token2, _ := c.GenerateToken()

	// Parse tokens to extract nonces
	claims1 := &TokenClaims{}
	claims2 := &TokenClaims{}

	jwt.ParseWithClaims(token1, claims1, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	jwt.ParseWithClaims(token2, claims2, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})

	if claims1.Nonce == claims2.Nonce {
		t.Error("Nonces should be unique across token generations")
	}
}
