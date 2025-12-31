// Package main shows how to use Bithumb WebSocket
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/hysuki/bithumb-go/client"
	"github.com/hysuki/bithumb-go/models/public"
	"github.com/hysuki/bithumb-go/websocket"
)

func main() {
	// Create client
	c, err := client.NewClient()
	if err != nil {
		log.Fatalf("클라이언트 생성 실패: %v", err)
	}

	ws := c.Websocket

	// Context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Connect to WebSocket
	fmt.Println("WebSocket 연결 중...")
	if err := ws.Connect(ctx); err != nil {
		log.Fatalf("WebSocket 연결 실패: %v", err)
	}
	defer ws.Close()

	fmt.Println("WebSocket 연결됨!")

	// Setup message handlers with thread-safe counting
	var mu sync.Mutex
	tickerCount := 0
	orderbookCount := 0

	handlers := websocket.MessageHandlers{
		Ticker: func(msg []byte) error {
			mu.Lock()
			tickerCount++
			count := tickerCount
			mu.Unlock()

			var ticker public.Ticker
			if err := json.Unmarshal(msg, &ticker); err != nil {
				return err
			}

			// Print first 10 tickers in detail
			if count <= 10 {
				fmt.Printf("[Ticker #%d] %s: %.0f KRW (%.2f%%)\n",
					count, ticker.Market, ticker.TradePrice, ticker.ChangeRate*100)
			} else if count == 11 {
				fmt.Println("[Ticker] ... (더 이상 표시 안 함)")
			}

			return nil
		},

		OrderBook: func(msg []byte) error {
			mu.Lock()
			orderbookCount++
			count := orderbookCount
			mu.Unlock()

			var orderbook public.OrderBook
			if err := json.Unmarshal(msg, &orderbook); err != nil {
				return err
			}

			// Print first 5 orderbooks in detail
			if count <= 5 {
				fmt.Printf("[OrderBook #%d] %s: %d 호가 단위\n",
					count, orderbook.Market, len(orderbook.OrderBookUnits))
			} else if count == 6 {
				fmt.Println("[OrderBook] ... (더 이상 표시 안 함)")
			}

			return nil
		},

		Trade: func(msg []byte) error {
			// Just acknowledge trade messages
			return nil
		},
	}

	// Subscribe to ticker and orderbook for BTC
	fmt.Println("\n구독 시작: KRW-BTC Ticker & OrderBook")
	params := []*websocket.SubscriptionParam{
		{
			Type:    websocket.SubscriptionTypeTicker,
			Symbols: []string{"KRW-BTC"},
		},
		{
			Type:    websocket.SubscriptionTypeOrderBook,
			Symbols: []string{"KRW-BTC"},
		},
	}

	if err := ws.Subscribe(params, handlers); err != nil {
		log.Fatalf("구독 실패: %v", err)
	}

	// Optional: Configure reconnection settings
	ws.SetReconnectDelay(5 * time.Second)
	ws.SetReconnectTimeout(30 * time.Second)
	fmt.Println("재연결 설정: 5초 지연, 30초 타임아웃")

	fmt.Println("\n실시간 데이터 수집 중... (Ctrl+C to quit)")

	// Wait for interrupt signal or 30 second timeout
	timeout := time.After(30 * time.Second)
	select {
	case <-sigCh:
		fmt.Println("\n\n사용자 중단 신호 수신")
	case <-timeout:
		fmt.Println("\n\n30초 타임아웃 도달")
	case <-ctx.Done():
		fmt.Println("\n\n컨텍스트 취소됨")
	}

	// Final statistics
	mu.Lock()
	fmt.Printf("\n수집된 메시지: Ticker %d개, OrderBook %d개\n",
		tickerCount, orderbookCount)
	mu.Unlock()

	fmt.Println("프로그램 종료")
}
