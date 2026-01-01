// Package public provides a client for Bithumb Public API.
package public

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	bithumbgo "github.com/hysuki/bithumb-go"
	"github.com/hysuki/bithumb-go/internal/base"
	"github.com/hysuki/bithumb-go/internal/query"
	"github.com/hysuki/bithumb-go/models/public"
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

// GetWeekCandles retrieves week candles.
func (c *Client) GetWeekCandles(req *public.GetWeekCandlesRequest) ([]public.WeekCandle, error) {
	return c.GetWeekCandlesWithContext(context.Background(), req)
}

// GetWeekCandlesWithContext retrieves week candles with context.
func (c *Client) GetWeekCandlesWithContext(ctx context.Context, req *public.GetWeekCandlesRequest) ([]public.WeekCandle, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	params := map[string]string{"market": req.Market}
	if req.To != "" {
		params["to"] = req.To
	}
	if req.Count > 0 {
		params["count"] = fmt.Sprintf("%d", req.Count)
	}
	if req.ConvertingPriceUnit != "" {
		params["convertingPriceUnit"] = req.ConvertingPriceUnit
	}

	resp, err := c.do(ctx, http.MethodGet, c.buildURL("/v1/candles/weeks", params), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var candles []public.WeekCandle
	if err := json.Unmarshal(body, &candles); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return candles, nil
}

// GetMonthCandles retrieves month candles.
func (c *Client) GetMonthCandles(req *public.GetMonthCandlesRequest) ([]public.MonthCandle, error) {
	return c.GetMonthCandlesWithContext(context.Background(), req)
}

// GetMonthCandlesWithContext retrieves month candles with context.
func (c *Client) GetMonthCandlesWithContext(ctx context.Context, req *public.GetMonthCandlesRequest) ([]public.MonthCandle, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	params := map[string]string{"market": req.Market}
	if req.To != "" {
		params["to"] = req.To
	}
	if req.Count > 0 {
		params["count"] = fmt.Sprintf("%d", req.Count)
	}
	if req.ConvertingPriceUnit != "" {
		params["convertingPriceUnit"] = req.ConvertingPriceUnit
	}

	resp, err := c.do(ctx, http.MethodGet, c.buildURL("/v1/candles/months", params), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var candles []public.MonthCandle
	if err := json.Unmarshal(body, &candles); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return candles, nil
}

// GetMarketAll retrieves all available market codes.
func (c *Client) GetMarketAll(details bool) ([]public.Market, error) {
	return c.GetMarketAllWithContext(context.Background(), details)
}

// GetMarketAllWithContext retrieves all available market codes with context.
func (c *Client) GetMarketAllWithContext(ctx context.Context, details bool) ([]public.Market, error) {
	params := map[string]string{}
	if details {
		params["isDetails"] = "true"
	}

	resp, err := c.do(ctx, http.MethodGet, c.buildURL("/v1/market/all", params), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var markets []public.Market
	if err := json.Unmarshal(body, &markets); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return markets, nil
}

// GetDayCandles retrieves day candles.
func (c *Client) GetDayCandles(req *public.GetDayCandlesRequest) ([]public.DayCandle, error) {
	return c.GetDayCandlesWithContext(context.Background(), req)
}

// GetDayCandlesWithContext retrieves day candles with context.
func (c *Client) GetDayCandlesWithContext(ctx context.Context, req *public.GetDayCandlesRequest) ([]public.DayCandle, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	params := map[string]string{"market": req.Market}
	if req.To != "" {
		params["to"] = req.To
	}
	if req.Count > 0 {
		params["count"] = fmt.Sprintf("%d", req.Count)
	}
	if req.ConvertingPriceUnit != "" {
		params["convertingPriceUnit"] = req.ConvertingPriceUnit
	}

	resp, err := c.do(ctx, http.MethodGet, c.buildURL("/v1/candles/days", params), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var candles []public.DayCandle
	if err := json.Unmarshal(body, &candles); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return candles, nil
}

// GetWarnings retrieves market warning alerts.
func (c *Client) GetWarnings() ([]public.Warning, *bithumbgo.Error) {
	return c.GetWarningsWithContext(context.Background())
}

// GetWarningsWithContext retrieves market warning alerts with context.
func (c *Client) GetWarningsWithContext(ctx context.Context) ([]public.Warning, *bithumbgo.Error) {
	resp, err := c.do(ctx, http.MethodGet, c.buildURL("/v1/market/warning", nil), nil)
	if err != nil {
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeNetwork, Message: "HTTP request failed", Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeHTTP, Message: fmt.Sprintf("API error: status %d: %s", resp.StatusCode, string(body)), HTTPStatus: resp.StatusCode}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "read response failed", Err: err}
	}

	var warnings []public.Warning
	if err := json.Unmarshal(body, &warnings); err != nil {
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "parse response failed", Err: err}
	}

	return warnings, nil
}

// GetNotices retrieves Bithumb announcements
func (c *Client) GetNotices(req *public.GetNoticesRequest) ([]public.Notice, *bithumbgo.Error) {
	return c.GetNoticesWithContext(context.Background(), req)
}

// GetNoticesWithContext retrieves Bithumb announcements with context
func (c *Client) GetNoticesWithContext(ctx context.Context, req *public.GetNoticesRequest) ([]public.Notice, *bithumbgo.Error) {
	if err := req.Validate(); err != nil {
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeAPI, Message: "invalid request", Err: err}
	}

	params := query.New().AddInt("count", req.Count)
	url := c.base.BaseURL() + "/v1/notice?" + params.Encode()

	resp, err := c.do(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeNetwork, Message: "HTTP request failed", Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, &bithumbgo.Error{
			Type:       bithumbgo.ErrorTypeHTTP,
			Message:    fmt.Sprintf("API error: status %d: %s", resp.StatusCode, string(body)),
			HTTPStatus: resp.StatusCode,
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "read response failed", Err: err}
	}

	var notices []public.Notice
	if err := json.Unmarshal(body, &notices); err != nil {
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "parse response failed", Err: err}
	}

	return notices, nil
}
