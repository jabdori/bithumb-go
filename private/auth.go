// Package private provides JWT authentication for Bithumb Private API.
package private

import (
	"time"

	"github.com/bithumb-go/bithumb-go"
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

// TokenClaims represents JWT token claims for Bithumb API authentication.
// It includes the access key, a unique nonce, and timestamp along with
// standard JWT registered claims like expiration and issued-at times.
type TokenClaims struct {
	// AccessKey is the API key for authentication.
	AccessKey string `json:"access_key"`
	// Nonce is a unique value to prevent replay attacks.
	Nonce string `json:"nonce"`
	// Timestamp is the Unix millisecond timestamp when the token was created.
	Timestamp int64 `json:"timestamp"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for API authentication.
// The token expires after 1 hour and includes the API access key,
// a unique nonce, and the current timestamp.
func (c *Client) GenerateToken() (string, error) {
	if !c.base.HasAPIKey() {
		return "", &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeAPI,
			Message: "API key and secret required",
		}
	}

	now := time.Now()
	claims := TokenClaims{
		AccessKey: c.base.APIKey(),
		Nonce:     generateNonce(),
		Timestamp: now.UnixMilli(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(DefaultTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.base.APISecret()))
}

// DefaultTokenExpiration is the default expiration time for JWT tokens.
const DefaultTokenExpiration = 1 * time.Hour

// generateNonce generates a unique nonce using UUID v4.
// The nonce prevents replay attacks by ensuring each token is unique.
func generateNonce() string {
	return uuid.New().String()
}
