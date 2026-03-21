package go_cqrs

import (
	"context"
	"fmt"
	"time"
)

// Logging returns a decorator that logs the start and completion of each use-case
// execution. It records the operation name (derived from the request type) and elapsed time.
// Successful calls are logged at Debug level; errors are logged at Error level.
func Logging[TRequest, TResponse any](logger Logger) UseCaseDecorator[TRequest, TResponse] {
	return func(next UseCase[TRequest, TResponse]) UseCase[TRequest, TResponse] {
		return func(ctx context.Context, req TRequest) (TResponse, error) {
			start := time.Now()
			op := fmt.Sprintf("%T", req)

			logger.Debug(ctx, "executing", "operation", op)

			resp, err := next(ctx, req)

			elapsed := time.Since(start)
			if err != nil {
				logger.Error(ctx, "completed", "operation", op, "elapsed", elapsed.String(), "error", err.Error())
			} else {
				logger.Debug(ctx, "completed", "operation", op, "elapsed", elapsed.String())
			}

			return resp, err
		}
	}
}
