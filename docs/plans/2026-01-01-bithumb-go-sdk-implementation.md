# Bithumb Go SDK Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 빗썸 API 공식 래퍼 Go SDK로 Public/Private API와 WebSocket을 지원하는 타입 안전한 클라이언트 라이브러리 구현

**Architecture:** 계층형 아키텍처로 메인 클라이언트가 Public/Private/WebSocket 서브클라이언트를 포함하며, 각 서브클라이언트는 독립적으로 동작. JWT 인증은 Private 클라이언트에서 처리하며, WebSocket은 배치 구독 방식과 자동 재연결을 지원.

**Tech Stack:** Go 1.21+, net/http, coder/websocket, golang-jwt/jwt/v5, 표준 라이브러리 중심

---

## Task 1: 프로젝트 기본 구조 설정

**Files:**
- Create: `go.mod`
- Create: `errors.go`
- Create: `models/common.go`

**Step 1: go.mod 초기화**

```bash
cd /Users/code/workspace/bithumb/bithumb-go
go mod init github.com/hysuki/bithumb-go
```

**Step 2: 의존성 추가**

```bash
go get github.com/coder/websocket
go get github.com/golang-jwt/jwt/v5
```

**Step 3: 에러 타입 정의 작성** (`errors.go`)

```go
package bithumbgo

import (
    "fmt"
    "errors"
)

type ErrorType int

const (
    ErrorTypeNetwork   ErrorType = iota
    ErrorTypeHTTP
    ErrorTypeAPI
    ErrorTypeParse
    ErrorTypeWebSocket
)

type Error struct {
    Type       ErrorType
    Code       string
    Message    string
    HTTPStatus int
    Err        error
}

func (e *Error) Error() string {
    return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *Error) Unwrap() error {
    return e.Err
}

func IsAPIError(err error) bool {
    var e *Error
    return errors.As(err, &e) && e.Type == ErrorTypeAPI
}

func IsRateLimitError(err error) bool {
    var e *Error
    return errors.As(err, &e) && e.HTTPStatus == 429
}

func HasErrorCode(err error, code string) bool {
    var e *Error
    return errors.As(err, &e) && e.Code == code
}
```

**Step 4: 공통 모델 작성** (`models/common.go`)

```go
package models

import "time"

type APIResponse struct {
    Status   string      `json:"status"`
    Data     interface{} `json:"data"`
    ErrorMessage string  `json:"message,omitempty"`
}

type Timestamp struct {
    Unix int64 `json:"timestamp"`
}

func (t *Timestamp) Time() time.Time {
    return time.Unix(t.Unix/1000, (t.Unix%1000)*1000000)
}
```

**Step 5: 커밋**

```bash
git add go.mod errors.go models/common.go
git commit -m "feat: initialize project structure and error types"
```

---

## Task 2: 메인 클라이언트 구현

**Files:**
- Create: `client/client.go`

**Step 1: 테스트 작성**

`client/client_test.go`:
```go
package client

import (
    "testing"
    "time"
)

func TestNewClient(t *testing.T) {
    client, err := NewClient()
    if err != nil {
        t.Fatalf("NewClient() error = %v", err)
    }
    if client == nil {
        t.Fatal("NewClient() returned nil")
    }
    if client.HasAPIKey() {
        t.Error("NewClient() should not have API key by default")
    }
}

func TestNewClientWithOptions(t *testing.T) {
    client, err := NewClient(
        WithAPIKey("test-key", "test-secret"),
        WithTimeout(10*time.Second),
    )
    if err != nil {
        t.Fatalf("NewClient() error = %v", err)
    }
    if !client.HasAPIKey() {
        t.Error("NewClient(WithAPIKey()) should have API key")
    }
}
```

**Step 2: 테스트 실행 (실패 확인)**

```bash
go test ./client -v
```

Expected: FAIL - package not found

**Step 3: 클라이언트 구현** (`client/client.go`)

```go
package client

import (
    "net/http"
    "time"
)

type Client struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
    apiSecret  string

    HasAPIKeyFunc func() bool
}

func NewClient(opts ...Option) (*Client, error) {
    c := &Client{
        baseURL: "https://api.bithumb.com",
        httpClient: &http.Client{Timeout: 30 * time.Second},
    }

    for _, opt := range opts {
        opt(c)
    }

    c.HasAPIKeyFunc = func() bool {
        return c.apiKey != "" && c.apiSecret != ""
    }

    return c, nil
}

type Option func(*Client)

func WithAPIKey(apiKey, apiSecret string) Option {
    return func(c *Client) {
        c.apiKey = apiKey
        c.apiSecret = apiSecret
    }
}

func WithHTTPClient(hc *http.Client) Option {
    return func(c *Client) { c.httpClient = hc }
}

func WithBaseURL(url string) Option {
    return func(c *Client) { c.baseURL = url }
}

func WithTimeout(timeout time.Duration) Option {
    return func(c *Client) {
        c.httpClient.Timeout = timeout
    }
}

func (c *Client) HasAPIKey() bool {
    return c.HasAPIKeyFunc()
}

func (c *Client) BaseURL() string {
    return c.baseURL
}

func (c *Client) HTTPClient() *http.Client {
    return c.httpClient
}

func (c *Client) APIKey() string {
    return c.apiKey
}

func (c *Client) APISecret() string {
    return c.apiSecret
}
```

**Step 4: 테스트 실행**

```bash
go test ./client -v
```

Expected: PASS

**Step 5: 커밋**

```bash
git add client/client.go client/client_test.go
git commit -m "feat: implement main client with options pattern"
```

---

## Task 3: Public API 모델 정의

**Files:**
- Create: `models/public/request.go`
- Create: `models/public/ticker.go`
- Create: `models/public/orderbook.go`
- Create: `models/public/trade.go`
- Create: `models/public/candle.go`

**Step 1: Ticker 모델 작성** (`models/public/ticker.go`)

```go
package public

type Ticker struct {
    Market                 string  `json:"market"`
    TradeDate              string  `json:"trade_date"`
    TradeTime              string  `json:"trade_time"`
    TradeDateKST           string  `json:"trade_date_kst"`
    TradeTimeKST           string  `json:"trade_time_kst"`
    TradeTimestamp         int64   `json:"trade_timestamp"`
    OpeningPrice           float64 `json:"opening_price"`
    HighPrice              float64 `json:"high_price"`
    LowPrice               float64 `json:"low_price"`
    TradePrice             float64 `json:"trade_price"`
    PrevClosingPrice       float64 `json:"prev_closing_price"`
    Change                 string  `json:"change"`
    ChangePrice            float64 `json:"change_price"`
    ChangeRate             float64 `json:"change_rate"`
    SignedChangePrice      float64 `json:"signed_change_price"`
    SignedChangeRate       float64 `json:"signed_change_rate"`
    TradeVolume            float64 `json:"trade_volume"`
    AccTradePrice          float64 `json:"acc_trade_price"`
    AccTradePrice24h       float64 `json:"acc_trade_price_24h"`
    AccTradeVolume         float64 `json:"acc_trade_volume"`
    AccTradeVolume24h      float64 `json:"acc_trade_volume_24h"`
    Highest52WeekPrice     float64 `json:"highest_52_week_price"`
    Highest52WeekDate      string  `json:"highest_52_week_date"`
    Lowest52WeekPrice      float64 `json:"lowest_52_week_price"`
    Lowest52WeekDate       string  `json:"lowest_52_week_date"`
    Timestamp              int64   `json:"timestamp"`
}

type GetTickerRequest struct {
    Markets []string // KRW-BTC, BTC-ETH 등 (콤마로 구분)
}
```

**Step 2: OrderBook 모델 작성** (`models/public/orderbook.go`)

```go
package public

type OrderBook struct {
    Market         string          `json:"market"`
    Timestamp      int64           `json:"timestamp"`
    TotalAskSize   float64         `json:"total_ask_size"`
    TotalBidSize   float64         `json:"total_bid_size"`
    OrderBookUnits []OrderBookUnit `json:"orderbook_units"`
}

type OrderBookUnit struct {
    AskPrice float64 `json:"ask_price"`
    BidPrice float64 `json:"bid_price"`
    AskSize  float64 `json:"ask_size"`
    BidSize  float64 `json:"bid_size"`
}

type GetOrderBookRequest struct {
    Markets []string
}
```

**Step 3: Trade 모델 작성** (`models/public/trade.go`)

```go
package public

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

type GetRecentTradesRequest struct {
    Market  string
    To      string // HHmmss 또는 HH:mm:ss
    Count   int    // 1~500, 기본값 1
    Cursor  string // sequentialId
    DaysAgo int    // 1~7
}
```

**Step 4: Candle 모델 작성** (`models/public/candle.go`)

```go
package public

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

type GetCandlestickRequest struct {
    Market string
    Unit   CandleInterval // 1, 3, 5, 10, 15, 30, 60, 240
    To     string         // yyyy-MM-dd HH:mm:ss 또는 yyyy-MM-ddTHH:mm:ss
    Count  int            // 최대 200
}
```

**Step 5: 커밋**

```bash
git add models/public/
git commit -m "feat: define public API models"
```

---

## Task 4: Public API 클라이언트 구현

**Files:**
- Create: `public/client.go`
- Create: `public/client_test.go`

**Step 1: GetTicker 테스트 작성** (`public/client_test.go`)

```go
package public

import (
    "context"
    "testing"
    "time"

    "github.com/hysuki/bithumb-go/client"
    "github.com/hysuki/bithumb-go/models/public"
)

func TestGetTicker(t *testing.T) {
    baseClient, _ := client.NewClient()
    c := NewClient(baseClient)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    ticker, err := c.GetTickerWithContext(ctx, &public.GetTickerRequest{
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
```

**Step 2: 테스트 실행 (실패 확인)**

```bash
go test ./public -v -run TestGetTicker
```

Expected: FAIL - undefined: NewClient

**Step 3: Public 클라이언트 구현** (`public/client.go`)

```go
package public

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"

    "github.com/hysuki/bithumb-go/client"
    "github.com/hysuki/bithumb-go/models/public"
)

type Client struct {
    base *client.Client
}

func NewClient(base *client.Client) *Client {
    return &Client{base: base}
}

func (c *Client) buildURL(path string, params map[string]string) string {
    u := c.base.BaseURL() + path
    if len(params) > 0 {
        v := url.Values{}
        for k, val := range params {
            v.Add(k, val)
        }
        u += "?" + v.Encode()
    }
    return u
}

func (c *Client) do(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
    req, err := http.NewRequestWithContext(ctx, method, url, body)
    if err != nil {
        return nil, fmt.Errorf("create request: %w", err)
    }

    req.Header.Set("Accept", "application/json")

    return c.base.HTTPClient().Do(req)
}

func (c *Client) GetTicker(req *public.GetTickerRequest) ([]public.Ticker, error) {
    return c.GetTickerWithContext(context.Background(), req)
}

func (c *Client) GetTickerWithContext(ctx context.Context, req *public.GetTickerRequest) ([]public.Ticker, error) {
    params := map[string]string{"market": ""}
    if len(req.Markets) > 0 {
        params["market"] = req.Markets[0]
    }

    resp, err := c.do(ctx, http.MethodGet, c.buildURL("/v1/ticker", params), nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read response: %w", err)
    }

    var ticker public.Ticker
    if err := json.Unmarshal(body, &ticker); err != nil {
        return nil, fmt.Errorf("parse response: %w", err)
    }

    return []public.Ticker{ticker}, nil
}

func (c *Client) GetOrderBook(req *public.GetOrderBookRequest) (*public.OrderBook, error) {
    return c.GetOrderBookWithContext(context.Background(), req)
}

func (c *Client) GetOrderBookWithContext(ctx context.Context, req *public.GetOrderBookRequest) (*public.OrderBook, error) {
    params := map[string]string{"markets": ""}
    if len(req.Markets) > 0 {
        params["markets"] = req.Markets[0]
    }

    resp, err := c.do(ctx, http.MethodGet, c.buildURL("/v1/orderbook", params), nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read response: %w", err)
    }

    var orderbook public.OrderBook
    if err := json.Unmarshal(body, &orderbook); err != nil {
        return nil, fmt.Errorf("parse response: %w", err)
    }

    return &orderbook, nil
}

func (c *Client) GetRecentTrades(req *public.GetRecentTradesRequest) ([]public.Trade, error) {
    return c.GetRecentTradesWithContext(context.Background(), req)
}

func (c *Client) GetRecentTradesWithContext(ctx context.Context, req *public.GetRecentTradesRequest) ([]public.Trade, error) {
    params := map[string]string{
        "market": req.Market,
    }
    if req.To != "" {
        params["to"] = req.To
    }
    if req.Count > 0 {
        params["count"] = fmt.Sprintf("%d", req.Count)
    }

    resp, err := c.do(ctx, http.MethodGet, c.buildURL("/v1/trades/ticks", params), nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read response: %w", err)
    }

    var trades []public.Trade
    if err := json.Unmarshal(body, &trades); err != nil {
        return nil, fmt.Errorf("parse response: %w", err)
    }

    return trades, nil
}

func (c *Client) GetCandlestick(req *public.GetCandlestickRequest) ([]public.Candle, error) {
    return c.GetCandlestickWithContext(context.Background(), req)
}

func (c *Client) GetCandlestickWithContext(ctx context.Context, req *public.GetCandlestickRequest) ([]public.Candle, error) {
    params := map[string]string{"market": req.Market}
    if req.To != "" {
        params["to"] = req.To
    }
    if req.Count > 0 {
        params["count"] = fmt.Sprintf("%d", req.Count)
    }

    url := fmt.Sprintf("/v1/candles/minutes/%s", req.Unit)
    resp, err := c.do(ctx, http.MethodGet, c.buildURL(url, params), nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read response: %w", err)
    }

    var candles []public.Candle
    if err := json.Unmarshal(body, &candles); err != nil {
        return nil, fmt.Errorf("parse response: %w", err)
    }

    return candles, nil
}
```

**Step 4: 테스트 실행**

```bash
go test ./public -v
```

Expected: PASS (실제 API 호출)

**Step 5: 커밋**

```bash
git add public/
git commit -m "feat: implement public API client"
```

---

## Task 5: Private API JWT 인증 구현

**Files:**
- Create: `private/auth.go`
- Create: `private/auth_test.go`

**Step 1: JWT 생성 테스트 작성** (`private/auth_test.go`)

```go
package private

import (
    "testing"
    "time"

    "github.com/hysuki/bithumb-go/client"
)

func TestGenerateToken(t *testing.T) {
    base, _ := client.NewClient(
        client.WithAPIKey("test-key", "test-secret"),
    )
    c := NewClient(base)

    token, err := c.GenerateToken()
    if err != nil {
        t.Fatalf("GenerateToken() error = %v", err)
    }

    if token == "" {
        t.Fatal("GenerateToken() returned empty token")
    }

    // 토큰은 JWT 형식이어야 함 (dot으로 구분)
    if len(token) < 10 {
        t.Errorf("Token too short: %s", token)
    }
}
```

**Step 2: 테스트 실행 (실패 확인)**

```bash
go test ./private -v -run TestGenerateToken
```

Expected: FAIL - undefined: NewClient

**Step 3: JWT 인증 구현** (`private/auth.go`)

```go
package private

import (
    "fmt"

    "github.com/golang-jwt/jwt/v5"

    "github.com/hysuki/bithumb-go/client"
)

type Client struct {
    base *client.Client
}

func NewClient(base *client.Client) *Client {
    return &Client{base: base}
}

type TokenClaims struct {
    AccessKey string `json:"access_key"`
    Nonce     string `json:"nonce"`
    Timestamp int64  `json:"timestamp"`
    jwt.RegisteredClaims
}

func (c *Client) GenerateToken() (string, error) {
    if !c.base.HasAPIKey() {
        return "", fmt.Errorf("API key and secret required")
    }

    now := time.Now()
    claims := TokenClaims{
        AccessKey: c.base.APIKey(),
        Nonce:     generateNonce(),
        Timestamp: now.UnixMilli(),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(c.base.APISecret()))
}

func generateNonce() string {
    return uuid.New().String()
}
```

**Step 4: 테스트 실행**

```bash
go test ./private -v -run TestGenerateToken
```

Expected: PASS

**Step 5: 커밋**

```bash
git add private/auth.go private/auth_test.go go.mod
git commit -m "feat: implement JWT authentication for private API"
```

---

## Task 6: Private API 모델 및 클라이언트 구현

**Files:**
- Create: `models/private/request.go`
- Create: `models/private/account.go`
- Create: `models/private/order.go`
- Modify: `private/client.go`

**Step 1: Private 모델 작성** (`models/private/account.go`)

```go
package private

type Account struct {
    Currency             string  `json:"currency"`
    Balance              string  `json:"balance"`
    Locked               string  `json:"locked"`
    AvgBuyPrice          string  `json:"avg_buy_price"`
    AvgBuyPriceModified  bool    `json:"avg_buy_price_modified"`
    UnitCurrency         string  `json:"unit_currency"`
}

type GetAccountRequest struct {
    Currency string // 비워두면 전체
}
```

**Step 2: 주문 모델 작성** (`models/private/order.go`)

```go
package private

type Order struct {
    OrderID   string    `json:"order_id"`
    Market    string    `json:"market"`
    Side      string    `json:"side"`
    OrderType string    `json:"order_type"`
    CreatedAt string    `json:"created_at"`
}

type PlaceOrderRequest struct {
    Market    string  // KRW-BTC
    Side      string  // bid(매수), ask(매도)
    OrderType string  // limit(지정가), price(시장가 매수), market(시장가 매도)
    Price     string  // 지정가, 시장가 매수 시 필수
    Volume    string  // 지정가, 시장가 매도 시 필수
}

type CancelOrderRequest struct {
    UUID string // 주문 ID
}
```

**Step 3: Private 클라이언트 구현** (`private/client.go`)

```go
package private

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"

    "github.com/hysuki/bithumb-go/client"
    "github.com/hysuki/bithumb-go/models/private"
)

func (c *Client) doWithAuth(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
    token, err := c.GenerateToken()
    if err != nil {
        return nil, fmt.Errorf("generate token: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, method, url, body)
    if err != nil {
        return nil, fmt.Errorf("create request: %w", err)
    }

    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")

    return c.base.HTTPClient().Do(req)
}

func (c *Client) GetAccount(req *private.GetAccountRequest) ([]private.Account, error) {
    return c.GetAccountWithContext(context.Background(), req)
}

func (c *Client) GetAccountWithContext(ctx context.Context, req *private.GetAccountRequest) ([]private.Account, error) {
    url := c.base.BaseURL() + "/v1/accounts"

    resp, err := c.doWithAuth(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read response: %w", err)
    }

    var accounts []private.Account
    if err := json.Unmarshal(body, &accounts); err != nil {
        return nil, fmt.Errorf("parse response: %w", err)
    }

    return accounts, nil
}

func (c *Client) PlaceOrder(req *private.PlaceOrderRequest) (*private.Order, error) {
    return c.PlaceOrderWithContext(context.Background(), req)
}

func (c *Client) PlaceOrderWithContext(ctx context.Context, req *private.PlaceOrderRequest) (*private.Order, error) {
    body, _ := json.Marshal(req)
    url := c.base.BaseURL() + "/v2/orders"

    resp, err := c.doWithAuth(ctx, http.MethodPost, url, bytes.NewReader(body))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read response: %w", err)
    }

    var order private.Order
    if err := json.Unmarshal(respBody, &order); err != nil {
        return nil, fmt.Errorf("parse response: %w", err)
    }

    return &order, nil
}

func (c *Client) CancelOrder(req *private.CancelOrderRequest) error {
    return c.CancelOrderWithContext(context.Background(), req)
}

func (c *Client) CancelOrderWithContext(ctx context.Context, req *private.CancelOrderRequest) error {
    url := fmt.Sprintf("%s/v2/order/%s", c.base.BaseURL(), req.UUID)

    reqBody, _ := json.Marshal(map[string]string{})
    resp, err := c.doWithAuth(ctx, http.MethodDelete, url, bytes.NewReader(reqBody))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("cancel order failed: status %d", resp.StatusCode)
    }

    return nil
}
```

**Step 4: 커밋**

```bash
git add models/private/ private/client.go
git commit -m "feat: implement private API client with account and order"
```

---

## Task 7: WebSocket 기본 구조 및 타입 정의

**Files:**
- Create: `websocket/types.go`
- Create: `websocket/subscription.go`

**Step 1: WebSocket 타입 정의** (`websocket/types.go`)

```go
package websocket

type SubscriptionType string

const (
    SubscriptionTypeTicker     SubscriptionType = "ticker"
    SubscriptionTypeOrderBook  SubscriptionType = "orderbook"
    SubscriptionTypeTrade      SubscriptionType = "transaction"
    SubscriptionTypeMyOrder    SubscriptionType = "myOrder"
    SubscriptionTypeMyAsset    SubscriptionType = "myAsset"
)

type MessageHandler func(msg []byte) error

type MessageHandlers struct {
    Ticker    MessageHandler
    OrderBook MessageHandler
    Trade     MessageHandler
    MyOrder   MessageHandler
    MyAsset   MessageHandler
}

type TickerMessage struct {
    Type                 string  `json:"type"`
    Content              TickerContent `json:"content"`
}

type TickerContent struct {
    // API 문서 참조하여 필드 정의
}

type OrderBookMessage struct {
    Type    string `json:"type"`
    Content OrderBookContent `json:"content"`
}

type OrderBookContent struct {
    // API 문서 참조하여 필드 정의
}
```

**Step 2: 구독 관리자 구현** (`websocket/subscription.go`)

```go
package websocket

import (
    "encoding/json"
    "fmt"
    "sync"
    "time"

    "github.com/google/uuid"
)

type SubscriptionParam struct {
    Type       SubscriptionType `json:"type"`
    Codes      []string         `json:"codes,omitempty"`
}

type SubscriptionInfo struct {
    Type       SubscriptionType
    Codes      []string
    Ticket     string
    CreatedAt  time.Time
    IsActive   bool
}

type SubscriptionManager struct {
    subscriptions map[string]*SubscriptionInfo
    mu            sync.RWMutex
}

func NewSubscriptionManager() *SubscriptionManager {
    return &SubscriptionManager{
        subscriptions: make(map[string]*SubscriptionInfo),
    }
}

func (sm *SubscriptionManager) generateTicket() string {
    return uuid.New().String()
}

func (sm *SubscriptionManager) CreateSubscriptionMessage(params []*SubscriptionParam) ([]byte, string, error) {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    ticket := sm.generateTicket()
    messages := []interface{}{
        map[string]string{"ticket": ticket},
    }

    for _, p := range params {
        messages = append(messages, p)

        // 구독 정보 저장
        key := string(p.Type)
        sm.subscriptions[key] = &SubscriptionInfo{
            Type:      p.Type,
            Codes:     p.Codes,
            Ticket:    ticket,
            CreatedAt: time.Now(),
            IsActive:  true,
        }
    }

    body, err := json.Marshal(messages)
    return body, ticket, err
}

func (sm *SubscriptionManager) CreateUnsubscribeMessage() ([]byte, string, error) {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    // 모든 구독 비활성화
    for _, sub := range sm.subscriptions {
        sub.IsActive = false
    }

    ticket := sm.generateTicket()
    message := []interface{}{
        map[string]string{"ticket": ticket},
    }
    body, err := json.Marshal(message)
    return body, ticket, err
}

func (sm *SubscriptionManager) RestoreSubscriptions() ([]byte, string, error) {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    ticket := sm.generateTicket()
    messages := []interface{}{
        map[string]string{"ticket": ticket},
    }

    for _, sub := range sm.subscriptions {
        if sub.IsActive {
            param := &SubscriptionParam{
                Type:  sub.Type,
                Codes: sub.Codes,
            }
            messages = append(messages, param)
        }
    }

    body, err := json.Marshal(messages)
    return body, ticket, err
}

func (sm *SubscriptionManager) GetSubscriptions() []*SubscriptionInfo {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    result := make([]*SubscriptionInfo, 0, len(sm.subscriptions))
    for _, sub := range sm.subscriptions {
        result = append(result, sub)
    }
    return result
}
```

**Step 3: 커밋**

```bash
git add websocket/types.go websocket/subscription.go
go get github.com/google/uuid
git add go.mod go.sum
git commit -m "feat: define WebSocket types and subscription manager"
```

---

## Task 8: WebSocket 클라이언트 구현

**Files:**
- Create: `websocket/client.go`
- Create: `websocket/handler.go`

**Step 1: WebSocket 클라이언트 구현** (`websocket/client.go`)

```go
package websocket

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/coder/websocket"

    "github.com/hysuki/bithumb-go/client"
)

const (
    DefaultPublicURL    = "wss://ws-api.bithumb.com/websocket/v1"
    DefaultPrivateURL   = "wss://ws-api.bithumb.com/websocket/v1/private"
    DefaultReconnectDelay = 5 * time.Second
)

type Client struct {
    base       *client.Client
    conn       *websocket.Conn
    url        string
    subs       *SubscriptionManager
    handlers   map[SubscriptionType]MessageHandler
    done       chan struct{}
    mu         sync.RWMutex
    reconnect  bool
    reconnectDelay time.Duration
    isConnected bool
}

func NewClient(base *client.Client) *Client {
    return &Client{
        base:          base,
        url:          DefaultPublicURL,
        subs:         NewSubscriptionManager(),
        handlers:     make(map[SubscriptionType]MessageHandler),
        done:         make(chan struct{}),
        reconnect:    true,
        reconnectDelay: DefaultReconnectDelay,
    }
}

func (c *Client) SetPrivateURL() {
    c.url = DefaultPrivateURL
}

func (c *Client) SetReconnect(enabled bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.reconnect = enabled
}

func (c *Client) Connect(ctx context.Context) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    if c.isConnected {
        return fmt.Errorf("already connected")
    }

    // Private 연결 시 인증 헤더 추가를 위해 http.Request 헤더 사용
    // coder/websocket의 Dial은 http.Header를 받지 않으므로,
    // Dial 후 Auth 헤더를 설정하거나 별도 처리 필요
    // 여기서는 기본 Dial 사용 후 연결
    conn, _, err := websocket.Dial(ctx, c.url, nil)
    if err != nil {
        return fmt.Errorf("dial: %w", err)
    }

    // Private 연결 시 추가 인증 처리가 필요할 수 있음
    // (빗썸 API 문서에 따라 별도 JWT 토큰 전송 방식 확인 필요)

    c.conn = conn
    c.isConnected = true

    // 메시지 수신 시작
    go c.readLoop()
    go c.reconnectLoop()

    return nil
}

func (c *Client) readLoop() {
    for {
        select {
        case <-c.done:
            return
        default:
        }

        c.mu.RLock()
        conn := c.conn
        c.mu.RUnlock()

        if conn == nil {
            time.Sleep(100 * time.Millisecond)
            continue
        }

        // coder/websocket: Read(ctx)은 (messageType, data, error) 반환
        messageType, message, err := conn.Read(context.Background())
        if err != nil {
            c.mu.Lock()
            c.isConnected = false
            c.mu.Unlock()
            return
        }

        // TextMessage만 처리
        if messageType == websocket.MessageText {
            c.handleMessage(message)
        }
    }
}

func (c *Client) reconnectLoop() {
    ticker := time.NewTicker(c.reconnectDelay)
    defer ticker.Stop()

    for {
        select {
        case <-c.done:
            return
        case <-ticker.C:
            c.mu.RLock()
            connected := c.isConnected
            shouldReconnect := c.reconnect
            c.mu.RUnlock()

            if !connected && shouldReconnect {
                ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
                if err := c.Connect(ctx); err == nil {
                    // 구독 복원
                    c.RestoreSubscriptions()
                }
                cancel()
            }
        }
    }
}

func (c *Client) Subscribe(params []*SubscriptionParam, handlers MessageHandlers) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    // 핸들러 등록
    if handlers.Ticker != nil {
        c.handlers[SubscriptionTypeTicker] = handlers.Ticker
    }
    if handlers.OrderBook != nil {
        c.handlers[SubscriptionTypeOrderBook] = handlers.OrderBook
    }
    if handlers.Trade != nil {
        c.handlers[SubscriptionTypeTrade] = handlers.Trade
    }
    if handlers.MyOrder != nil {
        c.handlers[SubscriptionTypeMyOrder] = handlers.MyOrder
    }
    if handlers.MyAsset != nil {
        c.handlers[SubscriptionTypeMyAsset] = handlers.MyAsset
    }

    // 구독 메시지 전송
    body, _, err := c.subs.CreateSubscriptionMessage(params)
    if err != nil {
        return fmt.Errorf("create subscription message: %w", err)
    }

    if c.conn == nil {
        return fmt.Errorf("not connected")
    }

    // coder/websocket: Write(ctx, messageType, data)
    err = c.conn.Write(context.Background(), websocket.MessageText, body)
    if err != nil {
        return fmt.Errorf("send subscription: %w", err)
    }

    return nil
}

func (c *Client) Unsubscribe() error {
    c.mu.Lock()
    defer c.mu.Unlock()

    body, _, err := c.subs.CreateUnsubscribeMessage()
    if err != nil {
        return fmt.Errorf("create unsubscribe message: %w", err)
    }

    if c.conn != nil {
        err = c.conn.Write(context.Background(), websocket.MessageText, body)
        if err != nil {
            return fmt.Errorf("send unsubscribe: %w", err)
        }
    }

    return nil
}

func (c *Client) RestoreSubscriptions() error {
    body, _, err := c.subs.RestoreSubscriptions()
    if err != nil {
        return fmt.Errorf("create restore message: %w", err)
    }

    c.mu.Lock()
    defer c.mu.Unlock()

    if c.conn != nil {
        return c.conn.Write(context.Background(), websocket.MessageText, body)
    }

    return nil
}

func (c *Client) handleMessage(data []byte) error {
    // 간단한 파싱 - 타입별 핸들러 분기
    // 실제 구현에서는 handler.go에서 처리
    return nil
}

func (c *Client) Close() error {
    close(c.done)

    c.mu.Lock()
    defer c.mu.Unlock()

    if c.conn != nil {
        // coder/websocket: Close(status, reason)
        err := c.conn.Close(websocket.StatusNormalClosure, "client closing")
        c.conn = nil
        c.isConnected = false
        return err
    }

    return nil
}
```

**Step 2: 메시지 핸들러 구현** (`websocket/handler.go`)

```go
package websocket

import (
    "encoding/json"
    "fmt"
)

func (c *Client) handleMessage(data []byte) error {
    // 타입 확인을 위한 간단한 파싱
    var raw map[string]interface{}
    if err := json.Unmarshal(data, &raw); err != nil {
        return fmt.Errorf("parse message: %w", err)
    }

    msgType, ok := raw["type"].(string)
    if !ok {
        return fmt.Errorf("missing type field")
    }

    c.mu.RLock()
    handler, exists := c.handlers[SubscriptionType(msgType)]
    c.mu.RUnlock()

    if !exists {
        return nil // 처리할 핸들러 없음
    }

    return handler(data)
}
```

**Step 3: 커밋**

```bash
git add websocket/client.go websocket/handler.go
git commit -m "feat: implement WebSocket client with auto-reconnect"
```

---

## Task 9: 통합 클라이언트 연결

**Files:**
- Modify: `client/client.go`

**Step 1: 메인 클라이언트 수정** (`client/client.go`)

```go
package client

import (
    "github.com/hysuki/bithumb-go/public"
    "github.com/hysuki/bithumb-go/private"
    "github.com/hysuki/bithumb-go/websocket"
)

type Client struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
    apiSecret  string

    HasAPIKeyFunc func() bool

    Public    *public.Client
    Private   *private.Client
    Websocket *websocket.Client
}

func NewClient(opts ...Option) (*Client, error) {
    c := &Client{
        baseURL:    "https://api.bithumb.com",
        httpClient: &http.Client{Timeout: 30 * time.Second},
    }

    for _, opt := range opts {
        opt(c)
    }

    c.HasAPIKeyFunc = func() bool {
        return c.apiKey != "" && c.apiSecret != ""
    }

    c.Public = public.NewClient(c)

    if c.HasAPIKey() {
        c.Private = private.NewClient(c)
    }

    c.Websocket = websocket.NewClient(c)

    return c, nil
}
```

**Step 2: 사용 예제 작성** (`examples/basic_usage.go`)

```go
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
    // Public API만 사용
    c, _ := client.NewClient()

    // Ticker 조회
    ticker, err := c.Public.GetTicker(&public.GetTickerRequest{
        Markets: []string{"KRW-BTC"},
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("BTC Price: %f\n", ticker[0].TradePrice)

    // Context 사용
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    orderbook, err := c.Public.GetOrderBookWithContext(ctx, &public.GetOrderBookRequest{
        Markets: []string{"KRW-BTC"},
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("OrderBook: %d levels\n", len(orderbook.OrderBookUnits))

    // Private API 사용 (API 키 필요)
    // cWithAuth, _ := client.NewClient(
    //     client.WithAPIKey("your-key", "your-secret"),
    // )
    // accounts, _ := cWithAuth.Private.GetAccount(&private.GetAccountRequest{})
}
```

**Step 3: 커밋**

```bash
git add client/client.go examples/basic_usage.go
git commit -m "feat: integrate sub-clients into main client"
```

---

## Task 10: README 및 문서 작성

**Files:**
- Create: `README.md`
- Create: `CONTRIBUTING.md`

**Step 1: README 작성** (`README.md`)

```markdown
# Bithumb Go SDK

빗썸 API 공식 래퍼 Go SDK

## 설치

```bash
go get github.com/hysuki/bithumb-go
```

## 사용법

### Public API

```go
package main

import (
    "fmt"
    "log"

    "github.com/hysuki/bithumb-go/client"
    "github.com/hysuki/bithumb-go/models/public"
)

func main() {
    c, _ := client.NewClient()

    ticker, err := c.Public.GetTicker(&public.GetTickerRequest{
        Markets: []string{"KRW-BTC"},
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("BTC: %f\n", ticker[0].TradePrice)
}
```

### Private API

```go
c, _ := client.NewClient(
    client.WithAPIKey("your-key", "your-secret"),
)

accounts, err := c.Private.GetAccount(&private.GetAccountRequest{})
```

### WebSocket

```go
c, _ := client.NewClient()

ws := c.Websocket
ws.Connect(context.Background())

params := []*websocket.SubscriptionParam{
    {Type: websocket.SubscriptionTypeTicker, Codes: []string{"KRW-BTC"}},
}

handlers := websocket.MessageHandlers{
    Ticker: func(msg []byte) error {
        fmt.Printf("Ticker: %s\n", string(msg))
        return nil
    },
}

ws.Subscribe(params, handlers)
```

## 기능

- [x] Public API (Ticker, OrderBook, Trade, Candlestick)
- [x] Private API (Account, Order)
- [x] WebSocket (Ticker, OrderBook, Trade, MyOrder, MyAsset)
- [x] JWT 인증
- [x] 자동 재연결
- [x] Context 지원

## 라이선스

MIT
```

**Step 2: 커밋**

```bash
git add README.md CONTRIBUTING.md
git commit -m "docs: add README and contributing guide"
```

---

## 완료 체크리스트

- [x] 프로젝트 기본 구조 설정
- [x] 메인 클라이언트 구현 (Options 패턴)
- [x] 에러 타입 정의
- [x] Public API 모델 정의
- [x] Public API 클라이언트 구현
- [x] Private API JWT 인증 구현
- [x] Private API 모델 및 클라이언트 구현
- [x] WebSocket 타입 및 구독 관리자 구현
- [x] WebSocket 클라이언트 구현 (자동 재연결)
- [x] 통합 클라이언트 연결
- [x] README 및 문서 작성

---

## 테스트 방법

```bash
# 모든 테스트
go test ./...

# Public API 테스트 (실제 API 호출)
go test ./public -v

# JWT 테스트
go test ./private -v -run TestGenerateToken
```

## 다음 단계

계획 실행을 위해 `superpowers:executing-plans` 스킬을 사용하세요.
