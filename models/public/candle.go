package public

import (
	"errors"
	"fmt"
)

// Candle represents OHLCV candlestick data.
type Candle struct {
	// Market is the market identifier.
	Market string `json:"market"`
	// CandleDateTimeUTC is the candle time in UTC.
	CandleDateTimeUTC string `json:"candle_date_time_utc"`
	// CandleDateTimeKST is the candle time in KST.
	CandleDateTimeKST string `json:"candle_date_time_kst"`
	// OpeningPrice is the opening price.
	OpeningPrice float64 `json:"opening_price"`
	// HighPrice is the highest price.
	HighPrice float64 `json:"high_price"`
	// LowPrice is the lowest price.
	LowPrice float64 `json:"low_price"`
	// TradePrice is the closing price.
	TradePrice float64 `json:"trade_price"`
	// Timestamp is the Unix timestamp in milliseconds.
	Timestamp int64 `json:"timestamp"`
	// CandleAccTradePrice is the accumulated trade price.
	CandleAccTradePrice float64 `json:"candle_acc_trade_price"`
	// CandleAccTradeVolume is the accumulated trade volume.
	CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`
	// Unit is the candle interval unit.
	Unit int `json:"unit"`
}

// CandleInterval represents a candlestick time interval.
type CandleInterval string

const (
	CandleInterval1m   CandleInterval = "1"
	CandleInterval3m   CandleInterval = "3"
	CandleInterval5m   CandleInterval = "5"
	CandleInterval10m  CandleInterval = "10"
	CandleInterval15m  CandleInterval = "15"
	CandleInterval30m  CandleInterval = "30"
	CandleInterval60m  CandleInterval = "60"
	CandleInterval240m CandleInterval = "240"
)

// GetCandlestickRequest is a request to get candlestick data.
type GetCandlestickRequest struct {
	Market string         // Market identifier (e.g., "KRW-BTC")
	Unit   CandleInterval // Candle interval (1, 3, 5, 10, 15, 30, 60, 240)
	To     string         // End time (yyyy-MM-dd HH:mm:ss or yyyy-MM-ddTHH:mm:ss)
	Count  int            // Number of candles (max 200)
}

// Validate checks if the request is valid.
func (r *GetCandlestickRequest) Validate() error {
	if r.Market == "" {
		return errors.New("market is required")
	}
	if r.Count < 0 || r.Count > 200 {
		return fmt.Errorf("count must be between 0 and 200, got %d", r.Count)
	}
	return nil
}

// WeekCandle represents a week candle.
type WeekCandle struct {
	// Market is the market identifier.
	Market string `json:"market"`
	// CandleDateTimeKST is the candle base time (KST).
	CandleDateTimeKST string `json:"candle_date_time_kst"`
	// CandleDateTimeUTC is the candle base time (UTC).
	CandleDateTimeUTC string `json:"candle_date_time_utc"`
	// OpeningPrice is the opening price.
	OpeningPrice float64 `json:"opening_price"`
	// HighPrice is the high price.
	HighPrice float64 `json:"high_price"`
	// LowPrice is the low price.
	LowPrice float64 `json:"low_price"`
	// TradePrice is the closing price.
	TradePrice float64 `json:"trade_price"`
	// Timestamp is the candle end time (KST).
	Timestamp int64 `json:"timestamp"`
	// CandleAccTradePrice is the accumulated trade price.
	CandleAccTradePrice float64 `json:"candle_acc_trade_price"`
	// CandleAccTradeVolume is the accumulated trade volume.
	CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`
	// PrevClosingPrice is the previous day closing price.
	PrevClosingPrice float64 `json:"prev_closing_price"`
	// ChangePrice is the change from previous close.
	ChangePrice float64 `json:"change_price"`
	// ChangeRate is the change rate from previous close.
	ChangeRate float64 `json:"change_rate"`
	// ConvertedTradePrice is the converted trade price (optional).
	ConvertedTradePrice float64 `json:"converted_trade_price,omitempty"`
}

// GetWeekCandlesRequest is a request to get week candles.
type GetWeekCandlesRequest struct {
	// Market is the market code (e.g., "KRW-BTC").
	Market string
	// To is the last candle time (exclusive), ISO8061 format.
	To string
	// Count is the number of candles (max 200, default 1).
	Count int
	// ConvertingPriceUnit is the price unit for conversion (e.g., "KRW").
	ConvertingPriceUnit string
}

// Validate checks if the request is valid.
func (r *GetWeekCandlesRequest) Validate() error {
	if r.Market == "" {
		return fmt.Errorf("market is required")
	}
	if r.Count < 1 || r.Count > 200 {
		r.Count = 1
	}
	return nil
}

// MonthCandle represents a month candle.
type MonthCandle struct {
	// Market is the market identifier.
	Market string `json:"market"`
	// CandleDateTimeKST is the candle base time (KST).
	CandleDateTimeKST string `json:"candle_date_time_kst"`
	// CandleDateTimeUTC is the candle base time (UTC).
	CandleDateTimeUTC string `json:"candle_date_time_utc"`
	// OpeningPrice is the opening price.
	OpeningPrice float64 `json:"opening_price"`
	// HighPrice is the high price.
	HighPrice float64 `json:"high_price"`
	// LowPrice is the low price.
	LowPrice float64 `json:"low_price"`
	// TradePrice is the closing price.
	TradePrice float64 `json:"trade_price"`
	// Timestamp is the candle end time (KST).
	Timestamp int64 `json:"timestamp"`
	// CandleAccTradePrice is the accumulated trade price.
	CandleAccTradePrice float64 `json:"candle_acc_trade_price"`
	// CandleAccTradeVolume is the accumulated trade volume.
	CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`
	// PrevClosingPrice is the previous day closing price.
	PrevClosingPrice float64 `json:"prev_closing_price"`
	// ChangePrice is the change from previous close.
	ChangePrice float64 `json:"change_price"`
	// ChangeRate is the change rate from previous close.
	ChangeRate float64 `json:"change_rate"`
	// ConvertedTradePrice is the converted trade price (optional).
	ConvertedTradePrice float64 `json:"converted_trade_price,omitempty"`
}

// GetMonthCandlesRequest is a request to get month candles.
type GetMonthCandlesRequest struct {
	// Market is the market code (e.g., "KRW-BTC").
	Market string
	// To is the last candle time (exclusive), ISO8061 format.
	To string
	// Count is the number of candles (max 200, default 1).
	Count int
	// ConvertingPriceUnit is the price unit for conversion (e.g., "KRW").
	ConvertingPriceUnit string
}

// Validate checks if the request is valid.
func (r *GetMonthCandlesRequest) Validate() error {
	if r.Market == "" {
		return fmt.Errorf("market is required")
	}
	if r.Count < 1 || r.Count > 200 {
		r.Count = 1
	}
	return nil
}

// DayCandle represents a day candle.
type DayCandle struct {
	// Market is the market identifier.
	Market string `json:"market"`
	// CandleDateTimeKST is the candle base time (KST).
	CandleDateTimeKST string `json:"candle_date_time_kst"`
	// CandleDateTimeUTC is the candle base time (UTC).
	CandleDateTimeUTC string `json:"candle_date_time_utc"`
	// OpeningPrice is the opening price.
	OpeningPrice float64 `json:"opening_price"`
	// HighPrice is the high price.
	HighPrice float64 `json:"high_price"`
	// LowPrice is the low price.
	LowPrice float64 `json:"low_price"`
	// TradePrice is the closing price.
	TradePrice float64 `json:"trade_price"`
	// Timestamp is the candle end time (KST).
	Timestamp int64 `json:"timestamp"`
	// CandleAccTradePrice is the accumulated trade price.
	CandleAccTradePrice float64 `json:"candle_acc_trade_price"`
	// CandleAccTradeVolume is the accumulated trade volume.
	CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`
	// PrevClosingPrice is the previous day closing price.
	PrevClosingPrice float64 `json:"prev_closing_price"`
	// ChangePrice is the change from previous close.
	ChangePrice float64 `json:"change_price"`
	// ChangeRate is the change rate from previous close.
	ChangeRate float64 `json:"change_rate"`
	// ConvertedTradePrice is the converted trade price (optional).
	ConvertedTradePrice float64 `json:"converted_trade_price,omitempty"`
}

// GetDayCandlesRequest is a request to get day candles.
type GetDayCandlesRequest struct {
	// Market is the market code (e.g., "KRW-BTC").
	Market string
	// To is the last candle time (exclusive), ISO8061 format.
	To string
	// Count is the number of candles (max 200, default 1).
	Count int
	// ConvertingPriceUnit is the price unit for conversion (e.g., "KRW").
	ConvertingPriceUnit string
}

// Validate checks if the request is valid.
func (r *GetDayCandlesRequest) Validate() error {
	if r.Market == "" {
		return fmt.Errorf("market is required")
	}
	if r.Count < 1 || r.Count > 200 {
		r.Count = 1
	}
	return nil
}
