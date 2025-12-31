package public

import (
	"errors"
	"fmt"
)

// TradeSide represents the side of a trade.
type TradeSide string

const (
	TradeSideAsk TradeSide = "ASK"
	TradeSideBid TradeSide = "BID"
)

// Trade represents a single trade execution.
type Trade struct {
	// Market is the market identifier.
	Market string `json:"market"`
	// TradeDateUTC is the trade date in UTC (YYYYMMDD).
	TradeDateUTC string `json:"trade_date_utc"`
	// TradeTimeUTC is the trade time in UTC (HHmmss).
	TradeTimeUTC string `json:"trade_time_utc"`
	// Timestamp is the Unix timestamp in milliseconds.
	Timestamp int64 `json:"timestamp"`
	// TradePrice is the execution price.
	TradePrice float64 `json:"trade_price"`
	// TradeVolume is the execution volume.
	TradeVolume float64 `json:"trade_volume"`
	// PrevClosingPrice is the previous closing price.
	PrevClosingPrice float64 `json:"prev_closing_price"`
	// ChangePrice is the price change.
	ChangePrice float64 `json:"change_price"`
	// AskBid indicates if trade was ASK or BID.
	AskBid TradeSide `json:"ask_bid"`
	// SequentialID is the sequential trade ID.
	SequentialID int64 `json:"sequential_id"`
}

// GetRecentTradesRequest is a request to get recent trades.
type GetRecentTradesRequest struct {
	Market  string // Market identifier (e.g., "KRW-BTC")
	To      string // End time in HHmmss or HH:mm:ss format
	Count   int    // Number of trades (1-500, default 1)
	Cursor  string // Sequential ID for pagination
	DaysAgo int    // Days ago (1-7)
}

// Validate checks if the request is valid.
func (r *GetRecentTradesRequest) Validate() error {
	if r.Market == "" {
		return errors.New("market is required")
	}
	if r.Count < 0 || r.Count > 500 {
		return fmt.Errorf("count must be between 0 and 500, got %d", r.Count)
	}
	if r.DaysAgo < 0 || r.DaysAgo > 7 {
		return fmt.Errorf("daysAgo must be between 0 and 7, got %d", r.DaysAgo)
	}
	return nil
}
