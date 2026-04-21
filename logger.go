package go_cqrs

import "context"

// Logger is a minimal interface for structured logging.
// Implement it on top of zerolog, slog, zap, or any other logger.
type Logger interface {
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}
