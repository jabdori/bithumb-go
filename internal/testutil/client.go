package testutil

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/hysuki/bithumb-go/client"
)

// NewTestClient creates a client for testing purposes.
func NewTestClient(apiKey, apiSecret string) (*client.Client, error) {
	return client.NewClient(
		client.WithAPIKey(apiKey, apiSecret),
		client.WithHTTPClient(&http.Client{
			Timeout: 10 * time.Second,
		}),
	)
}

// NewMockClient creates a mock base client for testing.
type MockClient struct {
	BaseURLValue    string
	APIKeyValue     string
	APISecretValue  string
	HasAPIKeyValue  bool
	HTTPClientValue *http.Client
}

// NewMockClient creates a new mock client.
func NewMockClient() *MockClient {
	return &MockClient{
		BaseURLValue:   "https://test.bithumb.com",
		APIKeyValue:    "test-key",
		APISecretValue: "test-secret",
		HasAPIKeyValue: true,
		HTTPClientValue: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// BaseURL implements base.Client.
func (m *MockClient) BaseURL() string {
	return m.BaseURLValue
}

// HTTPClient implements base.Client.
func (m *MockClient) HTTPClient() *http.Client {
	return m.HTTPClientValue
}

// HasAPIKey implements base.Client.
func (m *MockClient) HasAPIKey() bool {
	return m.HasAPIKeyValue
}

// APIKey implements base.Client.
func (m *MockClient) APIKey() string {
	return m.APIKeyValue
}

// APISecret implements base.Client.
func (m *MockClient) APISecret() string {
	return m.APISecretValue
}

// NewTestServer creates a test HTTP server.
func NewTestServer(handler http.HandlerFunc) (*httptest.Server, *client.Client) {
	server := httptest.NewServer(handler)

	c, _ := client.NewClient(
		client.WithBaseURL(server.URL),
		client.WithHTTPClient(server.Client()),
	)

	return server, c
}
