package go_cqrs

import (
	"context"
	"fmt"
)

// Recovery returns a decorator that catches panics from downstream handlers,
// logs them via [Logger.Error], and converts them into a regular error.
// This should typically be the outermost decorator in the pipeline.
func Recovery[TRequest, TResponse any](logger Logger) UseCaseDecorator[TRequest, TResponse] {
	return func(next UseCase[TRequest, TResponse]) UseCase[TRequest, TResponse] {
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
