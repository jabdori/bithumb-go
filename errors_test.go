package bithumbgo

import (
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	err := &Error{Type: ErrorTypeAPI, Message: "test error"}
	msg := err.Error()
	if msg == "" {
		t.Error("Error() should return message")
	}
}

func TestError_Unwrap(t *testing.T) {
	original := errors.New("original")
	err := &Error{Type: ErrorTypeNetwork, Err: original}
	if !errors.Is(err.Unwrap(), original) {
		t.Error("Unwrap() should return original error")
	}
}

func TestIsAPIError(t *testing.T) {
	apiErr := &Error{Type: ErrorTypeAPI}
	if !IsAPIError(apiErr) {
		t.Error("Should be API error")
	}
	networkErr := &Error{Type: ErrorTypeNetwork}
	if IsAPIError(networkErr) {
		t.Error("Network error should not be API error")
	}
}

func TestIsRateLimitError(t *testing.T) {
	rateLimitErr := &Error{HTTPStatus: 429}
	if !IsRateLimitError(rateLimitErr) {
		t.Error("Should be rate limit error")
	}
	normalErr := &Error{HTTPStatus: 200}
	if IsRateLimitError(normalErr) {
		t.Error("Normal status should not be rate limit")
	}
}

func TestHasErrorCode(t *testing.T) {
	err := &Error{Code: "INVALID_PARAMS"}
	if !HasErrorCode(err, "INVALID_PARAMS") {
		t.Error("Should have matching code")
	}
	if HasErrorCode(err, "OTHER_CODE") {
		t.Error("Should not have different code")
	}
}

func TestErrorType_String(t *testing.T) {
	tests := []struct {
		errType ErrorType
		want    string
	}{
		{ErrorTypeNetwork, "network error"},
		{ErrorTypeHTTP, "HTTP error"},
		{ErrorTypeAPI, "API error"},
		{ErrorTypeParse, "parse error"},
		{ErrorTypeWebSocket, "websocket error"},
	}

	for _, tt := range tests {
		if got := tt.errType.String(); got != tt.want {
			t.Errorf("ErrorType.String() = %v, want %v", got, tt.want)
		}
	}
}
