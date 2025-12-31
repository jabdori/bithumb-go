// Package bithumbgo provides a Go client for the Bithumb cryptocurrency exchange API.
package bithumbgo

import (
	"errors"
	"fmt"
)

// ErrorType represents the category of error that occurred.
type ErrorType int

const (
	ErrorTypeNetwork   ErrorType = iota
	ErrorTypeHTTP
	ErrorTypeAPI
	ErrorTypeParse
	ErrorTypeWebSocket
)

// String returns the string representation of ErrorType.
func (et ErrorType) String() string {
	switch et {
	case ErrorTypeNetwork:
		return "network error"
	case ErrorTypeHTTP:
		return "HTTP error"
	case ErrorTypeAPI:
		return "API error"
	case ErrorTypeParse:
		return "parse error"
	case ErrorTypeWebSocket:
		return "websocket error"
	default:
		return "unknown error"
	}
}

// Error represents an error from the Bithumb API.
// It supports error wrapping for determining the underlying cause.
type Error struct {
	Type       ErrorType
	Code       string
	Message    string
	HTTPStatus int
	Err        error
}

// Error returns the error message.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.Err
}

// IsAPIError returns true if err is an APIError with type ErrorTypeAPI.
func IsAPIError(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Type == ErrorTypeAPI
}

// IsRateLimitError returns true if err has HTTP status 429.
func IsRateLimitError(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.HTTPStatus == 429
}

// HasErrorCode returns true if err has the specified error code.
func HasErrorCode(err error, code string) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == code
}
