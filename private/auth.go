// Package private provides JWT authentication for Bithumb Private API.
package private

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"

	"github.com/bithumb-go/bithumb-go/client"
)

// Client provides JWT authentication for Private API.
type Client struct {
	base *client.Client
}

// NewClient creates a new Private API client.
func NewClient(base *client.Client) *Client {
	return &Client{base: base}
}

// TokenClaims represents JWT token claims for Bithumb API.
type TokenClaims struct {
	AccessKey string `json:"access_key"`
	Nonce     string `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for API authentication.
func (c *Client) GenerateToken() (string, error) {
	if !c.base.HasAPIKey() {
		return "", fmt.Errorf("API key and secret required")
	}

	now := time.Now()
	claims := TokenClaims{
		AccessKey: c.base.APIKey(),
		Nonce:     generateNonce(),
		Timestamp: now.UnixMilli(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.base.APISecret()))
}

// generateNonce generates a unique nonce using UUID.
func generateNonce() string {
	return uuid.New().String()
}
