package public

// OrderBook represents order book information.
type OrderBook struct {
	Market         string          `json:"market"`
	Timestamp      int64           `json:"timestamp"`
	TotalAskSize   float64         `json:"total_ask_size"`
	TotalBidSize   float64         `json:"total_bid_size"`
	OrderBookUnits []OrderBookUnit `json:"orderbook_units"`
}

// OrderBookUnit represents a single price level in the order book.
type OrderBookUnit struct {
	AskPrice float64 `json:"ask_price"`
	BidPrice float64 `json:"bid_price"`
	AskSize  float64 `json:"ask_size"`
	BidSize  float64 `json:"bid_size"`
}

// GetOrderBookRequest is a request to get order book information.
type GetOrderBookRequest struct {
	Markets []string
}
