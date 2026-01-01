package public_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hysuki/bithumb-go/client"
	publicmodels "github.com/hysuki/bithumb-go/models/public"
	"github.com/hysuki/bithumb-go/public"
)

// Test helpers
func assertEqual[T comparable](t *testing.T, expected, actual T, msgAndArgs ...interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %v, got %v. %s", expected, actual, msgAndArgs)
	}
}

func assertNil(t *testing.T, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if actual != nil {
		t.Errorf("Expected nil, got %v. %s", actual, msgAndArgs)
	}
}

func TestGetMarketAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/v1/market/all", r.URL.Path)
		assertEqual(t, "true", r.URL.Query().Get("isDetails"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"market":"KRW-BTC","korean_name":"비트코인","english_name":"Bitcoin","market_warning":"NONE"}]`))
	}))
	defer server.Close()

	baseClient, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithHTTPClient(server.Client()))
	c := public.NewClient(baseClient)

	markets, err := c.GetMarketAll(true)
	assertNil(t, err)
	assertEqual(t, 1, len(markets))
	assertEqual(t, "KRW-BTC", markets[0].Market)
	assertEqual(t, "비트코인", markets[0].KoreanName)
	assertEqual(t, "Bitcoin", markets[0].EnglishName)
	assertEqual(t, "NONE", markets[0].MarketWarning)
}

func TestGetDayCandles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/v1/candles/days", r.URL.Path)
		assertEqual(t, "KRW-BTC", r.URL.Query().Get("market"))
		assertEqual(t, "10", r.URL.Query().Get("count"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"market":"KRW-BTC","candle_date_time_kst":"2024-01-01T00:00:00","opening_price":50000000,"high_price":51000000,"low_price":49000000,"trade_price":50500000,"timestamp":1704067200000,"candle_acc_trade_price":1000000000,"candle_acc_trade_volume":20.0,"prev_closing_price":50000000,"change_price":500000,"change_rate":0.01}]`))
	}))
	defer server.Close()

	baseClient, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithHTTPClient(server.Client()))
	c := public.NewClient(baseClient)

	req := &publicmodels.GetDayCandlesRequest{
		Market: "KRW-BTC",
		Count:  10,
	}
	candles, err := c.GetDayCandles(req)
	assertNil(t, err)
	assertEqual(t, 1, len(candles))
	assertEqual(t, 50500000.0, candles[0].TradePrice)
}

func TestGetWeekCandles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/v1/candles/weeks", r.URL.Path)
		assertEqual(t, "KRW-BTC", r.URL.Query().Get("market"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"market":"KRW-BTC","candle_date_time_kst":"2024-01-01T00:00:00","opening_price":50000000,"high_price":51000000,"low_price":49000000,"trade_price":50500000,"timestamp":1704067200000,"candle_acc_trade_price":1000000000,"candle_acc_trade_volume":20.0,"prev_closing_price":50000000,"change_price":500000,"change_rate":0.01}]`))
	}))
	defer server.Close()

	baseClient, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithHTTPClient(server.Client()))
	c := public.NewClient(baseClient)

	req := &publicmodels.GetWeekCandlesRequest{
		Market: "KRW-BTC",
		Count:  10,
	}
	candles, err := c.GetWeekCandles(req)
	assertNil(t, err)
	assertEqual(t, 1, len(candles))
	assertEqual(t, 50500000.0, candles[0].TradePrice)
}

func TestGetMonthCandles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/v1/candles/months", r.URL.Path)
		assertEqual(t, "KRW-BTC", r.URL.Query().Get("market"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"market":"KRW-BTC","candle_date_time_kst":"2024-01-01T00:00:00","opening_price":50000000,"high_price":51000000,"low_price":49000000,"trade_price":50500000,"timestamp":1704067200000,"candle_acc_trade_price":1000000000,"candle_acc_trade_volume":20.0,"prev_closing_price":50000000,"change_price":500000,"change_rate":0.01}]`))
	}))
	defer server.Close()

	baseClient, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithHTTPClient(server.Client()))
	c := public.NewClient(baseClient)

	req := &publicmodels.GetMonthCandlesRequest{
		Market: "KRW-BTC",
		Count:  10,
	}
	candles, err := c.GetMonthCandles(req)
	assertNil(t, err)
	assertEqual(t, 1, len(candles))
	assertEqual(t, 50500000.0, candles[0].TradePrice)
}

func TestGetTicker(t *testing.T) {
	baseClient, _ := client.NewClient()
	c := public.NewClient(baseClient)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ticker, err := c.GetTickerWithContext(ctx, &publicmodels.GetTickerRequest{
		Markets: []string{"KRW-BTC"},
	})

	if err != nil {
		t.Fatalf("GetTickerWithContext() error = %v", err)
	}

	if len(ticker) == 0 {
		t.Fatal("GetTickerWithContext() returned empty slice")
	}

	if ticker[0].Market != "KRW-BTC" {
		t.Errorf("Market = %v, want KRW-BTC", ticker[0].Market)
	}
}

func TestGetOrderBook(t *testing.T) {
	baseClient, _ := client.NewClient()
	c := public.NewClient(baseClient)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderbook, err := c.GetOrderBookWithContext(ctx, &publicmodels.GetOrderBookRequest{
		Markets: []string{"KRW-BTC"},
	})

	if err != nil {
		t.Fatalf("GetOrderBookWithContext() error = %v", err)
	}

	if orderbook == nil {
		t.Fatal("GetOrderBookWithContext() returned nil")
	}

	if orderbook.Market != "KRW-BTC" {
		t.Errorf("Market = %v, want KRW-BTC", orderbook.Market)
	}
}

func TestGetRecentTrades(t *testing.T) {
	baseClient, _ := client.NewClient()
	c := public.NewClient(baseClient)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trades, err := c.GetRecentTradesWithContext(ctx, &publicmodels.GetRecentTradesRequest{
		Market: "KRW-BTC",
		Count:  10,
	})

	if err != nil {
		t.Fatalf("GetRecentTradesWithContext() error = %v", err)
	}

	if len(trades) == 0 {
		t.Fatal("GetRecentTradesWithContext() returned empty slice")
	}

	if trades[0].Market != "KRW-BTC" {
		t.Errorf("Market = %v, want KRW-BTC", trades[0].Market)
	}
}

func TestGetCandlestick(t *testing.T) {
	baseClient, _ := client.NewClient()
	c := public.NewClient(baseClient)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	candles, err := c.GetCandlestickWithContext(ctx, &publicmodels.GetCandlestickRequest{
		Market: "KRW-BTC",
		Unit:   publicmodels.CandleInterval1m,
		Count:  10,
	})

	if err != nil {
		t.Fatalf("GetCandlestickWithContext() error = %v", err)
	}

	if len(candles) == 0 {
		t.Fatal("GetCandlestickWithContext() returned empty slice")
	}

	if candles[0].Market != "KRW-BTC" {
		t.Errorf("Market = %v, want KRW-BTC", candles[0].Market)
	}
}

func TestGetWarnings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/v1/market/warning", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"market":"KRW-BTC","warning_type":"PRICE_SUDDEN_FLUCTUATION","end_date":"2026-01-31 23:59:59"}]`))
	}))
	defer server.Close()

	baseClient, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithHTTPClient(server.Client()))
	c := public.NewClient(baseClient)

	warnings, bithumbErr := c.GetWarnings()
	if bithumbErr != nil {
		t.Fatalf("Expected nil error, got %v", bithumbErr)
	}
	assertEqual(t, 1, len(warnings))
	assertEqual(t, "KRW-BTC", warnings[0].Market)
	assertEqual(t, publicmodels.WarningPriceSuddenFluctuation, warnings[0].WarningType)
	assertEqual(t, "2026-01-31 23:59:59", warnings[0].EndDate)
}
