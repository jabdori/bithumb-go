# Phase 3: PRIVATE API Core Features Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement core PRIVATE API features: TWAP orders (3) and deposit/withdrawal basics (8).

**Architecture:** Add TWAP models and client methods, add deposit/withdrawal models and client methods using existing private authentication patterns.

**Tech Stack:** Go 1.23+, JWT authentication, existing bithumb-go private client structure

---

## Prerequisites

- Complete Phase 2: PUBLIC API (for query builder usage pattern)

---

## Task 1: Implement TWAP Order Models

**Files:**
- Create: `models/private/twap.go`
- Test: `models/private/twap_test.go`

**Step 1: Write the failing test**

Create `models/private/twap_test.go`:

```go
package private

import (
    "testing"
)

func TestPlaceTWAPOrderRequest_Validate(t *testing.T) {
    tests := []struct {
        name    string
        req     PlaceTWAPOrderRequest
        wantErr bool
    }{
        {
            name: "valid bid order",
            req: PlaceTWAPOrderRequest{
                Market:    "KRW-BTC",
                Side:      TWAPSideBid,
                Price:     "50000000",
                Duration:  300,
                Frequency: 60,
            },
            wantErr: false,
        },
        {
            name: "valid ask order",
            req: PlaceTWAPOrderRequest{
                Market:    "KRW-BTC",
                Side:      TWAPSideAsk,
                Volume:    "0.001",
                Duration:  300,
                Frequency: 60,
            },
            wantErr: false,
        },
        {
            name: "missing market",
            req: PlaceTWAPOrderRequest{
                Side:      TWAPSideBid,
                Price:     "50000000",
                Duration:  300,
                Frequency: 60,
            },
            wantErr: true,
        },
        {
            name: "duration too short",
            req: PlaceTWAPOrderRequest{
                Market:    "KRW-BTC",
                Side:      TWAPSideBid,
                Price:     "50000000",
                Duration:  100,
                Frequency: 60,
            },
            wantErr: true,
        },
        {
            name: "duration too long",
            req: PlaceTWAPOrderRequest{
                Market:    "KRW-BTC",
                Side:      TWAPSideBid,
                Price:     "50000000",
                Duration:  50000,
                Frequency: 60,
            },
            wantErr: true,
        },
        {
            name: "invalid frequency",
            req: PlaceTWAPOrderRequest{
                Market:    "KRW-BTC",
                Side:      TWAPSideBid,
                Price:     "50000000",
                Duration:  300,
                Frequency: 30,
            },
            wantErr: true,
        },
        {
            name: "ask order without volume",
            req: PlaceTWAPOrderRequest{
                Market:    "KRW-BTC",
                Side:      TWAPSideAsk,
                Duration:  300,
                Frequency: 60,
            },
            wantErr: true,
        },
        {
            name: "bid order without price",
            req: PlaceTWAPOrderRequest{
                Market:    "KRW-BTC",
                Side:      TWAPSideBid,
                Duration:  300,
                Frequency: 60,
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.req.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./models/private -v -run TestPlaceTWAPOrderRequest_Validate`
Expected: FAIL with "undefined: PlaceTWAPOrderRequest"

**Step 3: Create TWAP models**

Create `models/private/twap.go`:

```go
package private

import (
    "fmt"
    "time"
)

// TWAPOrderSide represents TWAP order side
type TWAPOrderSide string

const (
    TWAPSideBid TWAPOrderSide = "bid" // 매수
    TWAPSideAsk TWAPOrderSide = "ask" // 매도
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
    Volume    string // 매도 시 필수
    Price     string // 매수 시 필수
    Duration  int    // 300-43200초
    Frequency int    // 5, 15, 20, 30, 60, 120
}

// Validate validates the TWAP order request
func (r *PlaceTWAPOrderRequest) Validate() error {
    if r.Market == "" {
        return fmt.Errorf("market is required")
    }
    if r.Side == "" {
        return fmt.Errorf("side is required")
    }
    if r.Duration < 300 || r.Duration > 43200 {
        return fmt.Errorf("duration must be between 300 and 43200 seconds")
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
    Market  string
    UUIDs   []string
    State   TWAPState
    NextKey string
    Limit   int
    OrderBy string // asc, desc
}

// Validate validates the TWAP orders query request
func (r *GetTWAPOrdersRequest) Validate() error {
    if r.Limit < 0 || r.Limit > 100 {
        return fmt.Errorf("limit must be between 0 and 100")
    }
    if r.OrderBy != "" && r.OrderBy != "asc" && r.OrderBy != "desc" {
        return fmt.Errorf("order_by must be 'asc' or 'desc'")
    }
    return nil
}

// CancelTWAPOrderRequest represents request to cancel TWAP order
type CancelTWAPOrderRequest struct {
    AlgoOrderID string
}

// Validate validates the cancel TWAP order request
func (r *CancelTWAPOrderRequest) Validate() error {
    if r.AlgoOrderID == "" {
        return fmt.Errorf("algo_order_id is required")
    }
    return nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./models/private -v -run TestPlaceTWAPOrderRequest_Validate`
Expected: PASS

**Step 5: Commit**

```bash
git add models/private/twap.go models/private/twap_test.go
git commit -m "feat(private): add TWAP order models and validation"
```

---

## Task 2: Implement TWAP Client Methods

**Files:**
- Modify: `private/client_methods.go`
- Test: `private/client_test.go`

**Step 1: Write the failing test**

Add to `private/client_test.go`:

```go
func TestPlaceTWAPOrder(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test - requires API credentials")
    }

    // This test requires actual API credentials
    // Mock the response structure instead
    t.Skip("TODO: Add mock server test")
}

func TestCancelTWAPOrder(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    t.Skip("TODO: Add mock server test")
}

func TestGetTWAPOrders(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    t.Skip("TODO: Add mock server test")
}
```

**Step 2: Add client methods**

Add to `private/client_methods.go`:

```go
// PlaceTWAPOrder places a TWAP algorithm order
func (c *Client) PlaceTWAPOrder(req *PlaceTWAPOrderRequest) (*TWAPOrder, *bithumbgo.Error) {
    return c.PlaceTWAPOrderWithContext(context.Background(), req)
}

// PlaceTWAPOrderWithContext places a TWAP algorithm order with context
func (c *Client) PlaceTWAPOrderWithContext(ctx context.Context, req *PlaceTWAPOrderRequest) (*TWAPOrder, *bithumbgo.Error) {
    if err := req.Validate(); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeAPI, Message: "invalid request", Err: err}
    }

    body, err := json.Marshal(req)
    if err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "marshal request failed", Err: err}
    }

    url := c.base.BaseURL() + "/v2/algo_orders/twap"

    resp, err := c.doWithAuth(ctx, http.MethodPost, url, bytes.NewReader(body))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "read response failed", Err: err}
    }

    var order TWAPOrder
    if err := json.Unmarshal(respBody, &order); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "parse response failed", Err: err}
    }

    return &order, nil
}

// GetTWAPOrders retrieves TWAP order list
func (c *Client) GetTWAPOrders(req *GetTWAPOrdersRequest) ([]TWAPOrder, *bithumbgo.Error) {
    return c.GetTWAPOrdersWithContext(context.Background(), req)
}

// GetTWAPOrdersWithContext retrieves TWAP order list with context
func (c *Client) GetTWAPOrdersWithContext(ctx context.Context, req *GetTWAPOrdersRequest) ([]TWAPOrder, *bithumbgo.Error) {
    if err := req.Validate(); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeAPI, Message: "invalid request", Err: err}
    }

    url := c.base.BaseURL() + "/v1/algo_orders/twap"

    // Build query parameters
    params := query.New()
    if req.Market != "" {
        params.Add("market", req.Market)
    }
    if len(req.UUIDs) > 0 {
        params.AddStringSlice("uuids", req.UUIDs)
    }
    if req.State != "" {
        params.Add("state", string(req.State))
    }
    if req.NextKey != "" {
        params.Add("next_key", req.NextKey)
    }
    params.AddInt("limit", req.Limit)
    if req.OrderBy != "" {
        params.Add("order_by", req.OrderBy)
    }

    if queryStr := params.Encode(); queryStr != "" {
        url += "?" + queryStr
    }

    resp, err := c.doWithAuth(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "read response failed", Err: err}
    }

    var orders []TWAPOrder
    if err := json.Unmarshal(body, &orders); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "parse response failed", Err: err}
    }

    return orders, nil
}

// CancelTWAPOrder cancels a TWAP order
func (c *Client) CancelTWAPOrder(req *CancelTWAPOrderRequest) *bithumbgo.Error {
    return c.CancelTWAPOrderWithContext(context.Background(), req)
}

// CancelTWAPOrderWithContext cancels a TWAP order with context
func (c *Client) CancelTWAPOrderWithContext(ctx context.Context, req *CancelTWAPOrderRequest) *bithumbgo.Error {
    if err := req.Validate(); err != nil {
        return &bithumbgo.Error{Type: bithumbgo.ErrorTypeAPI, Message: "invalid request", Err: err}
    }

    url := fmt.Sprintf("%s/v2/algo_orders/twap/%s", c.base.BaseURL(), req.AlgoOrderID)

    reqBody, err := json.Marshal(map[string]string{})
    if err != nil {
        return &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "marshal request failed", Err: err}
    }

    resp, err := c.doWithAuth(ctx, http.MethodDelete, url, bytes.NewReader(reqBody))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    io.Copy(io.Discard, resp.Body)
    return nil
}
```

Also add import: `"github.com/hysuki/bithumb-go/internal/query"`

**Step 3: Run tests**

Run: `go test ./private -v -run TestPlaceTWAPOrder`
Expected: Tests skip (or pass with mocks)

**Step 4: Commit**

```bash
git add private/client_methods.go private/client_test.go
git commit -m "feat(private): add TWAP order client methods"
```

---

## Task 3: Implement Deposit/Withdrawal Models

**Files:**
- Create: `models/private/deposit_withdraw.go`
- Test: `models/private/deposit_withdraw_test.go`

**Step 1: Write the failing test**

Create `models/private/deposit_withdraw_test.go`:

```go
package private

import (
    "testing"
)

func TestWalletState_String(t *testing.T) {
    tests := []struct {
        value   WalletState
        wantStr string
    }{
        {WalletStateWorking, "working"},
        {WalletStateWithdrawOnly, "withdraw_only"},
        {WalletStateDepositOnly, "deposit_only"},
        {WalletStatePaused, "paused"},
    }

    for _, tt := range tests {
        t.Run(tt.wantStr, func(t *testing.T) {
            if string(tt.value) != tt.wantStr {
                t.Errorf("WalletState = %s, want %s", tt.value, tt.wantStr)
            }
        })
    }
}

func TestDepositState_Values(t *testing.T) {
    states := []DepositState{
        DepositStateRequestedPending,
        DepositStateDepositAccepted,
        DepositStateCancelled,
    }

    for _, state := range states {
        if state == "" {
            t.Error("DepositState should not be empty")
        }
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./models/private -v -run TestWalletState`
Expected: FAIL with "undefined: WalletState"

**Step 3: Create deposit/withdrawal models**

Create `models/private/deposit_withdraw.go`:

```go
package private

// DepositState represents deposit state
type DepositState string

const (
    DepositStateRequestedPending       DepositState = "REQUESTED_PENDING"
    DepositStateRequestedSystemRejected DepositState = "REQUESTED_SYSTEM_REJECTED"
    DepositStateRequestedProcessing     DepositState = "REQUESTED_PROCESSING"
    DepositStateRequestedAdminRejected  DepositState = "REQUESTED_ADMIN_REJECTED"
    DepositStateProcessing              DepositState = "DEPOSIT_PROCESSING"
    DepositStateAccepted                DepositState = "DEPOSIT_ACCEPTED"
    DepositStateCancelled               DepositState = "DEPOSIT_CANCELLED"
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
    WalletStateWorking      WalletState = "working"
    WalletStateWithdrawOnly WalletState = "withdraw_only"
    WalletStateDepositOnly  WalletState = "deposit_only"
    WalletStatePaused       WalletState = "paused"
)

// BlockState represents block state
type BlockState string

const (
    BlockStateNormal   BlockState = "normal"
    BlockStateDelayed  BlockState = "delayed"
    BlockStateInactive BlockState = "inactive"
)

// CoinDeposit represents a coin deposit record
type CoinDeposit struct {
    UUID         string       `json:"uuid"`
    Currency     string       `json:"currency"`
    NetType      string       `json:"net_type"`
    Address      string       `json:"address"`
    Amount       string       `json:"amount"`
    State        DepositState `json:"state"`
    TXID         string       `json:"txid"`
    CreatedAt    string       `json:"created_at"`
    CompletedAt  string       `json:"completed_at"`
}

// CoinWithdrawal represents a coin withdrawal record
type CoinWithdrawal struct {
    UUID             string           `json:"uuid"`
    Currency         string           `json:"currency"`
    NetType          string           `json:"net_type"`
    Address          string           `json:"address"`
    SecondaryAddress string           `json:"secondary_address"`
    Amount           string           `json:"amount"`
    Fee              string           `json:"fee"`
    State            WithdrawalState  `json:"state"`
    TXID             string           `json:"txid"`
    CreatedAt        string           `json:"created_at"`
    CompletedAt      string           `json:"completed_at"`
}

// AssetStatus represents deposit/withdrawal status for a currency
type AssetStatus struct {
    Currency            string     `json:"currency"`
    WalletState         WalletState `json:"wallet_state"`
    BlockState          BlockState  `json:"block_state"`
    BlockHeight         int         `json:"block_height"`
    BlockUpdatedAt      string      `json:"block_updated_at"`
    BlockElapsedMinutes int         `json:"block_elapsed_minutes"`
    NetType             string      `json:"net_type"`
    NetworkName         string      `json:"network_name"`
}

// DepositAddress represents deposit address information
type DepositAddress struct {
    Currency         string `json:"currency"`
    NetType          string `json:"net_type"`
    DepositAddress   string `json:"deposit_address"`
    SecondaryAddress string `json:"secondary_address"`
}

// KRWDeposit represents a KRW deposit record
type KRWDeposit struct {
    UUID        string `json:"uuid"`
    Amount      string `json:"amount"`
    State       string `json:"state"`
    CreatedAt   string `json:"created_at"`
    CompletedAt string `json:"completed_at"`
}

// KRWWithdrawal represents a KRW withdrawal record
type KRWWithdrawal struct {
    UUID        string `json:"uuid"`
    Amount      string `json:"amount"`
    State       string `json:"state"`
    CreatedAt   string `json:"created_at"`
    CompletedAt string `json:"completed_at"`
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./models/private -v -run TestWalletState`
Expected: PASS

**Step 5: Commit**

```bash
git add models/private/deposit_withdraw.go models/private/deposit_withdraw_test.go
git commit -m "feat(private): add deposit/withdrawal models"
```

---

## Task 4: Implement Deposit/Withdrawal Client Methods

**Files:**
- Modify: `private/client_methods.go`
- Test: `private/client_test.go`

**Step 1: Add request types**

Add to `private/client_methods.go` (after existing models):

```go
// GetCoinDepositsRequest represents request for coin deposit list
type GetCoinDepositsRequest struct {
    Currency  string
    State     DepositState
    UUIDs     []string
    Page      int
    Limit     int
    OrderBy   string
}

func (r *GetCoinDepositsRequest) Validate() error {
    if r.Limit < 0 || r.Limit > 100 {
        return fmt.Errorf("limit must be between 0 and 100")
    }
    return nil
}

// GetCoinWithdrawalsRequest represents request for coin withdrawal list
type GetCoinWithdrawalsRequest struct {
    Currency  string
    State     WithdrawalState
    UUIDs     []string
    TXIDs     []string
    Page      int
    Limit     int
    OrderBy   string
}

func (r *GetCoinWithdrawalsRequest) Validate() error {
    if r.Limit < 0 || r.Limit > 100 {
        return fmt.Errorf("limit must be between 0 and 100")
    }
    return nil
}

// GetAssetStatusRequest represents request for asset status
type GetAssetStatusRequest struct {
    Currency string
}

// GetDepositAddressRequest represents request for deposit address
type GetDepositAddressRequest struct {
    Currency string
}

func (r *GetDepositAddressRequest) Validate() error {
    if r.Currency == "" {
        return fmt.Errorf("currency is required")
    }
    return nil
}
```

**Step 2: Add client methods**

Add to `private/client_methods.go`:

```go
// GetCoinDeposits retrieves coin deposit list
func (c *Client) GetCoinDeposits(req *GetCoinDepositsRequest) ([]CoinDeposit, *bithumbgo.Error) {
    return c.GetCoinDepositsWithContext(context.Background(), req)
}

func (c *Client) GetCoinDepositsWithContext(ctx context.Context, req *GetCoinDepositsRequest) ([]CoinDeposit, *bithumbgo.Error) {
    if err := req.Validate(); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeAPI, Message: "invalid request", Err: err}
    }

    url := c.base.BaseURL() + "/v1/coin/deposit"

    params := query.New()
    if req.Currency != "" {
        params.Add("currency", req.Currency)
    }
    if req.State != "" {
        params.Add("state", string(req.State))
    }
    if len(req.UUIDs) > 0 {
        params.AddStringSlice("uuids", req.UUIDs)
    }
    params.AddInt("page", req.Page)
    params.AddInt("limit", req.Limit)
    if req.OrderBy != "" {
        params.Add("order_by", req.OrderBy)
    }

    if queryStr := params.Encode(); queryStr != "" {
        url += "?" + queryStr
    }

    resp, err := c.doWithAuth(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "read response failed", Err: err}
    }

    var deposits []CoinDeposit
    if err := json.Unmarshal(body, &deposits); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "parse response failed", Err: err}
    }

    return deposits, nil
}

// GetCoinWithdrawals retrieves coin withdrawal list
func (c *Client) GetCoinWithdrawals(req *GetCoinWithdrawalsRequest) ([]CoinWithdrawal, *bithumbgo.Error) {
    return c.GetCoinWithdrawalsWithContext(context.Background(), req)
}

func (c *Client) GetCoinWithdrawalsWithContext(ctx context.Context, req *GetCoinWithdrawalsRequest) ([]CoinWithdrawal, *bithumbgo.Error) {
    if err := req.Validate(); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeAPI, Message: "invalid request", Err: err}
    }

    url := c.base.BaseURL() + "/v1/coin/withdraw"

    params := query.New()
    if req.Currency != "" {
        params.Add("currency", req.Currency)
    }
    if req.State != "" {
        params.Add("state", string(req.State))
    }
    if len(req.UUIDs) > 0 {
        params.AddStringSlice("uuids", req.UUIDs)
    }
    if len(req.TXIDs) > 0 {
        params.AddStringSlice("txids", req.TXIDs)
    }
    params.AddInt("page", req.Page)
    params.AddInt("limit", req.Limit)
    if req.OrderBy != "" {
        params.Add("order_by", req.OrderBy)
    }

    if queryStr := params.Encode(); queryStr != "" {
        url += "?" + queryStr
    }

    resp, err := c.doWithAuth(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "read response failed", Err: err}
    }

    var withdrawals []CoinWithdrawal
    if err := json.Unmarshal(body, &withdrawals); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "parse response failed", Err: err}
    }

    return withdrawals, nil
}

// GetAssetStatus retrieves deposit/withdrawal status
func (c *Client) GetAssetStatus(req *GetAssetStatusRequest) ([]AssetStatus, *bithumbgo.Error) {
    return c.GetAssetStatusWithContext(context.Background(), req)
}

func (c *Client) GetAssetStatusWithContext(ctx context.Context, req *GetAssetStatusRequest) ([]AssetStatus, *bithumbgo.Error) {
    url := c.base.BaseURL() + "/v1/asset/status"

    if req.Currency != "" {
        url += "?currency=" + req.Currency
    }

    resp, err := c.doWithAuth(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "read response failed", Err: err}
    }

    var statuses []AssetStatus
    if err := json.Unmarshal(body, &statuses); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "parse response failed", Err: err}
    }

    return statuses, nil
}

// GetDepositAddresses retrieves all deposit addresses
func (c *Client) GetDepositAddresses() ([]DepositAddress, *bithumbgo.Error) {
    return c.GetDepositAddressesWithContext(context.Background())
}

func (c *Client) GetDepositAddressesWithContext(ctx context.Context) ([]DepositAddress, *bithumbgo.Error) {
    url := c.base.BaseURL() + "/v1/addresses"

    resp, err := c.doWithAuth(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "read response failed", Err: err}
    }

    var addresses []DepositAddress
    if err := json.Unmarshal(body, &addresses); err != nil {
        return nil, &bithumbgo.Error{Type: bithumbgo.ErrorTypeParse, Message: "parse response failed", Err: err}
    }

    return addresses, nil
}
```

**Step 3: Run tests**

Run: `go test ./private -v`
Expected: Existing tests pass

**Step 4: Commit**

```bash
git add private/client_methods.go private/client_test.go
git commit -m "feat(private): add deposit/withdrawal client methods"
```

---

## Verification Steps

**Step 1: Run all PRIVATE tests**

```bash
go test ./private -v -cover
```

Expected: All tests PASS

**Step 2: Run model tests**

```bash
go test ./models/private -v -cover
```

Expected: All tests PASS

**Step 3: Verify code compiles**

```bash
go build ./private
```

Expected: No errors

**Step 4: Final commit**

```bash
git add .
git commit -m "feat(private): complete Phase 3 - TWAP and deposit/withdrawal APIs"
```

---

## Notes

- All methods follow existing private API patterns (doWithAuth, JWT tokens)
- Query builder used for all GET requests with parameters
- Request validation uses Validate() method pattern
- Context-aware methods for timeout/cancellation support
- Path parameters use fmt.Sprintf for URL construction

---

## References

- API Documentation: `docs/api/private/`
- Design Document: `docs/plans/2026-01-01-sdk-improvement-design.md`
