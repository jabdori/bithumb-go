package public

import (
	"testing"
)

func TestGetNoticesRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		wantErr bool
		errMsg  string
	}{
		// Edge cases - invalid values
		{
			name:    "count = 0 should error",
			count:   0,
			wantErr: true,
			errMsg:  "count must be between 1 and 20, got 0",
		},
		{
			name:    "count = -1 should error",
			count:   -1,
			wantErr: true,
			errMsg:  "count must be between 1 and 20, got -1",
		},
		{
			name:    "count = 21 should error",
			count:   21,
			wantErr: true,
			errMsg:  "count must be between 1 and 20, got 21",
		},
		{
			name:    "count = 100 should error",
			count:   100,
			wantErr: true,
			errMsg:  "count must be between 1 and 20, got 100",
		},
		// Boundary cases - valid values
		{
			name:    "count = 1 (lower boundary) should work",
			count:   1,
			wantErr: false,
		},
		{
			name:    "count = 20 (upper boundary) should work",
			count:   20,
			wantErr: false,
		},
		// Happy path
		{
			name:    "count = 5 (default) should work",
			count:   5,
			wantErr: false,
		},
		{
			name:    "count = 10 should work",
			count:   10,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GetNoticesRequest{
				Count: tt.count,
			}

			err := r.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}
