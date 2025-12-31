package private

import "fmt"

// Order represents an order.
type Order struct {
	// OrderID is the unique order identifier.
	OrderID string `json:"uuid"`
	// Market is the market identifier (e.g., "KRW-BTC").
	Market string `json:"market"`
	// Side is the order side ("bid" or "ask").
	Side string `json:"side"`
	// OrderType is the order type ("limit", "price", "market").
	OrderType string `json:"ord_type"`
	// CreatedAt is the order creation timestamp.
	CreatedAt string `json:"created_at"`
	// Price is the order price.
	Price string `json:"price"`
	// Volume is the order volume.
	Volume string `json:"volume"`
	// RemainingVolume is the unfilled volume.
	RemainingVolume string `json:"remaining_volume"`
	// ReservedFee is the reserved fee.
	ReservedFee string `json:"reserved_fee"`
	// RemainingFee is the remaining fee.
	RemainingFee string `json:"remaining_fee"`
	// PaidFee is the paid fee.
	PaidFee string `json:"paid_fee"`
	// State is the order state.
	State string `json:"state"`
	// Trades is the list of trades executed for this order.
	Trades []Trade `json:"trades"`
}

// Trade represents a single trade execution.
type Trade struct {
	// UUID is the trade identifier.
	UUID string `json:"uuid"`
	// Price is the trade price.
	Price string `json:"price"`
	// Volume is the trade volume.
	Volume string `json:"volume"`
	// Funds is the trade funds.
	Funds string `json:"funds"`
	// Side is the trade side.
	Side string `json:"side"`
}

// PlaceOrderRequest is a request to place an order.
type PlaceOrderRequest struct {
	// Market is the market identifier (e.g., "KRW-BTC").
	Market string
	// Side is the order side ("bid" for buy, "ask" for sell).
	Side string
	// OrderType is the order type ("limit", "price", "market").
	OrderType string
	// Price is the order price (required for limit and price orders).
	Price string
	// Volume is the order volume (required for limit and market orders).
	Volume string
}

// Validate checks if the request is valid.
func (r *PlaceOrderRequest) Validate() error {
	if r.Market == "" {
		return fmt.Errorf("market is required")
	}
	if r.Side != "bid" && r.Side != "ask" {
		return fmt.Errorf("side must be 'bid' or 'ask'")
	}
	return nil
}

// CancelOrderRequest is a request to cancel an order.
type CancelOrderRequest struct {
	// UUID is the order identifier to cancel.
	UUID string
}

// Validate checks if the request is valid.
func (r *CancelOrderRequest) Validate() error {
	if r.UUID == "" {
		return fmt.Errorf("UUID is required")
	}
	return nil
}
