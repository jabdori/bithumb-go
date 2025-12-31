package base

import (
	"net/http"
)

// Client defines the interface for base client functionality.
// This interface is used to avoid circular dependencies between packages.
type Client interface {
	// BaseURL returns the configured base URL.
	BaseURL() string

	// HTTPClient returns the HTTP client.
	HTTPClient() *http.Client

	// HasAPIKey returns true if API key and secret are configured.
	HasAPIKey() bool

	// APIKey returns the configured API key.
	APIKey() string

	// APISecret returns the configured API secret.
	APISecret() string
}
