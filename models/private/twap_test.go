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
			name: "invalid side enum value",
			req: PlaceTWAPOrderRequest{
				Market:    "KRW-BTC",
				Side:      "invalid",
				Price:     "50000000",
				Duration:  300,
				Frequency: 60,
			},
			wantErr: true,
		},
		{
			name: "empty side",
			req: PlaceTWAPOrderRequest{
				Market:    "KRW-BTC",
				Side:      "",
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
			name: "duration boundary - minimum valid (300)",
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
			name: "duration boundary - maximum valid (43200)",
			req: PlaceTWAPOrderRequest{
				Market:    "KRW-BTC",
				Side:      TWAPSideBid,
				Price:     "50000000",
				Duration:  43200,
				Frequency: 60,
			},
			wantErr: false,
		},
		{
			name: "duration boundary - just below minimum (299)",
			req: PlaceTWAPOrderRequest{
				Market:    "KRW-BTC",
				Side:      TWAPSideBid,
				Price:     "50000000",
				Duration:  299,
				Frequency: 60,
			},
			wantErr: true,
		},
		{
			name: "duration boundary - just above maximum (43201)",
			req: PlaceTWAPOrderRequest{
				Market:    "KRW-BTC",
				Side:      TWAPSideBid,
				Price:     "50000000",
				Duration:  43201,
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
				Frequency: 25,
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

func TestGetTWAPOrdersRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     GetTWAPOrdersRequest
		wantErr bool
	}{
		{
			name: "valid request with no filters",
			req: GetTWAPOrdersRequest{
				Market: "KRW-BTC",
			},
			wantErr: false,
		},
		{
			name: "valid request with state filter",
			req: GetTWAPOrdersRequest{
				Market: "KRW-BTC",
				State:  TWAPStateProgress,
			},
			wantErr: false,
		},
		{
			name: "valid request with done state",
			req: GetTWAPOrdersRequest{
				Market: "KRW-BTC",
				State:  TWAPStateDone,
			},
			wantErr: false,
		},
		{
			name: "valid request with cancel state",
			req: GetTWAPOrdersRequest{
				Market: "KRW-BTC",
				State:  TWAPStateCancel,
			},
			wantErr: false,
		},
		{
			name: "invalid state enum value",
			req: GetTWAPOrdersRequest{
				Market: "KRW-BTC",
				State:  "invalid",
			},
			wantErr: true,
		},
		{
			name: "valid with limit",
			req: GetTWAPOrdersRequest{
				Market: "KRW-BTC",
				Limit:  50,
			},
			wantErr: false,
		},
		{
			name: "limit boundary - minimum (0)",
			req: GetTWAPOrdersRequest{
				Market: "KRW-BTC",
				Limit:  0,
			},
			wantErr: false,
		},
		{
			name: "limit boundary - maximum (100)",
			req: GetTWAPOrdersRequest{
				Market: "KRW-BTC",
				Limit:  100,
			},
			wantErr: false,
		},
		{
			name: "limit too low (-1)",
			req: GetTWAPOrdersRequest{
				Market: "KRW-BTC",
				Limit:  -1,
			},
			wantErr: true,
		},
		{
			name: "limit too high (101)",
			req: GetTWAPOrdersRequest{
				Market: "KRW-BTC",
				Limit:  101,
			},
			wantErr: true,
		},
		{
			name: "valid order by asc",
			req: GetTWAPOrdersRequest{
				Market:  "KRW-BTC",
				OrderBy: "asc",
			},
			wantErr: false,
		},
		{
			name: "valid order by desc",
			req: GetTWAPOrdersRequest{
				Market:  "KRW-BTC",
				OrderBy: "desc",
			},
			wantErr: false,
		},
		{
			name: "invalid order by",
			req: GetTWAPOrdersRequest{
				Market:  "KRW-BTC",
				OrderBy: "invalid",
			},
			wantErr: true,
		},
		{
			name: "empty order by is valid",
			req: GetTWAPOrdersRequest{
				Market:  "KRW-BTC",
				OrderBy: "",
			},
			wantErr: false,
		},
		{
			name: "valid request with all fields",
			req: GetTWAPOrdersRequest{
				Market:  "KRW-BTC",
				UUIDs:   []string{"uuid1", "uuid2"},
				State:   TWAPStateProgress,
				NextKey: "next",
				Limit:   50,
				OrderBy: "asc",
			},
			wantErr: false,
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

func TestCancelTWAPOrderRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CancelTWAPOrderRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CancelTWAPOrderRequest{
				AlgoOrderID: "algo-order-123",
			},
			wantErr: false,
		},
		{
			name: "missing algo order id",
			req: CancelTWAPOrderRequest{
				AlgoOrderID: "",
			},
			wantErr: true,
		},
		{
			name: "whitespace only algo order id",
			req: CancelTWAPOrderRequest{
				AlgoOrderID: "   ",
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
