# bithumb-go SDK 개선 설계

## 개요

빗썸 Open API v2.1.5 문서에 기반하여, 현재 미구현된 API 기능을 추가하고 코드 품질을 개선하는 설계 문서입니다.

**작성일자**: 2026-01-01
**API 버전**: v2.1.5

---

## 1. 현재 상태

### 1.1 구현된 API (13개)

| 카테고리 | API | 수량 |
|---------|-----|------|
| PUBLIC | Market, Ticker, OrderBook, Trades, Candles(4종) | 7개 |
| PRIVATE | Account, Order(5종), OrderChance | 6개 |

### 1.2 미구현 API (35개)

| 카테고리 | 미구현 API |
|---------|-----------|
| PUBLIC | 경보제, 공지사항, 입출금수수료 (3개) |
| PRIVATE | TWAP(3), 입출금(15), 서비스정보(2), 코인대여(2), 기타(13개) |

### 1.3 코드 품질 이슈

| 이슈 | 설명 |
|------|------|
| 에러 처리 불일치 | PUBLIC은 plain `error`, PRIVATE는 `bithumbgo.Error` |
| 쿼리 파라미터 빌드 | 수동 문자열 연결로 버그 위험 |
| HTTP status check | PUBLIC API 일부 메서드에 누락 |

---

## 2. 개선 우선순위

### Phase 1: 웹소켓 에러 처리 (최우선)
### Phase 2: PUBLIC API 누락 (3개)
### Phase 3: PRIVATE API 핵심 (11개)
### Phase 4: 코드 품질 개선

---

## 3. Phase 1: 웹소켓 에러 처리

### 3.1 에러 타입 정의

**파일**: `websocket/errors.go`

```go
package websocket

// WebSocket 에러 타입 상수
const (
    WSErrWrongFormat        = "WRONG_FORMAT"
    WSErrNoTicket           = "NO_TICKET"
    WSErrNoType             = "NO_TYPE"
    WSErrNoCodes            = "NO_CODES"
    WSErrInvalidParam       = "INVALID_PARAM"
    WSErrInvalidParamFormat = "INVALID_PARAM_FORMAT"
)

// WSError represents WebSocket API error response
type WSError struct {
    Name    string `json:"name"`
    Message string `json:"message"`
}

// WSErrorResponse represents the error response structure
type WSErrorResponse struct {
    Error WSError `json:"error"`
}
```

### 3.2 인터페이스 확장

**파일**: `websocket/types.go`

```go
// MessageHandler 에러 핸들러 추가
type MessageHandler interface {
    Handle(data []byte) error
    Error(err WSError)  // 새로 추가
}
```

### 3.3 핸들러 변경

**파일**: `websocket/handler.go`

```go
func (c *Client) handleMessage(data []byte) error {
    var raw map[string]interface{}
    if err := json.Unmarshal(data, &raw); err != nil {
        return fmt.Errorf("parse message: %w", err)
    }

    // 에러 응답 체크 및 파싱
    if errResp, hasError := raw["error"]; hasError {
        // 에러 객체 파싱
        var wsErr WSError
        if errMap, ok := errResp.(map[string]interface{}); ok {
            if name, ok := errMap["name"].(string); ok {
                wsErr.Name = name
            }
            if msg, ok := errMap["message"].(string); ok {
                wsErr.Message = msg
            }
        }

        // 핸들러의 Error 콜백 호출
        c.mu.RLock()
        for _, h := range c.handlers {
            if h.Error != nil {
                h.Error(wsErr)
            }
        }
        c.mu.RUnlock()

        // 로깅 (타입별 레벨 구분)
        switch wsErr.Name {
        case WSErrInvalidParam:
            c.logger.Warn("WebSocket error", logger.F("type", wsErr.Name), logger.F("message", wsErr.Message))
        default:
            c.logger.Error("WebSocket error", logger.F("type", wsErr.Name), logger.F("message", wsErr.Message))
        }
        return nil
    }

    // ... 기존 메시지 처리 로직
}
```

---

## 4. Phase 2: PUBLIC API 누락 기능

### 4.1 추가할 API 목록

| API명 | 메서드 | 엔드포인트 | 파라미터 |
|-------|--------|-----------|----------|
| 경보제 | GET | `/v1/market/warning` | 없음 |
| 공지사항 | GET | `/v1/notice` | count (쿼리, optional) |
| 입출금수수료 | GET | `/v2/fee/inout/{currency}` | currency (Path, 필수) |

### 4.2 모델 정의

**파일**: `models/public/warning.go`

```go
package public

// WarningType represents warning alert types
type WarningType string

const (
    WarningPriceSuddenFluctuation           WarningType = "PRICE_SUDDEN_FLUCTUATION"
    WarningTradingVolumeSuddenFluctuation   WarningType = "TRADING_VOLUME_SUDDEN_FLUCTUATION"
    WarningDepositAmountSuddenFluctuation   WarningType = "DEPOSIT_AMOUNT_SUDDEN_FLUCTUATION"
    WarningPriceDifferenceHigh              WarningType = "PRICE_DIFFERENCE_HIGH"
    WarningSpecificAccountHighTransaction   WarningType = "SPECIFIC_ACCOUNT_HIGH_TRANSACTION"
    WarningExchangeTradingConcentration     WarningType = "EXCHANGE_TRADING_CONCENTRATION"
)

// Warning represents a market warning alert
type Warning struct {
    Market      string     `json:"market"`
    WarningType WarningType `json:"warning_type"`
    EndDate     string     `json:"end_date"`
}
```

**파일**: `models/public/notice.go`

```go
package public

// Notice represents a Bithumb announcement
type Notice struct {
    Categories  []string `json:"categories"`
    Title       string   `json:"title"`
    PCURL       string   `json:"pc_url"`
    PublishedAt string   `json:"published_at"`
    ModifiedAt  string   `json:"modified_at"`
}
```

**파일**: `models/public/fee.go`

```go
package public

// NetworkFee represents deposit/withdrawal fee information for a network
type NetworkFee struct {
    NetName                string `json:"net_name"`
    DepositFeeQuantity     string `json:"deposit_fee_quantity"`
    DepositMinimumQuantity string `json:"deposit_minimum_quantity"`
    WithdrawFeeQuantity    string `json:"withdraw_fee_quantity"`
    WithdrawMinimumQuantity string `json:"withdraw_minimum_quantity"`
}

// ChainFee represents fee information for a currency
type ChainFee struct {
    Name     string       `json:"name"`
    Currency string       `json:"currency"`
    Networks []NetworkFee `json:"networks"`
}
```

### 4.3 클라이언트 메서드

**파일**: `public/client.go`

```go
// GetWarnings retrieves market warning alerts
func (c *Client) GetWarnings() ([]Warning, *bithumbgo.Error) {
    return c.GetWarningsWithContext(context.Background())
}

func (c *Client) GetWarningsWithContext(ctx context.Context) ([]Warning, *bithumbgo.Error) {
    resp, err := c.do(ctx, http.MethodGet, c.base.BaseURL()+"/v1/market/warning", nil)
    if err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeNetwork, Err: err}
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return handleHTTPError(resp)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "read response failed", Err: err}
    }

    var warnings []Warning
    if err := json.Unmarshal(body, &warnings); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "parse response failed", Err: err}
    }

    return warnings, nil
}

// GetNotices retrieves Bithumb announcements
func (c *Client) GetNotices(req *GetNoticesRequest) ([]Notice, *bithumbgo.Error) {
    return c.GetNoticesWithContext(context.Background(), req)
}

func (c *Client) GetNoticesWithContext(ctx context.Context, req *GetNoticesRequest) ([]Notice, *bithumbgo.Error) {
    params := query.New().AddInt("count", req.Count)
    url := c.base.BaseURL() + "/v1/notice?" + params.Encode()

    resp, err := c.do(ctx, http.MethodGet, url, nil)
    // ... (유사한 에러 처리)
}

// GetChainFees retrieves deposit/withdrawal fees
func (c *Client) GetChainFees(currency string) ([]ChainFee, *bithumbgo.Error) {
    return c.GetChainFeesWithContext(context.Background(), currency)
}

func (c *Client) GetChainFeesWithContext(ctx context.Context, currency string) ([]ChainFee, *bithumbgo.Error) {
    url := c.base.BaseURL() + fmt.Sprintf("/v2/fee/inout/%s", currency)
    // ... (유사한 에러 처리)
}
```

---

## 5. Phase 3: PRIVATE API 핵심 기능

### 5.1 TWAP 주문 (3개)

| API명 | 메서드 | 엔드포인트 | 파라미터 |
|-------|--------|-----------|----------|
| TWAP 주문하기 | POST | `/v2/algo_orders/twap` | market, side, volume/price, duration, frequency |
| TWAP 주문 내역 조회 | GET | `/v1/algo_orders/twap` | market, uuids, state, next_key, limit, order_by |
| TWAP 주문 취소 | DELETE | `/v2/algo_orders/twap/{algo_order_id}` | algo_order_id |

### 5.2 입출금 기본 (8개)

| API명 | 메서드 | 엔드포인트 |
|-------|--------|-----------|
| 코인 입금 리스트 | GET | `/v1/coin/deposit` |
| 원화 입금 리스트 | GET | `/v1/krw/deposit` |
| 코인 출금 리스트 | GET | `/v1/coin/withdraw` |
| 원화 출금 리스트 | GET | `/v1/krw/withdraw` |
| 입출금 현황 | GET | `/v1/asset/status` |
| 전체 입금 주소 | GET | `/v1/addresses` |
| 개별 입금 조회 | GET | `/v1/coin/deposit/{uuid}` |
| 개별 출금 조회 | GET | `/v1/coin/withdraw/{uuid}` |

### 5.3 TWAP 모델

**파일**: `models/private/twap.go`

```go
package private

import "time"

// TWAPOrderSide represents TWAP order side
type TWAPOrderSide string

const (
    TWAPSideBid TWAPOrderSide = "bid"  // 매수
    TWAPSideAsk TWAPOrderSide = "ask"  // 매도
)

// TWAPState represents TWAP order state
type TWAPState string

const (
    TWAPStateProgress TWAPState = "progress"
    TWAPStateDone     TWAPState = "done"
    TWAPStateCancel   TWAPState = "cancel"
)

// TWAPOrder represents a TWAP algorithm order
type TWAPOrder struct {
    AlgoOrderID   string    `json:"algo_order_id"`
    Market        string    `json:"market"`
    Side          string    `json:"side"`
    Volume        string    `json:"volume"`
    Price         string    `json:"price"`
    Duration      int       `json:"duration"`
    Frequency     int       `json:"frequency"`
    State         TWAPState `json:"state"`
    RequestedTime time.Time `json:"requested_time"`
}

// PlaceTWAPOrderRequest represents request to place TWAP order
type PlaceTWAPOrderRequest struct {
    Market    string
    Side      TWAPOrderSide
    Volume    string  // 매도 시 필수
    Price     string  // 매수 시 필수
    Duration  int     // 300-43200초
    Frequency int     // 5, 15, 20, 30, 60, 120
}

func (r *PlaceTWAPOrderRequest) Validate() error {
    if r.Market == "" {
        return fmt.Errorf("market is required")
    }
    if r.Side == "" {
        return fmt.Errorf("side is required")
    }
    if r.Duration < 300 || r.Duration > 43200 {
        return fmt.Errorf("duration must be between 300 and 43200")
    }
    validFreq := map[int]bool{5: true, 15: true, 20: true, 30: true, 60: true, 120: true}
    if !validFreq[r.Frequency] {
        return fmt.Errorf("frequency must be one of: 5, 15, 20, 30, 60, 120")
    }
    if r.Side == TWAPSideAsk && r.Volume == "" {
        return fmt.Errorf("volume is required for ask order")
    }
    if r.Side == TWAPSideBid && r.Price == "" {
        return fmt.Errorf("price is required for bid order")
    }
    return nil
}

// GetTWAPOrdersRequest represents request to query TWAP orders
type GetTWAPOrdersRequest struct {
    Market   string
    UUIDs    []string
    State    TWAPState
    NextKey  string
    Limit    int
    OrderBy  string  // asc, desc
}

// CancelTWAPOrderRequest represents request to cancel TWAP order
type CancelTWAPOrderRequest struct {
    AlgoOrderID string
}
```

### 5.4 입출금 모델

**파일**: `models/private/deposit_withdraw.go`

```go
package private

// DepositState represents deposit state
type DepositState string

const (
    DepositStateRequestedPending      DepositState = "REQUESTED_PENDING"
    DepositStateRequestedSystemRejected DepositState = "REQUESTED_SYSTEM_REJECTED"
    DepositStateRequestedProcessing    DepositState = "REQUESTED_PROCESSING"
    DepositStateRequestedAdminRejected DepositState = "REQUESTED_ADMIN_REJECTED"
    DepositStateProcessing             DepositState = "DEPOSIT_PROCESSING"
    DepositStateAccepted               DepositState = "DEPOSIT_ACCEPTED"
    DepositStateCancelled              DepositState = "DEPOSIT_CANCELLED"
)

// WithdrawalState represents withdrawal state
type WithdrawalState string

const (
    WithdrawalStateProcessing WithdrawalState = "PROCESSING"
    WithdrawalStateDone       WithdrawalState = "DONE"
    WithdrawalStateCanceled   WithdrawalState = "CANCELED"
)

// WalletState represents wallet state
type WalletState string

const (
    WalletStateWorking     WalletState = "working"
    WalletStateWithdrawOnly WalletState = "withdraw_only"
    WalletStateDepositOnly  WalletState = "deposit_only"
    WalletStatePaused       WalletState = "paused"
)

// BlockState represents block state
type BlockState string

const (
    BlockStateNormal    BlockState = "normal"
    BlockStateDelayed   BlockState = "delayed"
    BlockStateInactive  BlockState = "inactive"
)

// CoinDeposit represents a coin deposit record
type CoinDeposit struct {
    UUID        string       `json:"uuid"`
    Currency    string       `json:"currency"`
    NetType     string       `json:"net_type"`
    Address     string       `json:"address"`
    Amount      string       `json:"amount"`
    State       DepositState `json:"state"`
    TXID        string       `json:"txid"`
    CreatedAt   string       `json:"created_at"`
    CompletedAt string       `json:"completed_at"`
}

// CoinWithdrawal represents a coin withdrawal record
type CoinWithdrawal struct {
    UUID         string           `json:"uuid"`
    Currency     string           `json:"currency"`
    NetType      string           `json:"net_type"`
    Address      string           `json:"address"`
    SecondaryAddress string       `json:"secondary_address"`
    Amount       string           `json:"amount"`
    Fee          string           `json:"fee"`
    State        WithdrawalState  `json:"state"`
    TXID         string           `json:"txid"`
    CreatedAt    string           `json:"created_at"`
    CompletedAt  string           `json:"completed_at"`
}

// AssetStatus represents deposit/withdrawal status for a currency
type AssetStatus struct {
    Currency              string     `json:"currency"`
    WalletState           WalletState `json:"wallet_state"`
    BlockState            BlockState  `json:"block_state"`
    BlockHeight           int         `json:"block_height"`
    BlockUpdatedAt        string      `json:"block_updated_at"`
    BlockElapsedMinutes   int         `json:"block_elapsed_minutes"`
    NetType               string      `json:"net_type"`
    NetworkName           string      `json:"network_name"`
}

// DepositAddress represents deposit address information
type DepositAddress struct {
    Currency          string `json:"currency"`
    NetType           string `json:"net_type"`
    DepositAddress    string `json:"deposit_address"`
    SecondaryAddress  string `json:"secondary_address"`
}
```

---

## 6. Phase 4: 코드 품질 개선

### 6.1 통합 에러 처리

모든 PUBLIC API 메서드가 `*bithumbgo.Error`를 반환하도록 변경합니다.

```go
// 변경 전
func (c *Client) GetTicker(req *GetTickerRequest) ([]Ticker, error)

// 변경 후
func (c *Client) GetTicker(req *GetTickerRequest) ([]Ticker, *bithumbgo.Error)
```

### 6.2 쿼리 파라미터 빌더

**파일**: `internal/query/builder.go`

```go
package query

import (
    "net/url"
    "strconv"
)

// Builder provides fluent API for building query strings
type Builder struct {
    values url.Values
}

// New creates a new query builder
func New() *Builder {
    return &Builder{values: make(url.Values)}
}

// Add adds a key-value pair if value is not empty
func (b *Builder) Add(key, value string) *Builder {
    if value != "" {
        b.values.Add(key, value)
    }
    return b
}

// AddInt adds an integer key-value pair if value > 0
func (b *Builder) AddInt(key string, value int) *Builder {
    if value > 0 {
        b.values.Add(key, strconv.Itoa(value))
    }
    return b
}

// AddStringSlice adds a string slice as multiple key[] parameters
func (b *Builder) AddStringSlice(key string, values []string) *Builder {
    for _, v := range values {
        if v != "" {
            b.values.Add(key+"[]", v)
        }
    }
    return b
}

// Encode returns the encoded query string
func (b *Builder) Encode() string {
    return b.values.Encode()
}
```

### 6.3 HTTP 에러 핸들러 공통화

**파일**: `public/client.go`

```go
// handleHTTPError creates a standardized error from HTTP response
func handleHTTPError(resp *http.Response) ([]byte, *bithumbgo.Error) {
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, &bithumbgo.Error{
            Type:       bithumbgo.ErrorTypeHTTP,
            Message:    fmt.Sprintf("API error: status %d (failed to read body)", resp.StatusCode),
            HTTPStatus: resp.StatusCode,
            Err:        err,
        }
    }
    return nil, &bithumbgo.Error{
        Type:       bithumbgo.ErrorTypeHTTP,
        Message:    fmt.Sprintf("API error: status %d: %s", resp.StatusCode, string(body)),
        HTTPStatus: resp.StatusCode,
    }
}
```

---

## 7. 파일 구조

```
bithumb-go/
├── websocket/
│   ├── errors.go          (신규) Phase 1
│   ├── client.go          (수정) Phase 1
│   ├── handler.go         (수정) Phase 1
│   └── types.go           (수정) Phase 1
│
├── public/
│   ├── client.go          (수정) Phase 2, 4
│   └── models/
│       ├── warning.go     (신규) Phase 2
│       ├── notice.go      (신규) Phase 2
│       └── fee.go         (신규) Phase 2
│
├── private/
│   ├── client_methods.go  (확장) Phase 3
│   └── models/
│       ├── twap.go        (신규) Phase 3
│       └── deposit_withdraw.go  (신규) Phase 3
│
└── internal/
    └── query/
        └── builder.go     (신규) Phase 4
```

---

## 8. 구현 순서

1. **Phase 1**: 웹소켓 에러 처리 (1-2일)
2. **Phase 2**: PUBLIC API 누락 (1일)
3. **Phase 4 (부분)**: 쿼리 빌더 및 에러 처리 (1일)
4. **Phase 3**: PRIVATE API 핵심 (2-3일)
5. **Phase 4 (완료)**: 나머지 코드 품질 개선 (1일)

**예상 기간**: 6-8일

---

## 9. 참고

- 빗썸 Open API 문서: `docs/api/`
- API 버전: v2.1.5
- 원본 문서: https://apidocs.bithumb.com/v2.1.5/reference
