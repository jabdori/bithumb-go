// Package main shows how to use the Bithumb Private API
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hysuki/bithumb-go/client"
	"github.com/hysuki/bithumb-go/models/private"
)

func main() {
	// Get API credentials from environment variables
	accessKey := os.Getenv("BITHUMB_ACCESS_KEY")
	secretKey := os.Getenv("BITHUMB_SECRET_KEY")

	if accessKey == "" || secretKey == "" {
		log.Fatal("BITHUMB_ACCESS_KEY and BITHUMB_SECRET_KEY environment variables are required")
	}

	// Create client with API credentials
	c, err := client.NewClient(
		client.WithAPIKey(accessKey, secretKey),
	)
	if err != nil {
		log.Fatalf("클라이언트 생성 실패: %v", err)
	}

	// Check if Private client is available
	if c.Private == nil {
		log.Fatal("Private API client is not available")
	}

	// Get account information
	fmt.Println("=== 계정 정보 조회 ===")
	accounts, err := c.Private.GetAccount(&private.GetAccountRequest{})
	if err != nil {
		log.Printf("계정 조회 실패: %v", err)
	} else {
		for _, acc := range accounts {
			fmt.Printf("%s: 보유 %s (잠김 %s)\n", acc.Currency, acc.Balance, acc.Locked)
		}
	}

	// Get BTC account specifically
	fmt.Println("\n=== BTC 계정 조회 ===")
	btcAccounts, err := c.Private.GetAccount(&private.GetAccountRequest{
		Currency: "BTC",
	})
	if err != nil {
		log.Printf("BTC 계정 조회 실패: %v", err)
	} else if len(btcAccounts) > 0 {
		acc := btcAccounts[0]
		fmt.Printf("BTC: 보유 %s (잠김 %s)\n", acc.Balance, acc.Locked)
		fmt.Printf("평균 매수가: %s\n", acc.AvgBuyPrice)
	}

	// Place a limit order (COMMENTED OUT - uncomment to test)
	/*
		fmt.Println("\n=== 지정가 주문 생성 ===")
		order, err := c.Private.PlaceOrder(&private.PlaceOrderRequest{
			Market:    "KRW-BTC",
			Side:      private.OrderSideBid,  // 매수
			OrderType: private.OrderTypeLimit,
			Price:     "50000000",           // 50,000,000 KRW
			Volume:    "0.001",              // 0.001 BTC
		})
		if err != nil {
			log.Printf("주문 실패: %v", err)
		} else {
			fmt.Printf("주문 생성됨: %s\n", order.OrderID)
		}
	*/

	fmt.Println("\n=== 예제 완료 ===")
}
