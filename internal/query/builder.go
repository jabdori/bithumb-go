package query

import (
	"net/url"
	"strconv"
)

// Builder provides fluent API for building query strings
type Builder struct {
	values url.Values
}

// New creates a new query builder
func New() *Builder {
	return &Builder{values: make(url.Values)}
}

// Add adds a key-value pair if value is not empty
func (b *Builder) Add(key, value string) *Builder {
	if value != "" {
		b.values.Add(key, value)
	}
	return b
}

// AddInt adds an integer key-value pair if value != 0
func (b *Builder) AddInt(key string, value int) *Builder {
	if value != 0 {
		b.values.Add(key, strconv.Itoa(value))
	}
	return b
}

// AddStringSlice adds a string slice as multiple key[] parameters
func (b *Builder) AddStringSlice(key string, values []string) *Builder {
	if values == nil {
		return b
	}
	for _, v := range values {
		if v != "" {
			b.values.Add(key+"[]", v)
		}
	}
	return b
}

// Encode returns the encoded query string
func (b *Builder) Encode() string {
	return b.values.Encode()
}
