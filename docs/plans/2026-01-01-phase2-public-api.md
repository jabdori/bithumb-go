# Phase 2: PUBLIC API Missing Features Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement missing PUBLIC API endpoints: 경보제 (Warning), 공지사항 (Notice), 입출금수수료 (Chain Fee).

**Architecture:** Add model types for each API, implement client methods using unified error handling, write comprehensive tests.

**Tech Stack:** Go 1.23+, existing bithumb-go public client structure

---

## Prerequisites

- Complete Phase 1: WebSocket Error Handling
- Note: This phase also implements query builder (Phase 4 partial) for clean parameter handling

---

## Task 1: Create Query Builder Utility

**Files:**
- Create: `internal/query/builder.go`
- Test: `internal/query/builder_test.go`

**Step 1: Write the failing test**

Create `internal/query/builder_test.go`:

```go
package query

import (
    "net/url"
    "testing"
)

func TestBuilder_Add(t *testing.T) {
    b := New()
    b.Add("key", "value")

    result := b.Encode()
    expected := "key=value"

    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}

func TestBuilder_Add_EmptyValue(t *testing.T) {
    b := New()
    b.Add("key", "")

    result := b.Encode()

    if result != "" {
        t.Errorf("Expected empty string, got %s", result)
    }
}

func TestBuilder_AddInt(t *testing.T) {
    b := New()
    b.AddInt("count", 10)

    result := b.Encode()
    expected := "count=10"

    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}

func TestBuilder_AddInt_ZeroValue(t *testing.T) {
    b := New()
    b.AddInt("count", 0)

    result := b.Encode()

    if result != "" {
        t.Errorf("Expected empty string, got %s", result)
    }
}

func TestBuilder_AddStringSlice(t *testing.T) {
    b := New()
    b.AddStringSlice("uuids", []string{"uuid1", "uuid2"})

    result := b.Encode()

    parsed, _ := url.ParseQuery(result)
    if len(parsed["uuids[]"]) != 2 {
        t.Errorf("Expected 2 uuids, got %d", len(parsed["uuids[]"]))
    }
}

func TestBuilder_Multiple(t *testing.T) {
    b := New()
    b.Add("market", "KRW-BTC")
    b.AddInt("count", 5)

    result := b.Encode()

    parsed, _ := url.ParseQuery(result)
    if parsed.Get("market") != "KRW-BTC" {
        t.Errorf("Expected market=KRW-BTC")
    }
    if parsed.Get("count") != "5" {
        t.Errorf("Expected count=5")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/query -v`
Expected: FAIL with "no such package internal/query"

**Step 3: Write minimal implementation**

Create `internal/query/builder.go`:

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

**Step 4: Run test to verify it passes**

Run: `go test ./internal/query -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/query/
git commit -m "feat(query): add query builder utility"
```

---

## Task 2: Implement Warning (경보제) API

**Files:**
- Create: `models/public/warning.go`
- Modify: `public/client.go`
- Test: `public/client_test.go`

**Step 1: Write the failing test**

Add to `public/client_test.go`:

```go
func TestGetWarnings(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    server := mockServer(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/v1/market/warning" {
            t.Errorf("Expected path /v1/market/warning, got %s", r.URL.Path)
        }

        warnings := []Warning{
            {
                Market:      "KRW-BTC",
                WarningType: "PRICE_SUDDEN_FLUCTUATION",
                EndDate:     "2026-01-31 23:59:59",
            },
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(warnings)
    })
    defer server.Close()

    base := base.NewClient(server.URL, &http.Client{})
    client := NewClient(base)

    warnings, err := client.GetWarnings()
    if err != nil {
        t.Fatalf("GetWarnings failed: %v", err)
    }

    if len(warnings) != 1 {
        t.Fatalf("Expected 1 warning, got %d", len(warnings))
    }

    if warnings[0].Market != "KRW-BTC" {
        t.Errorf("Expected market KRW-BTC, got %s", warnings[0].Market)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./public -v -run TestGetWarnings`
Expected: FAIL with "undefined: Warning" or "undefined: client.GetWarnings"

**Step 3: Create model**

Create `models/public/warning.go`:

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

**Step 4: Add client method**

Add to `public/client.go`:

```go
// GetWarnings retrieves market warning alerts
func (c *Client) GetWarnings() ([]Warning, *bithumbgo.Error) {
    return c.GetWarningsWithContext(context.Background())
}

// GetWarningsWithContext retrieves market warning alerts with context
func (c *Client) GetWarningsWithContext(ctx context.Context) ([]Warning, *bithumbgo.Error) {
    url := c.base.BaseURL() + "/v1/market/warning"

    resp, err := c.do(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeNetwork, Message: "HTTP request failed", Err: err}
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
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "read response failed", Err: err}
    }

    var warnings []Warning
    if err := json.Unmarshal(body, &warnings); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "parse response failed", Err: err}
    }

    return warnings, nil
}
```

**Step 5: Run test to verify it passes**

Run: `go test ./public -v -run TestGetWarnings`
Expected: PASS

**Step 6: Commit**

```bash
git add models/public/warning.go public/client.go public/client_test.go
git commit -m "feat(public): add GetWarnings API for market alerts"
```

---

## Task 3: Implement Notice (공지사항) API

**Files:**
- Create: `models/public/notice.go`
- Modify: `public/client.go`
- Test: `public/client_test.go`

**Step 1: Write the failing test**

Add to `public/client_test.go`:

```go
func TestGetNotices(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    server := mockServer(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/v1/notice" {
            t.Errorf("Expected path /v1/notice, got %s", r.URL.Path)
        }

        notices := []Notice{
            {
                Categories:  []string{"공지"},
                Title:       "테스트 공지",
                PCURL:       "https://example.com/notice/1",
                PublishedAt: "2026-01-01 10:00:00",
                ModifiedAt:  "2026-01-01 10:00:00",
            },
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(notices)
    })
    defer server.Close()

    base := base.NewClient(server.URL, &http.Client{})
    client := NewClient(base)

    notices, err := client.GetNotices(&GetNoticesRequest{Count: 5})
    if err != nil {
        t.Fatalf("GetNotices failed: %v", err)
    }

    if len(notices) != 1 {
        t.Fatalf("Expected 1 notice, got %d", len(notices))
    }

    if notices[0].Title != "테스트 공지" {
        t.Errorf("Expected title '테스트 공지', got %s", notices[0].Title)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./public -v -run TestGetNotices`
Expected: FAIL with undefined types/methods

**Step 3: Create model and request type**

Create `models/public/notice.go`:

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

// GetNoticesRequest represents request parameters for GetNotices
type GetNoticesRequest struct {
    Count int // Number of notices to retrieve (1-20, default: 5)
}

// Validate validates the request parameters
func (r *GetNoticesRequest) Validate() error {
    if r.Count < 0 || r.Count > 20 {
        return fmt.Errorf("count must be between 0 and 20")
    }
    return nil
}
```

**Step 4: Add client method using query builder**

Add to `public/client.go`:

```go
// GetNotices retrieves Bithumb announcements
func (c *Client) GetNotices(req *GetNoticesRequest) ([]Notice, *bithumbgo.Error) {
    return c.GetNoticesWithContext(context.Background(), req)
}

// GetNoticesWithContext retrieves Bithumb announcements with context
func (c *Client) GetNoticesWithContext(ctx context.Context, req *GetNoticesRequest) ([]Notice, *bithumbgo.Error) {
    if err := req.Validate(); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeAPI, Message: "invalid request", Err: err}
    }

    params := query.New().AddInt("count", req.Count)
    url := c.base.BaseURL() + "/v1/notice?" + params.Encode()

    resp, err := c.do(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeNetwork, Message: "HTTP request failed", Err: err}
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
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "read response failed", Err: err}
    }

    var notices []Notice
    if err := json.Unmarshal(body, &notices); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "parse response failed", Err: err}
    }

    return notices, nil
}
```

Also add import: `"github.com/hysuki/bithumb-go/internal/query"`

**Step 5: Run test to verify it passes**

Run: `go test ./public -v -run TestGetNotices`
Expected: PASS

**Step 6: Commit**

```bash
git add models/public/notice.go public/client.go public/client_test.go
git commit -m "feat(public): add GetNotices API for announcements"
```

---

## Task 4: Implement Chain Fee (입출금수수료) API

**Files:**
- Create: `models/public/fee.go`
- Modify: `public/client.go`
- Test: `public/client_test.go`

**Step 1: Write the failing test**

Add to `public/client_test.go`:

```go
func TestGetChainFees(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    server := mockServer(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/v2/fee/inout/ALL" {
            t.Errorf("Expected path /v2/fee/inout/ALL, got %s", r.URL.Path)
        }

        fees := []ChainFee{
            {
                Name:     "비트코인",
                Currency: "BTC",
                Networks: []NetworkFee{
                    {
                        NetName:                "BTC",
                        DepositFeeQuantity:     "0",
                        DepositMinimumQuantity: "0.0004",
                        WithdrawFeeQuantity:    "0.0005",
                        WithdrawMinimumQuantity: "0.001",
                    },
                },
            },
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(fees)
    })
    defer server.Close()

    base := base.NewClient(server.URL, &http.Client{})
    client := NewClient(base)

    fees, err := client.GetChainFees("ALL")
    if err != nil {
        t.Fatalf("GetChainFees failed: %v", err)
    }

    if len(fees) != 1 {
        t.Fatalf("Expected 1 fee entry, got %d", len(fees))
    }

    if fees[0].Currency != "BTC" {
        t.Errorf("Expected currency BTC, got %s", fees[0].Currency)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./public -v -run TestGetChainFees`
Expected: FAIL with undefined types/methods

**Step 3: Create model**

Create `models/public/fee.go`:

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

**Step 4: Add client method**

Add to `public/client.go`:

```go
// GetChainFees retrieves deposit/withdrawal fees
func (c *Client) GetChainFees(currency string) ([]ChainFee, *bithumbgo.Error) {
    return c.GetChainFeesWithContext(context.Background(), currency)
}

// GetChainFeesWithContext retrieves deposit/withdrawal fees with context
func (c *Client) GetChainFeesWithContext(ctx context.Context, currency string) ([]ChainFee, *bithumbgo.Error) {
    if currency == "" {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeAPI, Message: "currency is required"}
    }

    url := c.base.BaseURL() + fmt.Sprintf("/v2/fee/inout/%s", currency)

    resp, err := c.do(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeNetwork, Message: "HTTP request failed", Err: err}
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
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "read response failed", Err: err}
    }

    var fees []ChainFee
    if err := json.Unmarshal(body, &fees); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "parse response failed", Err: err}
    }

    return fees, nil
}
```

**Step 5: Run test to verify it passes**

Run: `go test ./public -v -run TestGetChainFees`
Expected: PASS

**Step 6: Commit**

```bash
git add models/public/fee.go public/client.go public/client_test.go
git commit -m "feat(public): add GetChainFees API for deposit/withdrawal fees"
```

---

## Verification Steps

**Step 1: Run all PUBLIC tests**

```bash
go test ./public -v -cover
```

Expected: All tests PASS

**Step 2: Run query builder tests**

```bash
go test ./internal/query -v -cover
```

Expected: All tests PASS

**Step 3: Verify code compiles**

```bash
go build ./public
```

Expected: No errors

**Step 4: Final commit**

```bash
git add .
git commit -m "feat(public): complete Phase 2 - missing API implementations"
```

---

## Notes

- All new methods use `*bithumbgo.Error` for consistency (Phase 4 partial)
- Query builder eliminates manual string concatenation bugs
- Request validation uses `Validate()` method pattern
- HTTP status check added to all new methods
- Context-aware methods (`*WithContext`) for timeout/cancellation support

---

## References

- API Documentation: `docs/api/public/`
- Design Document: `docs/plans/2026-01-01-sdk-improvement-design.md`
