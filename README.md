# Bithumb Go SDK

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

빗썸(Bithumb) 거래소 API를 위한 공식 Go SDK입니다. REST API와 WebSocket을 지원하며, 타입 안전성과 동시성 안전성을 보장합니다.

## 특징

- **공개 API 지원**: Ticker, OrderBook, 체결 내역, 캔들 데이터
- **전용 API 지원**: 계정 정보, 주문 조회/생성/취소 (JWT 인증)
- **WebSocket 실시간 데이터**: Ticker, OrderBook, 체결, 주문/자산 변경 실시간 구독
- **자동 재연결**: 연결 끊김 시 자동 재연결 및 구독 복구
- **스레드 안전**: 동시 요청 처리에 안전한 클라이언트 설계
- **Options 패턴**: 유연한 클라이언트 설정
- **컨텍스트 지원**: 요청 취소 및 타임아웃 제어
- **포괄적인 테스트**: 단위 테스트 및 통합 테스트 제공

## 설치

```bash
go get github.com/bithumb-go/bithumb-go
```

## 시작하기

### 공개 API (Public API)

API 키 없이 시장 데이터를 조회할 수 있습니다.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/bithumb-go/bithumb-go/client"
    "github.com/bithumb-go/bithumb-go/models/public"
)

func main() {
    // 클라이언트 생성
    c, err := client.NewClient()
    if err != nil {
        log.Fatalf("클라이언트 생성 실패: %v", err)
    }

    // Ticker 조회
    tickers, err := c.Public.GetTicker(&public.GetTickerRequest{
        Markets: []string{"KRW-BTC", "KRW-ETH"},
    })
    if err != nil {
        log.Printf("Ticker 조회 실패: %v", err)
        return
    }

    for _, ticker := range tickers {
        fmt.Printf("%s: %.2f KRW (변동율: %.2f%%)\n",
            ticker.Market, ticker.TradePrice, ticker.ChangeRate*100)
    }

    // 컨텍스트와 함께 사용 (타임아웃 설정)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    orderbook, err := c.Public.GetOrderBookWithContext(ctx, &public.GetOrderBookRequest{
        Markets: []string{"KRW-BTC"},
    })
    if err != nil {
        log.Printf("OrderBook 조회 실패: %v", err)
        return
    }

    fmt.Printf("호가 매도 총량: %.4f BTC\n", orderbook.TotalAskSize)
    fmt.Printf("호가 매수 총량: %.4f BTC\n", orderbook.TotalBidSize)
}
```

### 전용 API (Private API)

API 키가 필요한 계정 및 주문 관련 기능입니다.

```go
package main

import (
    "fmt"
    "log"

    "github.com/bithumb-go/bithumb-go/client"
    "github.com/bithumb-go/bithumb-go/models/private"
)

func main() {
    // API 키로 클라이언트 생성
    c, err := client.NewClient(
        client.WithAPIKey("your-api-key", "your-api-secret"),
    )
    if err != nil {
        log.Fatalf("클라이언트 생성 실패: %v", err)
    }

    // 계정 정보 조회
    accounts, err := c.Private.GetAccount(&private.GetAccountRequest{})
    if err != nil {
        log.Printf("계정 조회 실패: %v", err)
        return
    }

    for _, acc := range accounts {
        fmt.Printf("%s: 보유 %.8f (잠김 %.8f)\n",
            acc.Currency, acc.Balance, acc.Locked)
    }

    // 지정가 주문 생성
    order, err := c.Private.PlaceOrder(&private.PlaceOrderRequest{
        Market:    "KRW-BTC",
        Side:      private.OrderSideBid,  // 매수
        OrderType: private.OrderTypeLimit,
        Price:     "50000000",           // 50,000,000 KRW
        Volume:    "0.001",              // 0.001 BTC
    })
    if err != nil {
        log.Printf("주문 실패: %v", err)
        return
    }

    fmt.Printf("주문 생성됨: %s\n", order.OrderID)

    // 주문 취소
    err = c.Private.CancelOrder(&private.CancelOrderRequest{
        UUID: order.OrderID,
    })
    if err != nil {
        log.Printf("주문 취소 실패: %v", err)
        return
    }

    fmt.Println("주문 취소됨")
}
```

### WebSocket 실시간 데이터

실시간 시장 데이터 및 계정 업데이트를 구독할 수 있습니다.

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/bithumb-go/bithumb-go/client"
    "github.com/bithumb-go/bithumb-go/websocket"
)

func main() {
    c, _ := client.NewClient()
    ws := c.Websocket

    // WebSocket 핸들러 등록
    handlers := websocket.MessageHandlers{
        Ticker: func(msg []byte) error {
            var ticker websocket.TickerMessage
            if err := json.Unmarshal(msg, &ticker); err != nil {
                return err
            }
            fmt.Printf("Ticker 업데이트: %s\n", ticker.Content.MarketCode)
            return nil
        },
        OrderBook: func(msg []byte) error {
            fmt.Printf("OrderBook 업데이트: %s\n", string(msg))
            return nil
        },
    }

    // WebSocket 연결
    ctx := context.Background()
    if err := ws.Connect(ctx); err != nil {
        log.Fatalf("WebSocket 연결 실패: %v", err)
    }
    defer ws.Close()

    // 구독 설정
    params := []*websocket.SubscriptionParam{
        {
            Type:  websocket.SubscriptionTypeTicker,
            Codes: []string{"KRW-BTC", "KRW-ETH"},
        },
        {
            Type:  websocket.SubscriptionTypeOrderBook,
            Codes: []string{"KRW-BTC"},
        },
    }

    if err := ws.Subscribe(params, handlers); err != nil {
        log.Fatalf("구독 실패: %v", err)
    }

    // 프로그램이 종료되지 않도록 대기
    select {}
}
```

### 전용 WebSocket (Private WebSocket)

개인 주문 및 자산 변경 실시간 알림을 구독합니다.

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/bithumb-go/bithumb-go/client"
    "github.com/bithumb-go/bithumb-go/websocket"
)

func main() {
    // API 키로 클라이언트 생성
    c, err := client.NewClient(
        client.WithAPIKey("your-api-key", "your-api-secret"),
    )
    if err != nil {
        log.Fatalf("클라이언트 생성 실패: %v", err)
    }

    ws := c.Websocket
    ws.SetPrivateURL()  // 전용 WebSocket URL 설정

    // 핸들러 등록
    handlers := websocket.MessageHandlers{
        MyOrder: func(msg []byte) error {
            fmt.Printf("내 주문 업데이트: %s\n", string(msg))
            return nil
        },
        MyAsset: func(msg []byte) error {
            fmt.Printf("내 자산 업데이트: %s\n", string(msg))
            return nil
        },
    }

    // 연결 및 구독
    ctx := context.Background()
    if err := ws.Connect(ctx); err != nil {
        log.Fatalf("WebSocket 연결 실패: %v", err)
    }
    defer ws.Close()

    params := []*websocket.SubscriptionParam{
        {Type: websocket.SubscriptionTypeMyOrder},
        {Type: websocket.SubscriptionTypeMyAsset},
    }

    if err := ws.Subscribe(params, handlers); err != nil {
        log.Fatalf("구독 실패: %v", err)
    }

    select {}
}
```

## 클라이언트 옵션

```go
c, err := client.NewClient(
    client.WithAPIKey("your-key", "your-secret"),  // API 키 설정
    client.WithBaseURL("https://api.bithumb.com"), // 기본 URL 설정
    client.WithTimeout(30*time.Second),            // 타임아웃 설정
)

// 커스텀 HTTP 클라이언트 사용
httpClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 100,
    },
}
c, _ := client.NewClient(client.WithHTTPClient(httpClient))
```

## API 참조

### 공개 API 메서드

| 메서드 | 설명 |
|--------|------|
| `GetTicker(req)` | 현재가 정보 조회 |
| `GetOrderBook(req)` | 호가 정보 조회 |
| `GetRecentTrades(req)` | 최근 체결 내역 조회 |
| `GetCandlestick(req)` | 캔들 데이터 조회 |

### 전용 API 메서드

| 메서드 | 설명 |
|--------|------|
| `GetAccount(req)` | 계정 정보 조회 |
| `PlaceOrder(req)` | 주문 생성 |
| `CancelOrder(req)` | 주문 취소 |

### WebSocket 구독 타입

| 타입 | 설명 |
|------|------|
| `SubscriptionTypeTicker` | 실시간 현재가 |
| `SubscriptionTypeOrderBook` | 실시간 호가 |
| `SubscriptionTypeTrade` | 실시간 체결 |
| `SubscriptionTypeMyOrder` | 내 주문 변경 (전용) |
| `SubscriptionTypeMyAsset` | 내 자산 변경 (전용) |

## 에러 처리

```go
tickers, err := c.Public.GetTicker(&public.GetTickerRequest{
    Markets: []string{"KRW-BTC"},
})

if err != nil {
    // 에러 타입 확인
    var apiErr *bithumb.APIError
    if errors.As(err, &apiErr) {
        log.Printf("API 에러: 코드=%d, 메시지=%s",
            apiErr.Code, apiErr.Message)
        return
    }

    // 기타 에러 처리
    log.Printf("요청 실패: %v", err)
    return
}
```

## 재연결 설정

```go
ws := c.Websocket

// 자동 재연결 활성화/비활성화 (기본: 활성화)
ws.SetReconnect(true)

// 재연결 지연 설정 (기본: 5초)
ws.SetReconnectDelay(10 * time.Second)

// 재연결 타임아웃 설정 (기본: 10초)
ws.SetReconnectTimeout(15 * time.Second)
```

## 요구사항

- Go 1.21 이상

## 라이선스

이 프로젝트는 MIT 라이선스 하에 배포됩니다. [LICENSE](LICENSE) 파일을 참조하세요.

## 지원

- 버전 노트: [CHANGELOG.md](CHANGELOG.md)
- 기여 가이드: [CONTRIBUTING.md](CONTRIBUTING.md)
- 문제 신고: [GitHub Issues](https://github.com/bithumb-go/bithumb-go/issues)

## 공식 링크

- [빗썸 API 문서](https://apidocs.bithumb.com)
- [빗썸 거래소](https://www.bithumb.com)

---

**주의**: 이 SDK는 공식적으로 빗썸에서 지원되지 않습니다. 사용 시 자신의 책임하에 사용하세요.
