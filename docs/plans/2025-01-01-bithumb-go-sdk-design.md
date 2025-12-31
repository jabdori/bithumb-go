# Bithumb Go SDK Design

## Overview

Bithumb API 공식 래퍼 SDK로, Public/Private API와 WebSocket을 지원합니다.

## Design Goals

- **타입 안전성**: 강 타입 구조체로 응답 처리
- **일관성**: 모든 요청 파라미터를 구조체로 캕슐화
- **유연성**: API 키 선택사항, Context 선택적 사용
- **신뢰성**: WebSocket 자동 재연결 + 재구독

## Directory Structure

```
bithumb-go/
├── client/
│   └── client.go          # 메인 Client, NewClient(), Options
├── public/
│   └── client.go          # Public API 메서드
├── private/
│   ├── client.go          # Private API 메서드
│   └── auth.go            # JWT 토큰 생성
├── websocket/
│   ├── client.go          # WebSocket 연결 관리
│   ├── subscription.go    # 배치 구독 관리
│   └── handler.go         # 메시지 핸들링
├── models/
│   ├── common.go          # 공통 타입, API 응답 기본 구조
│   ├── public/
│   │   ├── request.go     # Public API 요청 구조체
│   │   ├── ticker.go      # Ticker 관련 타입
│   │   ├── orderbook.go   # 호가 관련 타입
│   │   ├── trade.go       # 체결 관련 타입
│   │   └── candle.go      # 캔들 관련 타입
│   ├── private/
│   │   ├── request.go     # Private API 요청 구조체
│   │   ├── account.go     # 계좌/잔고 타입
│   │   └── order.go       # 주문 타입
│   └── websocket/
│       ├── types.go       # WebSocket 메시지 타입
│       └── subscription.go # 구독 관련 타입
└── errors.go              # 통합 에러 타입, 헬퍼 함수
```

## 1. Client Structure

### 메인 클라이언트

```go
package client

type Client struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string  // 선택사항
    apiSecret  string  // 선택사항

    Public    *public.Client
    Private   *private.Client  // API 키 있을 때만
    Websocket *websocket.Client
}

// NewClient는 API 키 없이도 생성 가능
func NewClient(opts ...Option) (*Client, error) {
    client := &Client{
        baseURL:    "https://api.bithumb.com",
        httpClient: &http.Client{Timeout: 30 * time.Second},
    }

    for _, opt := range opts {
        opt(client)
    }

    client.Public = public.NewClient(client)

    if client.HasAPIKey() {
        client.Private = private.NewClient(client)
    }

    client.Websocket = websocket.NewClient(client)

    return client, nil
}

// Option 설정 함수
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
    return func(c *Client) { c.httpClient.Timeout = timeout }
}

func (c *Client) HasAPIKey() bool {
    return c.apiKey != "" && c.apiSecret != ""
}
```

### 사용 예시

```go
// Public API만
client, _ := bithumb.NewClient()

// Private API까지
client, _ := bithumb.NewClient(
    bithumb.WithAPIKey("api-key", "api-secret"),
    bithumb.WithTimeout(10*time.Second),
)
```

## 2. Public API

### 메서드

```go
package public

type Client struct {
    base *client.Client
}

func NewClient(base *client.Client) *Client {
    return &Client{base: base}
}

// 기본 메서드
func (c *Client) GetTicker(req *models.GetTickerRequest) (*models.Ticker, error) {
    return c.GetTickerWithContext(context.Background(), req)
}

// Context 버전
func (c *Client) GetTickerWithContext(ctx context.Context, req *models.GetTickerRequest) (*models.Ticker, error) {
}

func (c *Client) GetOrderBook(req *models.GetOrderBookRequest) (*models.OrderBook, error) {
}

func (c *Client) GetOrderBookWithContext(ctx context.Context, req *models.GetOrderBookRequest) (*models.OrderBook, error) {
}

func (c *Client) GetRecentTrades(req *models.GetRecentTradesRequest) ([]models.Trade, error) {
}

func (c *Client) GetRecentTradesWithContext(ctx context.Context, req *models.GetRecentTradesRequest) ([]models.Trade, error) {
}

func (c *Client) GetCandlestick(req *models.GetCandlestickRequest) ([]models.Candle, error) {
}

func (c *Client) GetCandlestickWithContext(ctx context.Context, req *models.GetCandlestickRequest) ([]models.Candle, error) {
}
```

### 요청 구조체

```go
package public

type GetTickerRequest struct {
    Currency string  // BTC, ETH...
}

type GetOrderBookRequest struct {
    Currency     string
    OrderBookType string  // 1: all, 0: 일부
}

type GetRecentTradesRequest struct {
    Currency string
    Count    int  // default 50
}

type GetCandlestickRequest struct {
    Currency string
    Interval CandleInterval  // 1m, 5m, 30m, 1h, 1d...
}
```

### 사용 예시

```go
ticker, _ := client.Public.GetTicker(&public.GetTickerRequest{Currency: "BTC"})
orderbook, _ := client.Public.GetOrderBook(&public.GetOrderBookRequest{Currency: "BTC"})

// Context 필요 시
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
ticker, _ := client.Public.GetTickerWithContext(ctx, &public.GetTickerRequest{Currency: "BTC"})
```

## 3. Private API

### 메서드

```go
package private

type Client struct {
    base *client.Client
}

func NewClient(base *client.Client) *Client {
    return &Client{base: base}
}

// 계좌 조회
func (c *Client) GetAccount(req *models.GetAccountRequest) (*models.Account, error) {
}

func (c *Client) GetAccountWithContext(ctx context.Context, req *models.GetAccountRequest) (*models.Account, error) {
}

// 잔고 조회
func (c *Client) GetBalance(req *models.GetBalanceRequest) ([]models.Balance, error) {
}

func (c *Client) GetBalanceWithContext(ctx context.Context, req *models.GetBalanceRequest) ([]models.Balance, error) {
}

// 주문
func (c *Client) PlaceOrder(req *models.PlaceOrderRequest) (*models.Order, error) {
}

func (c *Client) PlaceOrderWithContext(ctx context.Context, req *models.PlaceOrderRequest) (*models.Order, error) {
}

// 주문 취소
func (c *Client) CancelOrder(req *models.CancelOrderRequest) error {
}

func (c *Client) CancelOrderWithContext(ctx context.Context, req *models.CancelOrderRequest) error {
}
```

### JWT 인증

```go
package private

func (c *Client) generateToken() (string, error) {
    claims := jwt.MapClaims{
        "access_key": c.base.apiKey,
        "nonce":      time.Now().UnixMilli(),
        "exp":        time.Now().Add(1 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(c.base.apiSecret))
}
```

### 요청 구조체

```go
package private

type GetAccountRequest struct {
    Currency string  // 비워두면 전체
}

type GetBalanceRequest struct {
    Currency string  // 비워두면 전체
    IsDetail bool    // 상세 조회 여부
}

type PlaceOrderRequest struct {
    OrderType   string  // bid/ask
    Currency    string
    Price       float64
    Quantity    float64
}

type CancelOrderRequest struct {
    OrderID string
    Currency string
}
```

### 사용 예시

```go
account, _ := client.Private.GetAccount(&private.GetAccountRequest{Currency: "BTC"})
balance, _ := client.Private.GetBalance(&private.GetBalanceRequest{})
order, _ := client.Private.PlaceOrder(&private.PlaceOrderRequest{
    OrderType: "bid",
    Currency:  "BTC",
    Price:     50000000,
    Quantity:  0.001,
})
```

## 4. WebSocket

### 배치 구독 방식

Bithumb WebSocket은 하나의 티켓으로 여러 구독을 처리합니다. 새 티켓을 발급하면 기존 구독이 모두 제거됩니다.

```go
package websocket

type SubscriptionType string

const (
    SubscriptionTypeTicker     SubscriptionType = "ticker"
    SubscriptionTypeOrderBook  SubscriptionType = "orderbook"
    SubscriptionTypeTrade      SubscriptionType = "transaction"
    SubscriptionTypeMyOrder    SubscriptionType = "myOrder"       // Private
    SubscriptionTypeMyAsset    SubscriptionType = "myAsset"       // Private
)

type SubscriptionParam struct {
    Type       SubscriptionType `json:"type"`
    Currencies []string         `json:"symbols,omitempty"`
}

type SubscriptionInfo struct {
    Type       SubscriptionType
    Currencies []string
    CreatedAt  time.Time
    IsActive   bool
}

type SubscriptionManager struct {
    subscriptions map[string]*SubscriptionInfo
    mu            sync.RWMutex
}

// CreateSubscriptionMessage는 배치 구독 메시지 생성
func (sm *SubscriptionManager) CreateSubscriptionMessage(params []*SubscriptionParam) ([]byte, string, error) {
    ticket := sm.generateTicket()
    messages := []interface{}{
        map[string]string{"ticket": ticket},
    }

    for _, p := range params {
        messages = append(messages, p)
    }

    body, err := json.Marshal(messages)
    return body, ticket, err
}

// CreateUnsubscribeMessage는 빈 티켓으로 모든 구독 해제
func (sm *SubscriptionManager) CreateUnsubscribeMessage() ([]byte, string, error) {
    ticket := sm.generateTicket()
    message := []interface{}{
        map[string]string{"ticket": ticket},
    }
    body, err := json.Marshal(message)
    return body, ticket, err
}

// RestoreSubscriptions는 재연결 시 구독 복원
func (sm *SubscriptionManager) RestoreSubscriptions() ([]byte, string, error) {
    // 저장된 구독 목록으로 새 티켓 생성
}
```

### WebSocket 클라이언트

```go
package websocket

type Client struct {
    base      *client.Client
    conn      *websocket.Conn
    url       string
    subs      *SubscriptionManager
    handlers  map[SubscriptionType]MessageHandler
    done      chan struct{}
    reconnect bool
}

func NewClient(base *client.Client) *Client {
    return &Client{
        base:      base,
        url:       "wss://pubweb.bithumb.com/v1/WS",
        subs:      NewSubscriptionManager(),
        handlers:  make(map[SubscriptionType]MessageHandler),
        reconnect: true,
        done:      make(chan struct{}),
    }
}

// Subscribe는 배치 구독
func (c *Client) Subscribe(params []*SubscriptionParam, handlers MessageHandlers) error {
}

// Unsubscribe는 빈 티켓으로 모든 구독 해제
func (c *Client) Unsubscribe() error {
}

// reconnectAndResubscribe 재연결 후 구독 복원
func (c *Client) reconnectAndResubscribe(ctx context.Context) error {
}
```

### 메시지 핸들링

```go
package websocket

type MessageHandler func(msg []byte) error

type MessageHandlers struct {
    Ticker     MessageHandler
    OrderBook  MessageHandler
    Trade      MessageHandler
    MyOrder    MessageHandler  // Private
    MyAsset    MessageHandler  // Private
}

func (c *Client) handleMessage(data []byte) error {
    // 타입별 핸들러 분기
}
```

### 사용 예시

```go
// Public 구독 (API 키 불필요)
params := []*websocket.SubscriptionParam{
    {Type: websocket.SubscriptionTypeTicker, Currencies: []string{"BTC", "ETH", "XRP"}},
    {Type: websocket.SubscriptionTypeOrderBook, Currencies: []string{"BTC"}},
}

handlers := websocket.MessageHandlers{
    Ticker: func(msg []byte) error {
        ticker := &models.TickerMessage{}
        json.Unmarshal(msg, ticker)
        fmt.Printf("BTC: %f\n", ticker.ClosePrice)
        return nil
    },
    OrderBook: func(msg []byte) error {
        // 처리
        return nil
    },
}

ws.Subscribe(params, handlers)

// Private 구독 (API 키 필요)
privateParams := []*websocket.SubscriptionParam{
    {Type: websocket.SubscriptionTypeMyOrder},
}
ws.Subscribe(privateParams, privateHandlers)

// 모든 구독 해제
ws.Unsubscribe()
```

## 5. Error Handling

### 통합 에러 타입

```go
package bithumb

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
    Code       string   // API 에러 코드
    Message    string   // 에러 메시지
    HTTPStatus int      // HTTP 상태 코드
    Err        error    // 원본 에러
}

func (e *Error) Error() string {
    return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *Error) Unwrap() error {
    return e.Err
}
```

### 헬퍼 함수

```go
package bithumb

func IsAPIError(err error) bool {
    var e *Error
    return errors.As(err, &e) && e.Type == ErrorTypeAPI
}

func IsRateLimitError(err error) bool {
    var e *Error
    return errors.As(err, &e) && e.HTTPStatus == 429
}

func IsNetworkError(err error) bool {
    var e *Error
    return errors.As(err, &e) && e.Type == ErrorTypeNetwork
}

func HasErrorCode(err error, code string) bool {
    var e *Error
    return errors.As(err, &e) && e.Code == code
}
```

### 사용 예시

```go
ticker, err := client.Public.GetTicker(&public.GetTickerRequest{Currency: "BTC"})
if err != nil {
    if bithumb.IsRateLimitError(err) {
        time.Sleep(time.Second)
    } else if bithumb.IsAPIError(err) {
        log.Printf("API error: %v", err)
    }
}
```

## API Specification

### Public API

| 메서드 | 설명 |
|--------|------|
| GetTicker | 현재가 정보 |
| GetOrderBook | 호가 정보 |
| GetRecentTrades | 체결 내역 |
| GetCandlestick | 캔들 데이터 |

### Private API

| 메서드 | 설명 |
|--------|------|
| GetAccount | 계좌 조회 |
| GetBalance | 잔고 조회 |
| PlaceOrder | 주문 |
| CancelOrder | 주문 취소 |
| GetOrderDetail | 주문 상세 |
| GetOrders | 주문 목록 |

### WebSocket

| 타입 | 설명 | 인증 |
|------|------|------|
| ticker | 현재가 | 불필요 |
| orderbook | 호가 | 불필요 |
| transaction | 체결 | 불필요 |
| myOrder | 내 주문 | 필요 |
| myAsset | 내 자산 | 필요 |

## Implementation Checklist

- [ ] Client 기본 구조
- [ ] Errors 정의 및 헬퍼 함수
- [ ] Models 정의 (common, public, private, websocket)
- [ ] Public API 구현
- [ ] Private API + JWT 구현
- [ ] WebSocket 연결 관리
- [ ] WebSocket 구독 관리 (배치)
- [ ] WebSocket 재연결 + 재구독
- [ ] 테스트 코드
