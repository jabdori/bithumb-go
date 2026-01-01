# Phase 1: WebSocket Error Handling Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement WebSocket error handling according to Bithumb API documentation including error types, parsing, and handler callbacks.

**Architecture:** Add WSError types, extend MessageHandler interface with Error() callback, update handleMessage to parse and dispatch errors to handlers.

**Tech Stack:** Go 1.23+, github.com/coder/websocket, existing bithumb-go codebase

---

## Task 1: Create WebSocket Error Types

**Files:**
- Create: `websocket/errors.go`
- Test: `websocket/errors_test.go`

**Step 1: Write the failing test**

Create `websocket/errors_test.go`:

```go
package websocket

import (
    "encoding/json"
    "testing"
)

func TestWSError_Unmarshal(t *testing.T) {
    data := []byte(`{
        "error": {
            "name": "WRONG_FORMAT",
            "message": "Format 이 맞지 않습니다."
        }
    }`)

    var resp WSErrorResponse
    err := json.Unmarshal(data, &resp)
    if err != nil {
        t.Fatalf("Failed to unmarshal: %v", err)
    }

    if resp.Error.Name != WSErrWrongFormat {
        t.Errorf("Expected name %s, got %s", WSErrWrongFormat, resp.Error.Name)
    }

    if resp.Error.Message != "Format 이 맞지 않습니다." {
        t.Errorf("Expected message 'Format 이 맞지 않습니다.', got '%s'", resp.Error.Message)
    }
}

func TestWSErrorConstants(t *testing.T) {
    tests := []struct {
        name  string
        value string
    }{
        {"WrongFormat", "WRONG_FORMAT"},
        {"NoTicket", "NO_TICKET"},
        {"NoType", "NO_TYPE"},
        {"NoCodes", "NO_CODES"},
        {"InvalidParam", "INVALID_PARAM"},
        {"InvalidParamFormat", "INVALID_PARAM_FORMAT"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Constants should exist and match expected values
            switch tt.name {
            case "WrongFormat":
                if WSErrWrongFormat != tt.value {
                    t.Errorf("WSErrWrongFormat = %s, want %s", WSErrWrongFormat, tt.value)
                }
            case "NoTicket":
                if WSErrNoTicket != tt.value {
                    t.Errorf("WSErrNoTicket = %s, want %s", WSErrNoTicket, tt.value)
                }
            case "NoType":
                if WSErrNoType != tt.value {
                    t.Errorf("WSErrNoType = %s, want %s", WSErrNoType, tt.value)
                }
            case "NoCodes":
                if WSErrNoCodes != tt.value {
                    t.Errorf("WSErrNoCodes = %s, want %s", WSErrNoCodes, tt.value)
                }
            case "InvalidParam":
                if WSErrInvalidParam != tt.value {
                    t.Errorf("WSErrInvalidParam = %s, want %s", WSErrInvalidParam, tt.value)
                }
            case "InvalidParamFormat":
                if WSErrInvalidParamFormat != tt.value {
                    t.Errorf("WSErrInvalidParamFormat = %s, want %s", WSErrInvalidParamFormat, tt.value)
                }
            }
        })
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./websocket -v -run TestWSError`
Expected: FAIL with "undefined: WSErrWrongFormat" or similar

**Step 3: Write minimal implementation**

Create `websocket/errors.go`:

```go
package websocket

// WebSocket error type constants as defined in Bithumb API documentation
const (
    WSErrWrongFormat        = "WRONG_FORMAT"
    WSErrNoTicket           = "NO_TICKET"
    WSErrNoType             = "NO_TYPE"
    WSErrNoCodes            = "NO_CODES"
    WSErrInvalidParam       = "INVALID_PARAM"
    WSErrInvalidParamFormat = "INVALID_PARAM_FORMAT"
)

// WSError represents a WebSocket API error response
type WSError struct {
    Name    string `json:"name"`
    Message string `json:"message"`
}

// WSErrorResponse represents the full error response structure
type WSErrorResponse struct {
    Error WSError `json:"error"`
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./websocket -v -run TestWSError`
Expected: PASS

**Step 5: Commit**

```bash
git add websocket/errors.go websocket/errors_test.go
git commit -m "feat(websocket): add error types and constants"
```

---

## Task 2: Extend MessageHandler Interface with Error Callback

**Files:**
- Modify: `websocket/types.go`
- Modify: `websocket/subscription.go` (update existing handler implementations)

**Step 1: Write the failing test**

First, check current `types.go` to understand existing handler structure. Run:

```bash
cat websocket/types.go
```

**Step 2: Update types.go to add Error method to interface**

Add to `websocket/types.go`:

```go
// MessageHandler handles WebSocket messages for a subscription type
type MessageHandler interface {
    Handle(data []byte) error
    Error(err WSError)  // Error callback for WebSocket API errors
}
```

**Step 3: Update existing handler implementations**

Check `subscription.go` for existing handlers that need to implement Error:

```bash
grep -n "type.*struct" websocket/subscription.go
```

For each handler type, add an Error method. For example:

```go
func (h *TickerHandler) Error(err WSError) {
    // Default: no-op, users can embed this handler to override
}
```

**Step 4: Run tests to verify compilation**

Run: `go test ./websocket -v`
Expected: PASS (or failures if tests need updates)

**Step 5: Commit**

```bash
git add websocket/types.go websocket/subscription.go
git commit -m "feat(websocket): add Error callback to MessageHandler interface"
```

---

## Task 3: Update handleMessage to Parse and Dispatch Errors

**Files:**
- Modify: `websocket/handler.go`
- Test: `websocket/handler_test.go`

**Step 1: Write the failing test**

Create `websocket/handler_test.go`:

```go
package websocket

import (
    "encoding/json"
    "sync"
    "testing"
)

// mockHandler tracks whether Error was called
type mockHandler struct {
    mu          sync.Mutex
    errors      []WSError
    handleCalls int
}

func (m *mockHandler) Handle(data []byte) error {
    m.mu.Lock()
    m.handleCalls++
    m.mu.Unlock()
    return nil
}

func (m *mockHandler) Error(err WSError) {
    m.mu.Lock()
    m.errors = append(m.errors, err)
    m.mu.Unlock()
}

func TestHandleMessage_WithErrorResponse(t *testing.T) {
    client := NewClient(nil)

    errorData := []byte(`{
        "error": {
            "name": "NO_TICKET",
            "message": "티켓이 존재하지 않거나, 유효하지 않습니다."
        }
    }`)

    mock := &mockHandler{}
    client.handlers[SubscriptionTypeTicker] = mock

    err := client.handleMessage(errorData)
    if err != nil {
        t.Fatalf("handleMessage returned error: %v", err)
    }

    if len(mock.errors) != 1 {
        t.Fatalf("Expected 1 error, got %d", len(mock.errors))
    }

    if mock.errors[0].Name != WSErrNoTicket {
        t.Errorf("Expected error name %s, got %s", WSErrNoTicket, mock.errors[0].Name)
    }
}

func TestHandleMessage_WithNormalMessage(t *testing.T) {
    client := NewClient(nil)

    normalData := []byte(`{"type":"ticker","content":"test"}`)

    mock := &mockHandler{}
    client.handlers[SubscriptionTypeTicker] = mock

    err := client.handleMessage(normalData)
    if err != nil {
        t.Fatalf("handleMessage returned error: %v", err)
    }

    if len(mock.errors) != 0 {
        t.Fatalf("Expected 0 errors, got %d", len(mock.errors))
    }

    if mock.handleCalls != 1 {
        t.Errorf("Expected 1 handle call, got %d", mock.handleCalls)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./websocket -v -run TestHandleMessage`
Expected: FAIL (error response not parsed/handled correctly)

**Step 3: Update handleMessage implementation**

Modify `websocket/handler.go`:

```go
func (c *Client) handleMessage(data []byte) error {
    var raw map[string]interface{}
    if err := json.Unmarshal(data, &raw); err != nil {
        return fmt.Errorf("parse message: %w", err)
    }

    // Check for error response first
    if errResp, hasError := raw["error"]; hasError {
        var wsErr WSError

        // Try to parse as structured error response
        if errMap, ok := errResp.(map[string]interface{}); ok {
            if name, ok := errMap["name"].(string); ok {
                wsErr.Name = name
            }
            if msg, ok := errMap["message"].(string); ok {
                wsErr.Message = msg
            }
        }

        // Log error with appropriate level
        switch wsErr.Name {
        case WSErrInvalidParam:
            c.logger.Warn("WebSocket error", logger.F("type", wsErr.Name), logger.F("message", wsErr.Message))
        default:
            c.logger.Error("WebSocket error", logger.F("type", wsErr.Name), logger.F("message", wsErr.Message))
        }

        // Dispatch error to all registered handlers
        c.mu.RLock()
        for _, handler := range c.handlers {
            if handler.Error != nil {
                handler.Error(wsErr)
            }
        }
        c.mu.RUnlock()

        return nil  // Error responses are logged but don't stop processing
    }

    msgType, ok := raw["type"].(string)
    if !ok {
        return fmt.Errorf("missing type field")
    }

    c.mu.RLock()
    handler, exists := c.handlers[SubscriptionType(msgType)]
    c.mu.RUnlock()

    if !exists {
        return nil // No handler registered for this type
    }

    return handler.Handle(data)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./websocket -v -run TestHandleMessage`
Expected: PASS

**Step 5: Commit**

```bash
git add websocket/handler.go websocket/handler_test.go
git commit -m "feat(websocket): parse and dispatch WebSocket errors to handlers"
```

---

## Task 4: Add Comprehensive Error Tests

**Files:**
- Test: `websocket/errors_test.go`

**Step 1: Write tests for all error types**

Add to `websocket/errors_test.go`:

```go
func TestWSError_AllErrorTypes(t *testing.T) {
    errorCases := []struct {
        name          string
        errorName     string
        expectedConst string
    }{
        {"WrongFormat", "WRONG_FORMAT", WSErrWrongFormat},
        {"NoTicket", "NO_TICKET", WSErrNoTicket},
        {"NoType", "NO_TYPE", WSErrNoType},
        {"NoCodes", "NO_CODES", WSErrNoCodes},
        {"InvalidParam", "INVALID_PARAM", WSErrInvalidParam},
        {"InvalidParamFormat", "INVALID_PARAM_FORMAT", WSErrInvalidParamFormat},
    }

    for _, tc := range errorCases {
        t.Run(tc.name, func(t *testing.T) {
            data := []byte(`{
                "error": {
                    "name": "` + tc.errorName + `",
                    "message": "Test error message"
                }
            }`)

            var resp WSErrorResponse
            err := json.Unmarshal(data, &resp)
            if err != nil {
                t.Fatalf("Failed to unmarshal: %v", err)
            }

            if resp.Error.Name != tc.expectedConst {
                t.Errorf("Expected name %s, got %s", tc.expectedConst, resp.Error.Name)
            }
        })
    }
}
```

**Step 2: Run tests**

Run: `go test ./websocket -v -run TestWSError_AllErrorTypes`
Expected: PASS

**Step 3: Commit**

```bash
git add websocket/errors_test.go
git commit -m "test(websocket): add comprehensive error type tests"
```

---

## Verification Steps

After completing all tasks:

**Step 1: Run all WebSocket tests**

```bash
go test ./websocket -v -cover
```

Expected: All tests PASS with >80% coverage

**Step 2: Verify code compiles**

```bash
go build ./websocket
```

Expected: No errors

**Step 3: Check for race conditions**

```bash
go test ./websocket -race
```

Expected: No race warnings

**Step 4: Final commit**

```bash
git add .
git commit -m "feat(websocket): complete Phase 1 - error handling implementation"
```

---

## Notes

- Error parsing uses type assertion for flexibility with API response variations
- All registered handlers receive error callbacks for visibility
- Log levels differ by error type (INVALID_PARAM = Warn, others = Error)
- Handler.Error() is optional (no-op default in base implementations)

---

## References

- API Documentation: `docs/api/websocket/웹소켓-에러.md`
- Design Document: `docs/plans/2026-01-01-sdk-improvement-design.md`
