package public

// Candle represents OHLCV candlestick data.
type Candle struct {
	Market               string  `json:"market"`
	CandleDateTimeUTC    string  `json:"candle_date_time_utc"`
	CandleDateTimeKST    string  `json:"candle_date_time_kst"`
	OpeningPrice         float64 `json:"opening_price"`
	HighPrice            float64 `json:"high_price"`
	LowPrice             float64 `json:"low_price"`
	TradePrice           float64 `json:"trade_price"`
	Timestamp            int64   `json:"timestamp"`
	CandleAccTradePrice  float64 `json:"candle_acc_trade_price"`
	CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`
	Unit                 int     `json:"unit"`
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
	Market string // KRW-BTC, BTC-ETH etc.
	Unit   CandleInterval // 1, 3, 5, 10, 15, 30, 60, 240
	To     string         // yyyy-MM-dd HH:mm:ss or yyyy-MM-ddTHH:mm:ss
	Count  int            // max 200
}
