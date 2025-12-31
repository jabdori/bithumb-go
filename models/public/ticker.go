// Package public provides data models for Bithumb Public API responses.
package public

import "fmt"

// Ticker represents current market price information for a trading pair.
type Ticker struct {
	// Market is the market identifier (e.g., "KRW-BTC").
	Market string `json:"market"`
	// TradeDate is the trade date in YYYYMMDD format.
	TradeDate string `json:"trade_date"`
	// TradeTime is the trade time in HHmmss format.
	TradeTime string `json:"trade_time"`
	// TradeDateKST is the trade date in KST (YYYYMMDD).
	TradeDateKST string `json:"trade_date_kst"`
	// TradeTimeKST is the trade time in KST (HHmmss).
	TradeTimeKST string `json:"trade_time_kst"`
	// TradeTimestamp is the Unix timestamp in milliseconds.
	TradeTimestamp int64 `json:"trade_timestamp"`
	// OpeningPrice is the first trade price in the interval.
	OpeningPrice float64 `json:"opening_price"`
	// HighPrice is the highest trade price in the interval.
	HighPrice float64 `json:"high_price"`
	// LowPrice is the lowest trade price in the interval.
	LowPrice float64 `json:"low_price"`
	// TradePrice is the latest trade price.
	TradePrice float64 `json:"trade_price"`
	// PrevClosingPrice is the previous day's closing price.
	PrevClosingPrice float64 `json:"prev_closing_price"`
	// Change indicates price change direction ("RISE", "FALL", "EVEN").
	Change string `json:"change"`
	// ChangePrice is the absolute price change from previous close.
	ChangePrice float64 `json:"change_price"`
	// ChangeRate is the price change as a percentage.
	ChangeRate float64 `json:"change_rate"`
	// SignedChangePrice is the signed price change.
	SignedChangePrice float64 `json:"signed_change_price"`
	// SignedChangeRate is the signed change rate.
	SignedChangeRate float64 `json:"signed_change_rate"`
	// TradeVolume is the 24h trading volume in base currency.
	TradeVolume float64 `json:"trade_volume"`
	// AccTradePrice is the accumulated trade price in quote currency.
	AccTradePrice float64 `json:"acc_trade_price"`
	// AccTradePrice24h is the 24h accumulated trade price.
	AccTradePrice24h float64 `json:"acc_trade_price_24h"`
	// AccTradeVolume is the accumulated trading volume.
	AccTradeVolume float64 `json:"acc_trade_volume"`
	// AccTradeVolume24h is the 24h accumulated trading volume.
	AccTradeVolume24h float64 `json:"acc_trade_volume_24h"`
	// Highest52WeekPrice is the 52-week highest price.
	Highest52WeekPrice float64 `json:"highest_52_week_price"`
	// Highest52WeekDate is the date of 52-week high (YYYYMMDD).
	Highest52WeekDate string `json:"highest_52_week_date"`
	// Lowest52WeekPrice is the 52-week lowest price.
	Lowest52WeekPrice float64 `json:"lowest_52_week_price"`
	// Lowest52WeekDate is the date of 52-week low (YYYYMMDD).
	Lowest52WeekDate string `json:"lowest_52_week_date"`
	// Timestamp is the response timestamp in milliseconds.
	Timestamp int64 `json:"timestamp"`
}

// GetTickerRequest is a request to get ticker information.
type GetTickerRequest struct {
	Markets []string // Market identifiers (e.g., "KRW-BTC")
}

// Validate checks if the request is valid.
func (r *GetTickerRequest) Validate() error {
	if len(r.Markets) == 0 {
		return fmt.Errorf("markets is required")
	}
	return nil
}
