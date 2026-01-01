package query

import (
	"net/url"
	"testing"
)

func TestBuilder_Add(t *testing.T) {
	b := New()
	b.Add("key", "value")

	result := b.Encode()
	expected := "key=value"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestBuilder_Add_EmptyValue(t *testing.T) {
	b := New()
	b.Add("key", "")

	result := b.Encode()

	if result != "" {
		t.Errorf("Expected empty string, got %s", result)
	}
}

func TestBuilder_AddInt(t *testing.T) {
	b := New()
	b.AddInt("count", 10)

	result := b.Encode()
	expected := "count=10"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestBuilder_AddInt_ZeroValue(t *testing.T) {
	b := New()
	b.AddInt("count", 0)

	result := b.Encode()

	if result != "" {
		t.Errorf("Expected empty string, got %s", result)
	}
}

func TestBuilder_AddInt_NegativeValue(t *testing.T) {
	b := New()
	b.AddInt("offset", -10)

	result := b.Encode()
	expected := "offset=-10"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestBuilder_AddStringSlice(t *testing.T) {
	b := New()
	b.AddStringSlice("uuids", []string{"uuid1", "uuid2"})

	result := b.Encode()

	parsed, _ := url.ParseQuery(result)
	if len(parsed["uuids[]"]) != 2 {
		t.Errorf("Expected 2 uuids, got %d", len(parsed["uuids[]"]))
	}
}

func TestBuilder_AddStringSlice_Nil(t *testing.T) {
	b := New()
	b.AddStringSlice("uuids", nil)

	result := b.Encode()

	if result != "" {
		t.Errorf("Expected empty string for nil slice, got %s", result)
	}
}

func TestBuilder_Multiple(t *testing.T) {
	b := New()
	b.Add("market", "KRW-BTC")
	b.AddInt("count", 5)

	result := b.Encode()

	parsed, _ := url.ParseQuery(result)
	if parsed.Get("market") != "KRW-BTC" {
		t.Errorf("Expected market=KRW-BTC")
	}
	if parsed.Get("count") != "5" {
		t.Errorf("Expected count=5")
	}
}
