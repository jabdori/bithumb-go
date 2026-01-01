// Package private provides a client for Bithumb Private API.
package private

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	bithumbgo "github.com/hysuki/bithumb-go"
	"github.com/hysuki/bithumb-go/internal/query"
	"github.com/hysuki/bithumb-go/models/private"
)

// doWithAuth performs an authenticated HTTP request.
func (c *Client) doWithAuth(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	token, err := c.GenerateToken()
	if err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeAPI,
			Message: fmt.Sprintf("generate token: %v", err),
			Err:     err,
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeNetwork,
			Message: fmt.Sprintf("create request: %v", err),
			Err:     err,
		}
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.base.HTTPClient().Do(req)
	if err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeNetwork,
			Message: "HTTP request failed",
			Err:     err,
		}
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, &bithumbgo.Error{
				Type:       bithumbgo.ErrorTypeHTTP,
				Message:    fmt.Sprintf("API error: status %d (failed to read body)", resp.StatusCode),
				HTTPStatus: resp.StatusCode,
				Err:        err,
			}
		}
		return nil, &bithumbgo.Error{
			Type:       bithumbgo.ErrorTypeHTTP,
			Message:    fmt.Sprintf("API error: status %d: %s", resp.StatusCode, string(body)),
			HTTPStatus: resp.StatusCode,
		}
	}

	return resp, nil
}

// GetAccount retrieves account information.
func (c *Client) GetAccount(req *private.GetAccountRequest) ([]private.Account, error) {
	return c.GetAccountWithContext(context.Background(), req)
}

// GetAccountWithContext retrieves account information with context.
func (c *Client) GetAccountWithContext(ctx context.Context, req *private.GetAccountRequest) ([]private.Account, error) {
	if err := req.Validate(); err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeAPI,
			Message: fmt.Sprintf("invalid request: %v", err),
			Err:     err,
		}
	}

	url := c.base.BaseURL() + "/v1/accounts"

	resp, err := c.doWithAuth(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeParse,
			Message: "read response failed",
			Err:     err,
		}
	}

	var accounts []private.Account
	if err := json.Unmarshal(body, &accounts); err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeParse,
			Message: "parse response failed",
			Err:     err,
		}
	}

	return accounts, nil
}

// PlaceOrder places a new order.
func (c *Client) PlaceOrder(req *private.PlaceOrderRequest) (*private.Order, error) {
	return c.PlaceOrderWithContext(context.Background(), req)
}

// PlaceOrderWithContext places a new order with context.
func (c *Client) PlaceOrderWithContext(ctx context.Context, req *private.PlaceOrderRequest) (*private.Order, error) {
	if err := req.Validate(); err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeAPI,
			Message: fmt.Sprintf("invalid request: %v", err),
			Err:     err,
		}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeParse,
			Message: "marshal request failed",
			Err:     err,
		}
	}

	url := c.base.BaseURL() + "/v2/orders"

	resp, err := c.doWithAuth(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeParse,
			Message: "read response failed",
			Err:     err,
		}
	}

	var order private.Order
	if err := json.Unmarshal(respBody, &order); err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeParse,
			Message: "parse response failed",
			Err:     err,
		}
	}

	return &order, nil
}

// CancelOrder cancels an order.
func (c *Client) CancelOrder(req *private.CancelOrderRequest) error {
	return c.CancelOrderWithContext(context.Background(), req)
}

// CancelOrderWithContext cancels an order with context.
func (c *Client) CancelOrderWithContext(ctx context.Context, req *private.CancelOrderRequest) error {
	if err := req.Validate(); err != nil {
		return &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeAPI,
			Message: fmt.Sprintf("invalid request: %v", err),
			Err:     err,
		}
	}

	url := fmt.Sprintf("%s/v2/order/%s", c.base.BaseURL(), req.UUID)

	reqBody, err := json.Marshal(map[string]string{})
	if err != nil {
		return &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeParse,
			Message: "marshal request failed",
			Err:     err,
		}
	}

	resp, err := c.doWithAuth(ctx, http.MethodDelete, url, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Drain the response body to allow connection reuse
	io.Copy(io.Discard, resp.Body)

	return nil
}

// GetOrderDetail retrieves a single order by UUID.
func (c *Client) GetOrderDetail(req *private.GetOrderDetailRequest) (*private.Order, error) {
	return c.GetOrderDetailWithContext(context.Background(), req)
}

// GetOrderDetailWithContext retrieves a single order by UUID with context.
func (c *Client) GetOrderDetailWithContext(ctx context.Context, req *private.GetOrderDetailRequest) (*private.Order, error) {
	if err := req.Validate(); err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeAPI,
			Message: fmt.Sprintf("invalid request: %v", err),
			Err:     err,
		}
	}

	url := fmt.Sprintf("%s/v1/order?uuid=%s", c.base.BaseURL(), req.UUID)

	resp, err := c.doWithAuth(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeParse,
			Message: "read response failed",
			Err:     err,
		}
	}

	var order private.Order
	if err := json.Unmarshal(body, &order); err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeParse,
			Message: "parse response failed",
			Err:     err,
		}
	}

	return &order, nil
}

// GetOrders retrieves a list of orders.
func (c *Client) GetOrders(req *private.GetOrdersRequest) ([]private.Order, error) {
	return c.GetOrdersWithContext(context.Background(), req)
}

// GetOrdersWithContext retrieves a list of orders with context.
func (c *Client) GetOrdersWithContext(ctx context.Context, req *private.GetOrdersRequest) ([]private.Order, error) {
	if err := req.Validate(); err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeAPI,
			Message: fmt.Sprintf("invalid request: %v", err),
			Err:     err,
		}
	}

	url := c.base.BaseURL() + "/v1/orders"

	// Build query parameters
	query := ""
	if req.Market != "" {
		query = "market=" + req.Market
	}
	if len(req.UUIDs) > 0 {
		for _, uuid := range req.UUIDs {
			if query != "" {
				query += "&"
			}
			query += "uuids[]=" + uuid
		}
	}
	if req.State != "" {
		if query != "" {
			query += "&"
		}
		query += "state=" + req.State
	}
	if len(req.States) > 0 {
		for _, state := range req.States {
			if query != "" {
				query += "&"
			}
			query += "states[]=" + state
		}
	}
	if req.Page > 0 {
		if query != "" {
			query += "&"
		}
		query += fmt.Sprintf("page=%d", req.Page)
	}
	if req.Limit > 0 {
		if query != "" {
			query += "&"
		}
		query += fmt.Sprintf("limit=%d", req.Limit)
	}
	if req.OrderBy != "" {
		if query != "" {
			query += "&"
		}
		query += "order_by=" + req.OrderBy
	}

	if query != "" {
		url += "?" + query
	}

	resp, err := c.doWithAuth(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeParse,
			Message: "read response failed",
			Err:     err,
		}
	}

	var orders []private.Order
	if err := json.Unmarshal(body, &orders); err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeParse,
			Message: "parse response failed",
			Err:     err,
		}
	}

	return orders, nil
}

// GetOrderChance retrieves order chance information.
func (c *Client) GetOrderChance(req *private.GetOrderChanceRequest) (*private.OrderChance, error) {
	return c.GetOrderChanceWithContext(context.Background(), req)
}

// GetOrderChanceWithContext retrieves order chance information with context.
func (c *Client) GetOrderChanceWithContext(ctx context.Context, req *private.GetOrderChanceRequest) (*private.OrderChance, error) {
	if err := req.Validate(); err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeAPI,
			Message: fmt.Sprintf("invalid request: %v", err),
			Err:     err,
		}
	}

	url := fmt.Sprintf("%s/v1/orders/chance?market=%s", c.base.BaseURL(), req.Market)

	resp, err := c.doWithAuth(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeParse,
			Message: "read response failed",
			Err:     err,
		}
	}

	var chance private.OrderChance
	if err := json.Unmarshal(body, &chance); err != nil {
		return nil, &bithumbgo.Error{
			Type:    bithumbgo.ErrorTypeParse,
			Message: "parse response failed",
			Err:     err,
		}
	}

	return &chance, nil
}

// PlaceTWAPOrder places a TWAP algorithm order
func (c *Client) PlaceTWAPOrder(req *private.PlaceTWAPOrderRequest) (*private.TWAPOrder, error) {
	return c.PlaceTWAPOrderWithContext(context.Background(), req)
}

// PlaceTWAPOrderWithContext places a TWAP algorithm order with context
func (c *Client) PlaceTWAPOrderWithContext(ctx context.Context, req *private.PlaceTWAPOrderRequest) (*private.TWAPOrder, error) {
	if err := req.Validate(); err != nil {
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeAPI, Message: fmt.Sprintf("invalid request: %v", err), Err: err}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "marshal request failed", Err: err}
	}

	url := c.base.BaseURL() + "/v2/algo_orders/twap"

	resp, err := c.doWithAuth(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "read response failed", Err: err}
	}

	var order private.TWAPOrder
	if err := json.Unmarshal(respBody, &order); err != nil {
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "parse response failed", Err: err}
	}

	return &order, nil
}

// GetTWAPOrders retrieves TWAP order list
func (c *Client) GetTWAPOrders(req *private.GetTWAPOrdersRequest) ([]private.TWAPOrder, error) {
	return c.GetTWAPOrdersWithContext(context.Background(), req)
}

// GetTWAPOrdersWithContext retrieves TWAP order list with context
func (c *Client) GetTWAPOrdersWithContext(ctx context.Context, req *private.GetTWAPOrdersRequest) ([]private.TWAPOrder, error) {
	if err := req.Validate(); err != nil {
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeAPI, Message: fmt.Sprintf("invalid request: %v", err), Err: err}
	}

	url := c.base.BaseURL() + "/v1/algo_orders/twap"

	// Build query parameters
	params := query.New()
	if req.Market != "" {
		params.Add("market", req.Market)
	}
	if len(req.UUIDs) > 0 {
		params.AddStringSlice("uuids", req.UUIDs)
	}
	if req.State != "" {
		params.Add("state", string(req.State))
	}
	if req.NextKey != "" {
		params.Add("next_key", req.NextKey)
	}
	params.AddInt("limit", req.Limit)
	if req.OrderBy != "" {
		params.Add("order_by", req.OrderBy)
	}

	if queryStr := params.Encode(); queryStr != "" {
		url += "?" + queryStr
	}

	resp, err := c.doWithAuth(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "read response failed", Err: err}
	}

	var orders []private.TWAPOrder
	if err := json.Unmarshal(body, &orders); err != nil {
		return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "parse response failed", Err: err}
	}

	return orders, nil
}

// CancelTWAPOrder cancels a TWAP order
func (c *Client) CancelTWAPOrder(req *private.CancelTWAPOrderRequest) error {
	return c.CancelTWAPOrderWithContext(context.Background(), req)
}

// CancelTWAPOrderWithContext cancels a TWAP order with context
func (c *Client) CancelTWAPOrderWithContext(ctx context.Context, req *private.CancelTWAPOrderRequest) error {
	if err := req.Validate(); err != nil {
		return &bithumbgo.Error{Type: bithumbgo.ErrorTypeAPI, Message: fmt.Sprintf("invalid request: %v", err), Err: err}
	}

	url := fmt.Sprintf("%s/v2/algo_orders/twap/%s", c.base.BaseURL(), req.AlgoOrderID)

	reqBody, err := json.Marshal(map[string]string{})
	if err != nil {
		return &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "marshal request failed", Err: err}
	}

	resp, err := c.doWithAuth(ctx, http.MethodDelete, url, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Drain the response body to allow connection reuse
	io.Copy(io.Discard, resp.Body)
	return nil
}
