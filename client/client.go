package client

import (
	"net/http"
	"time"

	"github.com/hysuki/bithumb-go/private"
	"github.com/hysuki/bithumb-go/public"
	"github.com/hysuki/bithumb-go/websocket"
)

const (
	// DefaultBaseURL is the default base URL for the Bithumb API.
	DefaultBaseURL = "https://api.bithumb.com"
	// DefaultTimeout is the default HTTP request timeout.
	DefaultTimeout = 30 * time.Second
)

// Client represents the Bithumb API client.
// Client values are safe for concurrent use after initialization.
// The underlying HTTP client is safe for concurrent use by default.
type Client struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
	apiSecret  string

	// HasAPIKeyFunc is a function that checks if API key is configured.
	HasAPIKeyFunc func() bool

	// Public provides access to Public API endpoints.
	Public *public.Client

	// Private provides access to Private API endpoints (requires API key).
	// Will be nil if no API credentials were provided to NewClient.
	Private *private.Client

	// Websocket provides access to WebSocket connections.
	Websocket *websocket.Client
}

// NewClient creates a new Bithumb API client.
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		baseURL:    DefaultBaseURL,
		httpClient: &http.Client{Timeout: DefaultTimeout},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.HasAPIKeyFunc = func() bool {
		return c.apiKey != "" && c.apiSecret != ""
	}

	// Initialize Public API client
	c.Public = public.NewClient(c)

	// Initialize Private API client (only if API key is configured)
	if c.HasAPIKey() {
		c.Private = private.NewClient(c)
	}

	// Initialize WebSocket client
	c.Websocket = websocket.NewClient(c)

	return c, nil
}

// Option is a function that configures a Client.
type Option func(*Client)

// WithAPIKey sets the API key and secret for authentication.
func WithAPIKey(apiKey, apiSecret string) Option {
	return func(c *Client) {
		c.apiKey = apiKey
		c.apiSecret = apiSecret
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithBaseURL sets the base URL for API requests.
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// HasAPIKey returns true if API key and secret are configured.
func (c *Client) HasAPIKey() bool {
	return c.HasAPIKeyFunc()
}

// BaseURL returns the configured base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// HTTPClient returns the HTTP client.
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// APIKey returns the configured API key.
func (c *Client) APIKey() string {
	return c.apiKey
}

// APISecret returns the configured API secret.
func (c *Client) APISecret() string {
	return c.apiSecret
}
