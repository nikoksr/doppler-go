// Package logging provides a simple logging interface.
package logging

// Logger is the interface that wraps the basic logging methods. It is used by the SDK to log detailed information about
// requests and responses.
type Logger interface {
	// Debugw logs a message at debug level.
	Debugw(format string, keysAndValues ...any)

	// Infow logs a message at info level.
	Infow(format string, keysAndValues ...any)

	// Warnw logs a message at warn level.
	Warnw(format string, keysAndValues ...any)

	// Errorw logs a message at error level.
	Errorw(format string, keysAndValues ...any)
}

// NopLogger is a logger that does nothing. It is used by default.
type NopLogger struct{}

// Debugw logs a message at debug level.
func (l NopLogger) Debugw(_ string, _ ...any) {}

// Infow logs a message at info level.
func (l NopLogger) Infow(_ string, _ ...any) {}

// Warnw logs a message at warn level.
func (l NopLogger) Warnw(_ string, _ ...any) {}

// Errorw logs a message at error level.
func (l NopLogger) Errorw(_ string, _ ...any) {}
