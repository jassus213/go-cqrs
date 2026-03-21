package go_cqrs

import (
	"context"
)

// Validation returns a decorator that checks whether the request implements
// the [Validator] interface. If it does, Validate() is called before the
// request reaches the handler. A non-nil error short-circuits the pipeline.
//
// Requests that do not implement [Validator] pass through unchanged.
func Validation[TRequest, TResponse any]() UseCaseDecorator[TRequest, TResponse] {
	return func(next UseCase[TRequest, TResponse]) UseCase[TRequest, TResponse] {
		return func(ctx context.Context, req TRequest) (TResponse, error) {
			if v, ok := any(req).(Validator); ok {
				if err := v.Validate(); err != nil {
					return *new(TResponse), err
				}
			}

			return next(ctx, req)
		}
	}
}
