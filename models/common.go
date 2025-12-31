// Package models provides data types for Bithumb API responses.
package models

import "time"

// APIResponse represents a standard API response structure.
type APIResponse struct {
	Status       string      `json:"status"`
	Data         interface{} `json:"data"`
	ErrorMessage string      `json:"message,omitempty"`
}

// Timestamp represents a millisecond-precision Unix timestamp.
type Timestamp struct {
	Unix int64 `json:"timestamp"`
}

// Time converts the millisecond timestamp to time.Time.
func (t *Timestamp) Time() time.Time {
	return time.Unix(t.Unix/1000, (t.Unix%1000)*1000000)
}
