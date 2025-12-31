package public

import "errors"

// OrderBook represents order book information.
type OrderBook struct {
	// Market is the market identifier.
	Market string `json:"market"`
	// Timestamp is the Unix timestamp in milliseconds.
	Timestamp int64 `json:"timestamp"`
	// TotalAskSize is the total ask (sell) volume.
	TotalAskSize float64 `json:"total_ask_size"`
	// TotalBidSize is the total bid (buy) volume.
	TotalBidSize float64 `json:"total_bid_size"`
	// OrderBookUnits is the list of price levels.
	OrderBookUnits []OrderBookUnit `json:"orderbook_units"`
}

// OrderBookUnit represents a single price level in the order book.
type OrderBookUnit struct {
	// AskPrice is the ask (sell) price.
	AskPrice float64 `json:"ask_price"`
	// BidPrice is the bid (buy) price.
	BidPrice float64 `json:"bid_price"`
	// AskSize is the ask (sell) volume.
	AskSize float64 `json:"ask_size"`
	// BidSize is the bid (buy) volume.
	BidSize float64 `json:"bid_size"`
}

// GetOrderBookRequest is a request to get order book information.
type GetOrderBookRequest struct {
	Markets []string // Market identifiers
}

// Validate checks if the request is valid.
func (r *GetOrderBookRequest) Validate() error {
	if len(r.Markets) == 0 {
		return errors.New("markets is required")
	}
	return nil
}
