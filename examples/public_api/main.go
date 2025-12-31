// Package main shows basic usage of the Bithumb Go SDK.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hysuki/bithumb-go/client"
	"github.com/hysuki/bithumb-go/models/public"
)

func main() {
	// Create client for Public API
	c, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	// Get ticker information
	ticker, err := c.Public.GetTicker(&public.GetTickerRequest{
		Markets: []string{"KRW-BTC"},
	})
	if err != nil {
		log.Printf("Failed to get ticker: %v", err)
		return
	}
	fmt.Printf("BTC Price: %.2f\n", ticker[0].TradePrice)

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orderbook, err := c.Public.GetOrderBookWithContext(ctx, &public.GetOrderBookRequest{
		Markets: []string{"KRW-BTC"},
	})
	if err != nil {
		log.Printf("Failed to get orderbook: %v", err)
		return
	}
	fmt.Printf("OrderBook: %d levels\n", len(orderbook.OrderBookUnits))

	// Private API usage (requires API key)
	// cWithAuth, _ := client.NewClient(
	//     client.WithAPIKey("your-key", "your-secret"),
	// )
	// accounts, _ := cWithAuth.Private.GetAccount(&private.GetAccountRequest{})
}
