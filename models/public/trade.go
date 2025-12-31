package public

// Trade represents a single trade execution.
type Trade struct {
	Market          string  `json:"market"`
	TradeDateUTC    string  `json:"trade_date_utc"`
	TradeTimeUTC    string  `json:"trade_time_utc"`
	Timestamp       int64   `json:"timestamp"`
	TradePrice      float64 `json:"trade_price"`
	TradeVolume     float64 `json:"trade_volume"`
	PrevClosingPrice float64 `json:"prev_closing_price"`
	ChangePrice     float64 `json:"change_price"`
	AskBid         string  `json:"ask_bid"`
	SequentialID   int64   `json:"sequential_id"`
}

// GetRecentTradesRequest is a request to get recent trades.
type GetRecentTradesRequest struct {
	Market  string // KRW-BTC, BTC-ETH etc.
	To      string // HHmmss or HH:mm:ss
	Count   int    // 1~500, default 1
	Cursor  string // sequentialId
	DaysAgo int    // 1~7
}
