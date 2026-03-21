package go_cqrs

// UseCaseBuilder constructs a [UseCase] pipeline by chaining decorators.
// Decorators are applied in the order they are added: the first decorator
// wraps the outermost layer (executes first on the way in, last on the way out).
type UseCaseBuilder[TRequest any, TResponse any] struct {
	useCase    UseCase[TRequest, TResponse]
	decorators []UseCaseDecorator[TRequest, TResponse]
}

// NewBuilder creates a bare [UseCaseBuilder] with no decorators.
// Add decorators with [UseCaseBuilder.Use], then call [UseCaseBuilder.Build].
func NewBuilder[TRequest any, TResponse any](useCase UseCase[TRequest, TResponse]) *UseCaseBuilder[TRequest, TResponse] {
	return &UseCaseBuilder[TRequest, TResponse]{
		useCase: useCase,
	}
}

// NewDefaultBuilder creates a [UseCaseBuilder] pre-configured with a sensible
// decorator stack: [Recovery] → [Logging] → [Validation] → handler.
//
// You can append additional decorators with [UseCaseBuilder.Use] before calling Build.
func NewDefaultBuilder[Req, Res any](l Logger, handler UseCase[Req, Res]) *UseCaseBuilder[Req, Res] {
	return NewBuilder(handler).
		Use(
			Recovery[Req, Res](l),
			Logging[Req, Res](l),
			Validation[Req, Res](),
		)
}

// Use appends one or more decorators to the pipeline.
// Decorators run in the order they are added.
func (b *UseCaseBuilder[TRequest, TResponse]) Use(decorators ...UseCaseDecorator[TRequest, TResponse]) *UseCaseBuilder[TRequest, TResponse] {
	b.decorators = append(b.decorators, decorators...)
	return b
}

// Build compiles the decorator chain and returns a single [UseCase] function
// ready to be called. The builder should not be reused after calling Build.
func (b *UseCaseBuilder[TRequest, TResponse]) Build() UseCase[TRequest, TResponse] {
	h := b.useCase

	for i := len(b.decorators) - 1; i >= 0; i-- {
		h = b.decorators[i](h)
	}

	return h
}
