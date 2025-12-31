package bithumbgo

import (
	"errors"
	"fmt"
)

type ErrorType int

const (
	ErrorTypeNetwork ErrorType = iota
	ErrorTypeHTTP
	ErrorTypeAPI
	ErrorTypeParse
	ErrorTypeWebSocket
)

type Error struct {
	Type       ErrorType
	Code       string
	Message    string
	HTTPStatus int
	Err        error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func IsAPIError(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Type == ErrorTypeAPI
}

func IsRateLimitError(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.HTTPStatus == 429
}

func HasErrorCode(err error, code string) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == code
}
