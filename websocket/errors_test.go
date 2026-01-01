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
