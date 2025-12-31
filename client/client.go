package client

import (
	"net/http"
	"time"
)

// Client represents the Bithumb API client.
type Client struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
	apiSecret  string

	HasAPIKeyFunc func() bool
}

// NewClient creates a new Bithumb API client.
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		baseURL:    "https://api.bithumb.com",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.HasAPIKeyFunc = func() bool {
		return c.apiKey != "" && c.apiSecret != ""
	}

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
