package go_cqrs

import (
	"context"
	"fmt"
	"time"
)

// QueryLogging returns a decorator that logs query execution.
// Successful calls are logged at Debug level; errors at Error level.
func QueryLogging[TRequest, TResponse any](logger Logger) QueryDecorator[TRequest, TResponse] {
	return func(next QueryHandler[TRequest, TResponse]) QueryHandler[TRequest, TResponse] {
		return func(ctx context.Context, req TRequest) (TResponse, error) {
			start := time.Now()
			op := fmt.Sprintf("%T", req)

			logger.Debug(ctx, "query started", "operation", op)

			resp, err := next(ctx, req)
			elapsed := time.Since(start)

			if err != nil {
				logger.Error(ctx, "query failed", "operation", op, "elapsed", elapsed.String(), "error", err.Error())
			} else {
				logger.Debug(ctx, "query completed", "operation", op, "elapsed", elapsed.String())
			}

			return resp, err
		}
	}
}

// CommandLogging returns a decorator that logs command execution.
// Successful calls are logged at Debug level; errors at Error level.
func CommandLogging[TRequest any, TResponse any](logger Logger) CommandDecorator[TRequest, TResponse] {
	return func(next CommandHandler[TRequest, TResponse]) CommandHandler[TRequest, TResponse] {
		return func(ctx context.Context, req TRequest) (TResponse, error) {
			start := time.Now()
			op := fmt.Sprintf("%T", req)

			logger.Debug(ctx, "command started", "operation", op)

			resp, err := next(ctx, req)
			elapsed := time.Since(start)

			if err != nil {
				logger.Error(ctx, "command failed", "operation", op, "elapsed", elapsed.String(), "error", err.Error())
			} else {
				logger.Debug(ctx, "command completed", "operation", op, "elapsed", elapsed.String())
			}

			return resp, err
		}
	}
}
