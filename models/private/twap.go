package private

import (
	"fmt"
	"strings"
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
	AlgoOrderID   string       `json:"algo_order_id"`
	Market        string       `json:"market"`
	Side          TWAPOrderSide `json:"side"`
	Volume        string       `json:"volume"`
	Price         string       `json:"price"`
	Duration      int          `json:"duration"`
	Frequency     int          `json:"frequency"`
	State         TWAPState    `json:"state"`
	RequestedTime time.Time    `json:"requested_time"`
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
	if r.Side != TWAPSideBid && r.Side != TWAPSideAsk {
		return fmt.Errorf("side must be '%s' or '%s'", TWAPSideBid, TWAPSideAsk)
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
	if r.State != "" && r.State != TWAPStateProgress && r.State != TWAPStateDone && r.State != TWAPStateCancel {
		return fmt.Errorf("state must be '%s', '%s', or '%s'", TWAPStateProgress, TWAPStateDone, TWAPStateCancel)
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
	if r.AlgoOrderID == "" || strings.TrimSpace(r.AlgoOrderID) == "" {
		return fmt.Errorf("algo_order_id is required")
	}
	return nil
}
