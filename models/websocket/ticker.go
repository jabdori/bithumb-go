// Package websocket provides data models for Bithumb WebSocket API messages.
package websocket

// TickerMessage represents a ticker WebSocket message.
type TickerMessage struct {
	// Type is the message type ("ticker").
	Type string `json:"type"`
	// Code is the market code (e.g., "KRW-BTC").
	Code string `json:"code"`
	// OpeningPrice is the opening price.
	OpeningPrice float64 `json:"opening_price"`
	// HighPrice is the highest price.
	HighPrice float64 `json:"high_price"`
	// LowPrice is the lowest price.
	LowPrice float64 `json:"low_price"`
	// TradePrice is the current trade price.
	TradePrice float64 `json:"trade_price"`
	// PrevClosingPrice is the previous closing price.
	PrevClosingPrice float64 `json:"prev_closing_price"`
	// Change is the change type ("RISE", "EVEN", "FALL").
	Change string `json:"change"`
	// ChangePrice is the price change (unsigned).
	ChangePrice float64 `json:"change_price"`
	// ChangeRate is the change rate (percentage).
	ChangeRate float64 `json:"change_rate"`
	// SignedChangeRate is the signed change rate (percentage).
	SignedChangeRate float64 `json:"signed_change_rate"`
	// TradeVolume is the trade volume.
	TradeVolume float64 `json:"trade_volume"`
	// AccTradeVolume is the accumulated trade volume.
	AccTradeVolume float64 `json:"acc_trade_volume"`
	// AccTradePrice is the accumulated trade price.
	AccTradePrice float64 `json:"acc_trade_price"`
	// AccTradeVolume24h is the 24-hour accumulated trade volume.
	AccTradeVolume24h float64 `json:"acc_trade_volume_24h"`
	// Highest52WeekPrice is the 52-week highest price.
	Highest52WeekPrice float64 `json:"highest_52_week_price"`
	// Lowest52WeekPrice is the 52-week lowest price.
	Lowest52WeekPrice float64 `json:"lowest_52_week_price"`
	// Highest52WeekDate is the date of 52-week highest price.
	Highest52WeekDate string `json:"highest_52_week_date"`
	// Lowest52WeekDate is the date of 52-week lowest price.
	Lowest52WeekDate string `json:"lowest_52_week_date"`
	// StreamType is the stream type ("SNAPSHOT" or "REALTIME").
	StreamType string `json:"stream_type"`
	// Timestamp is the message timestamp (milliseconds).
	Timestamp int64 `json:"timestamp"`
}
