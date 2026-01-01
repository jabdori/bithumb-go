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
	// Locked is the amount locked in the order.
	Locked string `json:"locked"`
	// ExecutedVolume is the executed volume.
	ExecutedVolume string `json:"executed_volume"`
	// TradesCount is the number of trades.
	TradesCount int `json:"trades_count"`
	// Trades is the list of trades executed for this order.
	Trades []Trade `json:"trades"`
}

// Trade represents a single trade execution.
type Trade struct {
	// UUID is the trade identifier.
	UUID string `json:"uuid"`
	// Market is the market identifier.
	Market string `json:"market"`
	// Price is the trade price.
	Price string `json:"price"`
	// Volume is the trade volume.
	Volume string `json:"volume"`
	// Funds is the trade funds.
	Funds string `json:"funds"`
	// Side is the trade side.
	Side string `json:"side"`
	// CreatedAt is the trade creation timestamp.
	CreatedAt string `json:"created_at"`
}

// Order side constants.
const (
	OrderSideBid = "bid"
	OrderSideAsk = "ask"
)

// Order type constants.
const (
	OrderTypeLimit  = "limit"
	OrderTypePrice  = "price"
	OrderTypeMarket = "market"
)

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

	if r.Side != OrderSideBid && r.Side != OrderSideAsk {
		return fmt.Errorf("side must be '%s' or '%s'", OrderSideBid, OrderSideAsk)
	}

	// Validate order type
	if r.OrderType != OrderTypeLimit && r.OrderType != OrderTypePrice && r.OrderType != OrderTypeMarket {
		return fmt.Errorf("order_type must be '%s', '%s', or '%s'", OrderTypeLimit, OrderTypePrice, OrderTypeMarket)
	}

	// Validate price/volume based on order type
	switch r.OrderType {
	case OrderTypeLimit, OrderTypePrice:
		if r.Price == "" {
			return fmt.Errorf("price is required for %s orders", r.OrderType)
		}
		if r.Volume == "" {
			return fmt.Errorf("volume is required for %s orders", r.OrderType)
		}
	case OrderTypeMarket:
		if r.Volume == "" {
			return fmt.Errorf("volume is required for market orders")
		}
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

// GetOrderDetailRequest is a request to get order detail.
type GetOrderDetailRequest struct {
	// UUID is the order identifier.
	UUID string
}

// Validate checks if the request is valid.
func (r *GetOrderDetailRequest) Validate() error {
	if r.UUID == "" {
		return fmt.Errorf("UUID is required")
	}
	return nil
}

// GetOrdersRequest is a request to get order list.
type GetOrdersRequest struct {
	// Market is the market identifier (e.g., "KRW-BTC").
	Market string
	// UUIDs is the list of order UUIDs.
	UUIDs []string
	// State is the order state filter.
	State string
	// States is the list of order states.
	States []string
	// Page is the page number (default: 1).
	Page int
	// Limit is the limit of orders per page (default: 100, max: 100).
	Limit int
	// OrderBy is the sort order ("asc" or "desc", default: "desc").
	OrderBy string
}

// Validate checks if the request is valid.
func (r *GetOrdersRequest) Validate() error {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.Limit < 1 || r.Limit > 100 {
		r.Limit = 100
	}
	if r.OrderBy != "asc" && r.OrderBy != "desc" {
		r.OrderBy = "desc"
	}
	return nil
}
