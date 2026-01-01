// Package main shows how to use zap logger with Bithumb WebSocket
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hysuki/bithumb-go/client"
	"github.com/hysuki/bithumb-go/logger"
	wsmodels "github.com/hysuki/bithumb-go/models/websocket"
	"github.com/hysuki/bithumb-go/websocket"
	"go.uber.org/zap"
)

// ZapLogger implements logger.Logger interface using zap
type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(l *zap.Logger) *ZapLogger {
	return &ZapLogger{logger: l}
}

func (z *ZapLogger) Debug(msg string, fields ...logger.Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	z.logger.Debug(msg, zapFields...)
}

func (z *ZapLogger) Info(msg string, fields ...logger.Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	z.logger.Info(msg, zapFields...)
}

func (z *ZapLogger) Error(msg string, fields ...logger.Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	z.logger.Error(msg, zapFields...)
}

func main() {
	// Create zap logger
	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()

	// Create Bithumb client
	c, _ := client.NewClient()
	ws := c.Websocket

	// Set custom logger
	ws.SetLogger(NewZapLogger(zapLogger))

	// Connect to WebSocket
	ctx := context.Background()
	if err := ws.Connect(ctx); err != nil {
		panic(err)
	}
	defer ws.Close()

	// Setup message handlers
	handlers := websocket.MessageHandlers{
		Ticker: websocket.HandlerFunc(func(msg []byte) error {
			var ticker wsmodels.TickerMessage
			if err := json.Unmarshal(msg, &ticker); err != nil {
				return err
			}
			fmt.Printf("[Ticker] %s: %.0f KRW\n", ticker.Code, ticker.TradePrice)
			return nil
		}),

		OrderBook: websocket.HandlerFunc(func(msg []byte) error {
			var orderbook wsmodels.OrderBookMessage
			if err := json.Unmarshal(msg, &orderbook); err != nil {
				return err
			}
			fmt.Printf("[OrderBook] %s: %d levels\n", orderbook.Code, len(orderbook.OrderBookUnits))
			return nil
		}),
	}

	// Subscribe to ticker and orderbook for BTC
	params := []*websocket.SubscriptionParam{
		{Type: websocket.SubscriptionTypeTicker, Codes: []string{"KRW-BTC"}},
		{Type: websocket.SubscriptionTypeOrderBook, Codes: []string{"KRW-BTC"}},
	}

	if err := ws.Subscribe(params, handlers); err != nil {
		panic(err)
	}

	fmt.Println("WebSocket connected with zap logger! Press Ctrl+C to quit")

	// Wait for interrupt signal
	<-ctx.Done()
}
