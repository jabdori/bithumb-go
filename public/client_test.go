package public_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
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
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
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
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
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
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
	assertEqual(t, 1, len(candles))
	assertEqual(t, 50500000.0, candles[0].TradePrice)
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

func TestGetNotices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/v1/notice", r.URL.Path)
		assertEqual(t, "5", r.URL.Query().Get("count"))

		notices := []publicmodels.Notice{
			{
				Categories:  []string{"공지"},
				Title:       "테스트 공지",
				PCURL:       "https://example.com/notice/1",
				PublishedAt: "2026-01-01 10:00:00",
				ModifiedAt:  "2026-01-01 10:00:00",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(notices)
	}))
	defer server.Close()

	baseClient, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithHTTPClient(server.Client()))
	c := public.NewClient(baseClient)

	notices, bithumbErr := c.GetNotices(&publicmodels.GetNoticesRequest{Count: 5})
	if bithumbErr != nil {
		t.Fatalf("GetNotices failed: %v", bithumbErr)
	}

	if len(notices) != 1 {
		t.Fatalf("Expected 1 notice, got %d", len(notices))
	}

	if notices[0].Title != "테스트 공지" {
		t.Errorf("Expected title '테스트 공지', got %s", notices[0].Title)
	}
}

func TestGetChainFees(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/fee/inout/ALL" {
			t.Errorf("Expected path /v2/fee/inout/ALL, got %s", r.URL.Path)
		}

		fees := []publicmodels.ChainFee{
			{
				Name:     "비트코인",
				Currency: "BTC",
				Networks: []publicmodels.NetworkFee{
					{
						NetName:                 "BTC",
						DepositFeeQuantity:      "0",
						DepositMinimumQuantity:  "0.0004",
						WithdrawFeeQuantity:     "0.0005",
						WithdrawMinimumQuantity: "0.001",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fees)
	}))
	defer server.Close()

	baseClient, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithHTTPClient(server.Client()))
	c := public.NewClient(baseClient)

	fees, bithumbErr := c.GetChainFees("ALL")
	if bithumbErr != nil {
		t.Fatalf("GetChainFees failed: %v", bithumbErr)
	}

	if len(fees) != 1 {
		t.Fatalf("Expected 1 fee entry, got %d", len(fees))
	}

	if fees[0].Currency != "BTC" {
		t.Errorf("Expected currency BTC, got %s", fees[0].Currency)
	}

	if len(fees[0].Networks) != 1 {
		t.Fatalf("Expected 1 network, got %d", len(fees[0].Networks))
	}

	if fees[0].Networks[0].NetName != "BTC" {
		t.Errorf("Expected net_name BTC, got %s", fees[0].Networks[0].NetName)
	}
}

func TestGetChainFees_EmptyCurrency(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	baseClient, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithHTTPClient(server.Client()))
	c := public.NewClient(baseClient)

	_, bithumbErr := c.GetChainFees("")
	if bithumbErr == nil {
		t.Fatal("Expected error for empty currency, got nil")
	}

	if bithumbErr.Message != "currency is required" {
		t.Errorf("Expected error message 'currency is required', got '%s'", bithumbErr.Message)
	}
}
