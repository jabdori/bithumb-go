package models

import (
	"testing"
	"time"
)

func TestTimestamp_Time(t *testing.T) {
	// 2024-01-01 00:00:00 UTC in milliseconds
	ts := &Timestamp{Unix: 1704067200000}
	tm := ts.Time()

	expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !tm.Equal(expected) {
		t.Errorf("Time() = %v, want %v", tm, expected)
	}
}

func TestTimestamp_Time_Zero(t *testing.T) {
	ts := &Timestamp{Unix: 0}
	tm := ts.Time()

	// Unix timestamp 0 represents the epoch: 1970-01-01 00:00:00 UTC
	expected := time.Unix(0, 0)
	if !tm.Equal(expected) {
		t.Errorf("Time() = %v, want %v", tm, expected)
	}
}
