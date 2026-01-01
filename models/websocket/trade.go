// Package websocket provides data models for Bithumb WebSocket API messages.
package websocket

// TradeMessage represents a trade WebSocket message.
type TradeMessage struct {
	// Type is the message type ("trade").
	Type string `json:"type"`
	// Code is the market code (e.g., "KRW-BTC").
	Code string `json:"code"`
	// TradePrice is the trade price.
	TradePrice float64 `json:"trade_price"`
	// TradeVolume is the trade volume.
	TradeVolume float64 `json:"trade_volume"`
	// AskBid is the ask/bid type ("ASK": sell, "BID": buy).
	AskBid string `json:"ask_bid"`
	// PrevClosingPrice is the previous closing price.
	PrevClosingPrice float64 `json:"prev_closing_price"`
	// Change is the change type ("RISE", "EVEN", "FALL").
	Change string `json:"change"`
	// ChangePrice is the price change (unsigned).
	ChangePrice float64 `json:"change_price"`
	// TradeDate is the trade date (KST, yyyy-MM-dd).
	TradeDate string `json:"trade_date"`
	// TradeTime is the trade time (KST, HH:mm:ss).
	TradeTime string `json:"trade_time"`
	// TradeTimestamp is the trade timestamp (milliseconds).
	TradeTimestamp int64 `json:"trade_timestamp"`
	// SequentialID is the unique sequential ID.
	SequentialID int64 `json:"sequential_id"`
	// StreamType is the stream type ("SNAPSHOT" or "REALTIME").
	StreamType string `json:"stream_type"`
	// Timestamp is the message timestamp (milliseconds).
	Timestamp int64 `json:"timestamp"`
}
