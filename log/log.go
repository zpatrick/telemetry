package log

import (
	"context"

	"github.com/zpatrick/telemetry/tag"
)

// The Logger interface is used to log messages.
type Logger interface {
	// Debug logs a message at the debug level.
	Debug(ctx context.Context, message string, tags ...tag.Tag)
	// Info logs a message at the info level.
	Info(ctx context.Context, message string, tags ...tag.Tag)
	// Warn logs a message at the warn level.
	Error(ctx context.Context, message string, tags ...tag.Tag)
	// With returns a new Logger which includes the given tags in each log.
	With(tags ...tag.Tag) Logger
}
