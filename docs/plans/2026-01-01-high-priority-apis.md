# High-Priority Bithumb APIs Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement high-priority missing Bithumb APIs to bring SDK coverage from 31.7% to approximately 60%

**Architecture:** Follow existing SDK patterns - models in `models/`, client methods in `public/` or `private/`, request/response types with validation, TDD approach with tests

**Tech Stack:** Go 1.21+, standard library, existing `client.Base` pattern

---

## Task 1: Market Code API (마켓 코드 조회)

**Files:**
- Create: `models/public/market.go`
- Modify: `public/client_methods.go` (append at end)
- Test: `public/client_methods_test.go`

**Step 1: Write the failing test**

```go
func TestGetMarketAll(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assertEquals(t, "/v1/market/all", r.URL.Path)
        assertEquals(t, "true", r.URL.Query().Get("isDetails"))
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`[{"market":"KRW-BTC","korean_name":"비트코인","english_name":"Bitcoin","market_warning":"NONE"}]`))
    }))
    defer server.Close()

    base := client.NewBase(server.URL, http.DefaultClient)
    c := public.NewClient(base)

    markets, err := c.GetMarketAll(true)
    assertNil(t, err)
    assertEqual(t, 1, len(markets))
    assertEqual(t, "KRW-BTC", markets[0].Market)
    assertEqual(t, "비트코인", markets[0].KoreanName)
    assertEqual(t, "Bitcoin", markets[0].EnglishName)
    assertEqual(t, "NONE", markets[0].MarketWarning)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./public -run TestGetMarketAll -v`
Expected: FAIL with "undefined: c.GetMarketAll"

**Step 3: Create Market type in models/public/market.go**

```go
package public

// Market represents a market code information.
type Market struct {
    // Market is the market identifier (e.g., "KRW-BTC").
    Market string `json:"market"`
    // KoreanName is the Korean name of the asset.
    KoreanName string `json:"korean_name"`
    // EnglishName is the English name of the asset.
    EnglishName string `json:"english_name"`
    // MarketWarning is the warning status ("NONE" or "CAUTION").
    MarketWarning string `json:"market_warning"`
}
```

**Step 4: Implement GetMarketAll method in public/client_methods.go**

```go
// GetMarketAll retrieves all available market codes.
func (c *Client) GetMarketAll(details bool) ([]Market, error) {
    return c.GetMarketAllWithContext(context.Background(), details)
}

// GetMarketAllWithContext retrieves all available market codes with context.
func (c *Client) GetMarketAllWithContext(ctx context.Context, details bool) ([]Market, error) {
    url := c.base.BaseURL() + "/v1/market/all"
    if details {
        url += "?isDetails=true"
    }

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, &bithumbgo.Error{
            Type:    bithumbgo.ErrorTypeNetwork,
            Message: "create request failed",
            Err:     err,
        }
    }

    resp, err := c.base.HTTPClient().Do(req)
    if err != nil {
        return nil, &bithumbgo.Error{
            Type:    bithumbgo.ErrorTypeNetwork,
            Message: "HTTP request failed",
            Err:     err,
        }
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        body, _ := io.ReadAll(resp.Body)
        return nil, &bithumbgo.Error{
            Type:       bithumbgo.ErrorTypeHTTP,
            Message:    fmt.Sprintf("API error: status %d: %s", resp.StatusCode, string(body)),
            HTTPStatus: resp.StatusCode,
        }
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, &bithumbgo.Error{
            Type:    bithumbgo.ErrorTypeParse,
            Message: "read response failed",
            Err:     err,
        }
    }

    var markets []Market
    if err := json.Unmarshal(body, &markets); err != nil {
        return nil, &bithumbgo.Error{
            Type:    bithumbgo.ErrorTypeParse,
            Message: "parse response failed",
            Err:     err,
        }
    }

    return markets, nil
}
```

**Step 5: Run test to verify it passes**

Run: `go test ./public -run TestGetMarketAll -v`
Expected: PASS

**Step 6: Commit**

```bash
git add models/public/market.go public/client_methods.go public/client_methods_test.go
git commit -m "feat: add GetMarketAll API for market code lookup"
```

---

## Task 2: Day Candle API (일 캔들)

**Files:**
- Modify: `models/public/candle.go` (add DayCandle type)
- Modify: `public/client_methods.go` (append method)
- Test: `public/client_methods_test.go`

**Step 1: Write the failing test**

```go
func TestGetDayCandles(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assertEquals(t, "/v1/candles/days", r.URL.Path)
        assertEquals(t, "KRW-BTC", r.URL.Query().Get("market"))
        assertEquals(t, "10", r.URL.Query().Get("count"))
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`[{"market":"KRW-BTC","candle_date_time_kst":"2024-01-01T00:00:00","opening_price":50000000,"high_price":51000000,"low_price":49000000,"trade_price":50500000,"timestamp":1704067200000,"candle_acc_trade_price":1000000000,"candle_acc_trade_volume":20.0,"prev_closing_price":50000000,"change_price":500000,"change_rate":0.01}]`))
    }))
    defer server.Close()

    base := client.NewBase(server.URL, http.DefaultClient)
    c := public.NewClient(base)

    req := &public.GetDayCandlesRequest{
        Market: "KRW-BTC",
        Count:  10,
    }
    candles, err := c.GetDayCandles(req)
    assertNil(t, err)
    assertEqual(t, 1, len(candles))
    assertEqual(t, 50500000.0, candles[0].TradePrice)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./public -run TestGetDayCandles -v`
Expected: FAIL with "undefined: c.GetDayCandles"

**Step 3: Add DayCandle type and request to models/public/candle.go**

```go
// DayCandle represents a day candle.
type DayCandle struct {
    // Market is the market identifier.
    Market string `json:"market"`
    // CandleDateTimeKST is the candle base time (KST).
    CandleDateTimeKST string `json:"candle_date_time_kst"`
    // CandleDateTimeUTC is the candle base time (UTC).
    CandleDateTimeUTC string `json:"candle_date_time_utc"`
    // OpeningPrice is the opening price.
    OpeningPrice float64 `json:"opening_price"`
    // HighPrice is the high price.
    HighPrice float64 `json:"high_price"`
    // LowPrice is the low price.
    LowPrice float64 `json:"low_price"`
    // TradePrice is the closing price.
    TradePrice float64 `json:"trade_price"`
    // Timestamp is the candle end time (KST).
    Timestamp int64 `json:"timestamp"`
    // CandleAccTradePrice is the accumulated trade price.
    CandleAccTradePrice float64 `json:"candle_acc_trade_price"`
    // CandleAccTradeVolume is the accumulated trade volume.
    CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`
    // PrevClosingPrice is the previous day closing price.
    PrevClosingPrice float64 `json:"prev_closing_price"`
    // ChangePrice is the change from previous close.
    ChangePrice float64 `json:"change_price"`
    // ChangeRate is the change rate from previous close.
    ChangeRate float64 `json:"change_rate"`
    // ConvertedTradePrice is the converted trade price (optional).
    ConvertedTradePrice float64 `json:"converted_trade_price,omitempty"`
}

// GetDayCandlesRequest is a request to get day candles.
type GetDayCandlesRequest struct {
    // Market is the market code (e.g., "KRW-BTC").
    Market string
    // To is the last candle time (exclusive), ISO8061 format.
    To string
    // Count is the number of candles (max 200, default 1).
    Count int
    // ConvertingPriceUnit is the price unit for conversion (e.g., "KRW").
    ConvertingPriceUnit string
}

// Validate checks if the request is valid.
func (r *GetDayCandlesRequest) Validate() error {
    if r.Market == "" {
        return fmt.Errorf("market is required")
    }
    if r.Count < 1 || r.Count > 200 {
        r.Count = 1
    }
    return nil
}
```

**Step 4: Implement GetDayCandles method in public/client_methods.go**

```go
// GetDayCandles retrieves day candles.
func (c *Client) GetDayCandles(req *GetDayCandlesRequest) ([]DayCandle, error) {
    return c.GetDayCandlesWithContext(context.Background(), req)
}

// GetDayCandlesWithContext retrieves day candles with context.
func (c *Client) GetDayCandlesWithContext(ctx context.Context, req *GetDayCandlesRequest) ([]DayCandle, error) {
    if err := req.Validate(); err != nil {
        return nil, &bithumbgo.Error{
            Type:    bithumbgo.ErrorTypeAPI,
            Message: fmt.Sprintf("invalid request: %v", err),
            Err:     err,
        }
    }

    url := c.base.BaseURL() + "/v1/candles/days"
    query := "market=" + req.Market
    if req.To != "" {
        query += "&to=" + req.To
    }
    if req.Count > 0 {
        query += fmt.Sprintf("&count=%d", req.Count)
    }
    if req.ConvertingPriceUnit != "" {
        query += "&convertingPriceUnit=" + req.ConvertingPriceUnit
    }
    url += "?" + query

    // ... HTTP request and response parsing similar to GetMarketAll ...
    // (Use same pattern as existing GetCandlestick method)

    var candles []DayCandle
    // ... unmarshal and return ...
}
```

**Step 5: Run test to verify it passes**

Run: `go test ./public -run TestGetDayCandles -v`
Expected: PASS

**Step 6: Commit**

```bash
git add models/public/candle.go public/client_methods.go public/client_methods_test.go
git commit -m "feat: add GetDayCandles API for day candles"
```

---

## Task 3: Week Candle API (주 캔들)

**Files:**
- Modify: `models/public/candle.go`
- Modify: `public/client_methods.go`
- Test: `public/client_methods_test.go`

**Step 1: Write the failing test**

```go
func TestGetWeekCandles(t *testing.T) {
    // Similar to TestGetDayCandles but with week candles
    // URL: /v1/candles/weeks
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./public -run TestGetWeekCandles -v`
Expected: FAIL

**Step 3: Add WeekCandle type to models/public/candle.go**

```go
// WeekCandle represents a week candle (same structure as DayCandle).
type WeekCandle struct {
    // Same fields as DayCandle
    Market string `json:"market"`
    CandleDateTimeKST string `json:"candle_date_time_kst"`
    CandleDateTimeUTC string `json:"candle_date_time_utc"`
    OpeningPrice float64 `json:"opening_price"`
    HighPrice float64 `json:"high_price"`
    LowPrice float64 `json:"low_price"`
    TradePrice float64 `json:"trade_price"`
    Timestamp int64 `json:"timestamp"`
    CandleAccTradePrice float64 `json:"candle_acc_trade_price"`
    CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`
    PrevClosingPrice float64 `json:"prev_closing_price"`
    ChangePrice float64 `json:"change_price"`
    ChangeRate float64 `json:"change_rate"`
    ConvertedTradePrice float64 `json:"converted_trade_price,omitempty"`
}

// GetWeekCandlesRequest is a request to get week candles.
type GetWeekCandlesRequest struct {
    Market string
    To string
    Count int
    ConvertingPriceUnit string
}

// Validate checks if the request is valid.
func (r *GetWeekCandlesRequest) Validate() error {
    if r.Market == "" {
        return fmt.Errorf("market is required")
    }
    if r.Count < 1 || r.Count > 200 {
        r.Count = 1
    }
    return nil
}
```

**Step 4: Implement GetWeekCandles method**

```go
// GetWeekCandles retrieves week candles.
func (c *Client) GetWeekCandles(req *GetWeekCandlesRequest) ([]WeekCandle, error) {
    return c.GetWeekCandlesWithContext(context.Background(), req)
}

// GetWeekCandlesWithContext retrieves week candles with context.
func (c *Client) GetWeekCandlesWithContext(ctx context.Context, req *GetWeekCandlesRequest) ([]WeekCandle, error) {
    // Similar to GetDayCandles but URL: /v1/candles/weeks
}
```

**Step 5: Run test to verify it passes**

Run: `go test ./public -run TestGetWeekCandles -v`
Expected: PASS

**Step 6: Commit**

```bash
git add models/public/candle.go public/client_methods.go public/client_methods_test.go
git commit -m "feat: add GetWeekCandles API for week candles"
```

---

## Task 4: Month Candle API (월 캔들)

**Files:**
- Modify: `models/public/candle.go`
- Modify: `public/client_methods.go`
- Test: `public/client_methods_test.go`

**Step 1-4:** Similar pattern to Day/Week candles
- URL: `/v1/candles/months`
- Type: `MonthCandle` (same structure)
- Request: `GetMonthCandlesRequest`

**Step 5: Run test**

Run: `go test ./public -run TestGetMonthCandles -v`
Expected: PASS

**Step 6: Commit**

```bash
git add models/public/candle.go public/client_methods.go public/client_methods_test.go
git commit -m "feat: add GetMonthCandles API for month candles"
```

---

## Task 5: Order Chance API (주문 가능 정보)

**Files:**
- Create: `models/private/order_chance.go`
- Modify: `private/client_methods.go`
- Test: `private/client_methods_test.go`

**Step 1: Write the failing test**

```go
func TestGetOrderChance(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assertEquals(t, "/v1/orders/chance", r.URL.Path)
        assertEquals(t, "KRW-BTC", r.URL.Query().Get("market"))
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"bid_fee":"0.0005","ask_fee":"0.0005","market":{"id":"KRW-BTC","name":"BTC","order_types":["limit"]}}`))
    }))
    defer server.Close()

    base := client.NewBaseWithKeys(server.URL, "key", "secret", http.DefaultClient)
    c := private.NewClient(base)

    chance, err := c.GetOrderChance(&private.GetOrderChanceRequest{Market: "KRW-BTC"})
    assertNil(t, err)
    assertEqual(t, "0.0005", chance.BidFee)
    assertEqual(t, "KRW-BTC", chance.Market.ID)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./private -run TestGetOrderChance -v`
Expected: FAIL with "undefined: c.GetOrderChance"

**Step 3: Create OrderChance type in models/private/order_chance.go**

```go
package private

// OrderChance represents order chance information.
type OrderChance struct {
    // BidFee is the bid fee rate.
    BidFee string `json:"bid_fee"`
    // AskFee is the ask fee rate.
    AskFee string `json:"ask_fee"`
    // MakerBidFee is the maker bid fee rate.
    MakerBidFee string `json:"maker_bid_fee"`
    // MakerAskFee is the maker ask fee rate.
    MakerAskFee string `json:"maker_ask_fee"`
    // Market is the market information.
    Market MarketInfo `json:"market"`
    // BidAccount is the bid account information.
    BidAccount AccountInfo `json:"bid_account"`
    // AskAccount is the ask account information.
    AskAccount AccountInfo `json:"ask_account"`
}

// MarketInfo represents market information.
type MarketInfo struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    OrderTypes  []string `json:"order_types"`
    AskTypes    []string `json:"ask_types"`
    BidTypes    []string `json:"bid_types"`
    OrderSides  []string `json:"order_sides"`
    Bid         Constraint `json:"bid"`
    Ask         Constraint `json:"ask"`
    MaxTotal    string    `json:"max_total"`
    State       string    `json:"state"`
}

// Constraint represents trading constraints.
type Constraint struct {
    Currency  string `json:"currency"`
    PriceUnit string `json:"price_unit"`
    MinTotal  string `json:"min_total"`
}

// AccountInfo represents account information.
type AccountInfo struct {
    Currency              string `json:"currency"`
    Balance               string `json:"balance"`
    Locked                string `json:"locked"`
    AvgBuyPrice           string `json:"avg_buy_price"`
    AvgBuyPriceModified   bool   `json:"avg_buy_price_modified"`
    UnitCurrency          string `json:"unit_currency"`
}

// GetOrderChanceRequest is a request to get order chance.
type GetOrderChanceRequest struct {
    // Market is the market code (e.g., "KRW-BTC").
    Market string
}

// Validate checks if the request is valid.
func (r *GetOrderChanceRequest) Validate() error {
    if r.Market == "" {
        return fmt.Errorf("market is required")
    }
    return nil
}
```

**Step 4: Implement GetOrderChance method in private/client_methods.go**

```go
// GetOrderChance retrieves order chance information.
func (c *Client) GetOrderChance(req *GetOrderChanceRequest) (*OrderChance, error) {
    return c.GetOrderChanceWithContext(context.Background(), req)
}

// GetOrderChanceWithContext retrieves order chance information with context.
func (c *Client) GetOrderChanceWithContext(ctx context.Context, req *GetOrderChanceRequest) (*OrderChance, error) {
    if err := req.Validate(); err != nil {
        return nil, &bithumbgo.Error{
            Type:    bithumbgo.ErrorTypeAPI,
            Message: fmt.Sprintf("invalid request: %v", err),
            Err:     err,
        }
    }

    url := fmt.Sprintf("%s/v1/orders/chance?market=%s", c.base.BaseURL(), req.Market)

    resp, err := c.doWithAuth(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, &bithumbgo.Error{
            Type:    bithumbgo.ErrorTypeParse,
            Message: "read response failed",
            Err:     err,
        }
    }

    var chance OrderChance
    if err := json.Unmarshal(body, &chance); err != nil {
        return nil, &bithumbgo.Error{
            Type:    bithumbgo.ErrorTypeParse,
            Message: "parse response failed",
            Err:     err,
        }
    }

    return &chance, nil
}
```

**Step 5: Run test to verify it passes**

Run: `go test ./private -run TestGetOrderChance -v`
Expected: PASS

**Step 6: Commit**

```bash
git add models/private/order_chance.go private/client_methods.go private/client_methods_test.go
git commit -m "feat: add GetOrderChance API for order availability check"
```

---

## Task 6: WebSocket Private Message Types

**Files:**
- Create: `models/websocket/my_order.go`
- Create: `models/websocket/my_asset.go`

**Step 1: Create MyOrderMessage type**

```go
package websocket

// MyOrderMessage represents a user's order WebSocket message.
type MyOrderMessage struct {
    Type string `json:"type"`
    // Order fields based on API documentation
    OrderID string `json:"uuid"`
    Market  string `json:"market"`
    Side    string `json:"side"`
    // ... add remaining fields from Order type
}
```

**Step 2: Create MyAssetMessage type**

```go
package websocket

// MyAssetMessage represents a user's asset WebSocket message.
type MyAssetMessage struct {
    Type string `json:"type"`
    // Asset fields based on API documentation
    // ... add fields from Account type
}
```

**Step 3: Run go vet**

Run: `go vet ./models/websocket/...`
Expected: PASS (no errors)

**Step 4: Commit**

```bash
git add models/websocket/my_order.go models/websocket/my_asset.go
git commit -m "feat: add WebSocket private message types (MyOrder, MyAsset)"
```

---

## Final Verification

**Step 1: Build all packages**

```bash
go build ./...
```

Expected: PASS (no errors)

**Step 2: Run all tests**

```bash
go test ./... -v
```

Expected: All tests pass

**Step 3: Run go vet**

```bash
go vet ./...
```

Expected: PASS (no warnings)

**Step 4: Final commit**

```bash
git commit --allow-empty -m "docs: complete high-priority API implementation plan"
```

---

## Summary

This plan implements:
1. **Market Code API** - Get all available markets
2. **Day/Week/Month Candles** - Extended time frame candles
3. **Order Chance API** - Check order availability before placing orders
4. **WebSocket Private Messages** - MyOrder and MyAsset message types

**Estimated new coverage: ~60%** (from 31.7%)
**New endpoints: 9 APIs** (Market, 3 candles, OrderChance, 2 WebSocket types)
**Files to create: 3** (market.go, order_chance.go, 2 websocket types)
**Files to modify: 3** (candle.go, client_methods.go x2, test files)

**Testing Strategy:** TDD for each new method, following existing patterns
**Commit Strategy:** One commit per API for easy rollback if needed
