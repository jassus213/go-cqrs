package go_cqrs

import "context"

// QueryValidation returns a decorator that calls Validate() on requests implementing
// [Validator]. Requests without validation pass through unchanged.
func QueryValidation[TRequest, TResponse any]() QueryDecorator[TRequest, TResponse] {
	return func(next QueryHandler[TRequest, TResponse]) QueryHandler[TRequest, TResponse] {
		return func(ctx context.Context, req TRequest) (TResponse, error) {
			if v, ok := any(req).(Validator); ok {
				if err := v.Validate(); err != nil {
					var zero TResponse
					return zero, err
				}
			}
			return next(ctx, req)
		}
	}
}

// CommandValidation returns a decorator that calls Validate() on commands implementing
// [Validator]. Commands without validation pass through unchanged.
func CommandValidation[TRequest any, TResponse any]() CommandDecorator[TRequest, TResponse] {
	return func(next CommandHandler[TRequest, TResponse]) CommandHandler[TRequest, TResponse] {
		return func(ctx context.Context, req TRequest) (TResponse, error) {
			if v, ok := any(req).(Validator); ok {
				if err := v.Validate(); err != nil {
					var zero TResponse
					return zero, err
				}
			}
			return next(ctx, req)
		}
	}
}
