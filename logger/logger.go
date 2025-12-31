// Package logger provides a simple logger interface for the SDK.
package logger

// Logger is a simple logging interface.
// Users can provide their own logger implementation (e.g., zap, logrus, zerolog).
type Logger interface {
	// Debug logs a debug message.
	Debug(msg string, fields ...Field)
	// Info logs an info message.
	Info(msg string, fields ...Field)
	// Error logs an error message.
	Error(msg string, fields ...Field)
}

// Field is a key-value pair for structured logging.
type Field struct {
	Key   string
	Value interface{}
}

// F creates a new Field.
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// DefaultLogger is a no-op logger that discards all logs.
type DefaultLogger struct{}

// Debug logs a debug message (no-op).
func (DefaultLogger) Debug(msg string, fields ...Field) {}

// Info logs an info message (no-op).
func (DefaultLogger) Info(msg string, fields ...Field) {}

// Error logs an error message (no-op).
func (DefaultLogger) Error(msg string, fields ...Field) {}

// StdLogger wraps the standard log package.
type StdLogger struct{}

// Debug logs a debug message to stdout.
func (StdLogger) Debug(msg string, fields ...Field) {
	debugLog(msg, fields...)
}

// Info logs an info message to stdout.
func (StdLogger) Info(msg string, fields ...Field) {
	infoLog(msg, fields...)
}

// Error logs an error message to stderr.
func (StdLogger) Error(msg string, fields ...Field) {
	errorLog(msg, fields...)
}
