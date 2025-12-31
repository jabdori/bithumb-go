// Package public provides a client for Bithumb Public API.
package public

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/bithumb-go/bithumb-go/internal/base"
	"github.com/bithumb-go/bithumb-go/models/public"
)

// Client provides access to Bithumb Public API.
type Client struct {
	base base.Client
}

// NewClient creates a new Public API client.
func NewClient(base base.Client) *Client {
	return &Client{base: base}
}

// buildURL constructs a URL with query parameters.
func (c *Client) buildURL(path string, params map[string]string) string {
	u := c.base.BaseURL() + path
	if len(params) > 0 {
		v := url.Values{}
		for k, val := range params {
			v.Add(k, val)
		}
		u += "?" + v.Encode()
	}
	return u
}

// do performs an HTTP request.
func (c *Client) do(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	return c.base.HTTPClient().Do(req)
}

// GetTicker retrieves current ticker information.
func (c *Client) GetTicker(req *public.GetTickerRequest) ([]public.Ticker, error) {
	return c.GetTickerWithContext(context.Background(), req)
}

// GetTickerWithContext retrieves current ticker information with context.
func (c *Client) GetTickerWithContext(ctx context.Context, req *public.GetTickerRequest) ([]public.Ticker, error) {
	params := map[string]string{}
	if len(req.Markets) > 0 {
		params["markets"] = req.Markets[0]
	}

	resp, err := c.do(ctx, http.MethodGet, c.buildURL("/v1/ticker", params), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var tickers []public.Ticker
	if err := json.Unmarshal(body, &tickers); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return tickers, nil
}

// GetOrderBook retrieves order book information.
func (c *Client) GetOrderBook(req *public.GetOrderBookRequest) (*public.OrderBook, error) {
	return c.GetOrderBookWithContext(context.Background(), req)
}

// GetOrderBookWithContext retrieves order book information with context.
func (c *Client) GetOrderBookWithContext(ctx context.Context, req *public.GetOrderBookRequest) (*public.OrderBook, error) {
	params := map[string]string{}
	if len(req.Markets) > 0 {
		params["markets"] = req.Markets[0]
	}

	resp, err := c.do(ctx, http.MethodGet, c.buildURL("/v1/orderbook", params), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var orderbooks []public.OrderBook
	if err := json.Unmarshal(body, &orderbooks); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if len(orderbooks) == 0 {
		return nil, fmt.Errorf("no orderbook data returned")
	}

	return &orderbooks[0], nil
}

// GetRecentTrades retrieves recent trade information.
func (c *Client) GetRecentTrades(req *public.GetRecentTradesRequest) ([]public.Trade, error) {
	return c.GetRecentTradesWithContext(context.Background(), req)
}

// GetRecentTradesWithContext retrieves recent trade information with context.
func (c *Client) GetRecentTradesWithContext(ctx context.Context, req *public.GetRecentTradesRequest) ([]public.Trade, error) {
	params := map[string]string{
		"market": req.Market,
	}
	if req.To != "" {
		params["to"] = req.To
	}
	if req.Count > 0 {
		params["count"] = fmt.Sprintf("%d", req.Count)
	}

	resp, err := c.do(ctx, http.MethodGet, c.buildURL("/v1/trades/ticks", params), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var trades []public.Trade
	if err := json.Unmarshal(body, &trades); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return trades, nil
}

// GetCandlestick retrieves candlestick data.
func (c *Client) GetCandlestick(req *public.GetCandlestickRequest) ([]public.Candle, error) {
	return c.GetCandlestickWithContext(context.Background(), req)
}

// GetCandlestickWithContext retrieves candlestick data with context.
func (c *Client) GetCandlestickWithContext(ctx context.Context, req *public.GetCandlestickRequest) ([]public.Candle, error) {
	params := map[string]string{"market": req.Market}
	if req.To != "" {
		params["to"] = req.To
	}
	if req.Count > 0 {
		params["count"] = fmt.Sprintf("%d", req.Count)
	}

	url := fmt.Sprintf("/v1/candles/minutes/%s", req.Unit)
	resp, err := c.do(ctx, http.MethodGet, c.buildURL(url, params), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var candles []public.Candle
	if err := json.Unmarshal(body, &candles); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return candles, nil
}
