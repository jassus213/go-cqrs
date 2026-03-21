package go_cqrs

import "context"

// Logger is a minimal interface for structured logging.
// Implement this interface to adapt your logger (zerolog, zap, slog, etc.).
type Logger interface {
	Info(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}
