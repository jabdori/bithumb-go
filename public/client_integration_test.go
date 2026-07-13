//go:build integration

package public_test

import (
	"context"
	"testing"
	"time"

	"github.com/hysuki/bithumb-go/client"
	publicmodels "github.com/hysuki/bithumb-go/models/public"
	"github.com/hysuki/bithumb-go/public"
)

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
