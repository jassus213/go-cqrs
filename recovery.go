package go_cqrs

import (
	"context"
	"fmt"
)

// QueryRecovery returns a decorator that catches panics from downstream handlers,
// logs them via [Logger.Error], and converts them into a regular error.
// It should typically be the outermost decorator in the pipeline.
func QueryRecovery[TRequest, TResponse any](logger Logger) QueryDecorator[TRequest, TResponse] {
	return func(next QueryHandler[TRequest, TResponse]) QueryHandler[TRequest, TResponse] {
		return func(ctx context.Context, req TRequest) (res TResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error(ctx, "recovered from panic", "panic", fmt.Sprintf("%v", r))
					err = fmt.Errorf("panic: %v", r)
				}
			}()
			return next(ctx, req)
		}
	}
}

// CommandRecovery returns a decorator that catches panics from downstream handlers,
// logs them via [Logger.Error], and converts them into a regular error.
// It should typically be the outermost decorator in the pipeline.
func CommandRecovery[TRequest any, TResponse any](logger Logger) CommandDecorator[TRequest, TResponse] {
	return func(next CommandHandler[TRequest, TResponse]) CommandHandler[TRequest, TResponse] {
		return func(ctx context.Context, req TRequest) (res TResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error(ctx, "recovered from panic", "panic", fmt.Sprintf("%v", r))
					err = fmt.Errorf("panic: %v", r)
				}
			}()
			return next(ctx, req)
		}
	}
}
