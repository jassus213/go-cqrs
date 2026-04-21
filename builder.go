package go_cqrs

// QueryBuilder constructs a [QueryHandler] pipeline by chaining decorators.
type QueryBuilder[TRequest any, TResponse any] struct {
	handler    QueryHandler[TRequest, TResponse]
	decorators []QueryDecorator[TRequest, TResponse]
}

// NewQueryBuilder creates a bare [QueryBuilder] with no decorators.
// Add decorators with [QueryBuilder.Use], then call [QueryBuilder.Build].
func NewQueryBuilder[TRequest any, TResponse any](handler QueryHandler[TRequest, TResponse]) *QueryBuilder[TRequest, TResponse] {
	return &QueryBuilder[TRequest, TResponse]{handler: handler}
}

// NewDefaultQueryBuilder creates a [QueryBuilder] pre-configured with a sensible
// decorator stack: [QueryRecovery] → [QueryLogging] → [QueryValidation] → handler.
//
// You can append additional decorators with [QueryBuilder.Use] before calling Build.
func NewDefaultQueryBuilder[TRequest any, TResponse any](l Logger, handler QueryHandler[TRequest, TResponse]) *QueryBuilder[TRequest, TResponse] {
	return NewQueryBuilder(handler).
		Use(
			QueryRecovery[TRequest, TResponse](l),
			QueryLogging[TRequest, TResponse](l),
			QueryValidation[TRequest, TResponse](),
		)
}

// Use appends one or more decorators to the pipeline.
// Decorators run in the order they are added.
func (b *QueryBuilder[TRequest, TResponse]) Use(decorators ...QueryDecorator[TRequest, TResponse]) *QueryBuilder[TRequest, TResponse] {
	b.decorators = append(b.decorators, decorators...)
	return b
}

// Build compiles the decorator chain and returns a ready [QueryHandler].
// The builder should not be reused after calling Build.
func (b *QueryBuilder[TRequest, TResponse]) Build() QueryHandler[TRequest, TResponse] {
	h := b.handler
	for i := len(b.decorators) - 1; i >= 0; i-- {
		h = b.decorators[i](h)
	}
	return h
}

// CommandBuilder constructs a [CommandHandler] pipeline by chaining decorators.
type CommandBuilder[TRequest any, TResponse any] struct {
	handler    CommandHandler[TRequest, TResponse]
	decorators []CommandDecorator[TRequest, TResponse]
}

// NewCommandBuilder creates a bare [CommandBuilder] with no decorators.
// Add decorators with [CommandBuilder.Use], then call [CommandBuilder.Build].
func NewCommandBuilder[TRequest any, TResponse any](handler CommandHandler[TRequest, TResponse]) *CommandBuilder[TRequest, TResponse] {
	return &CommandBuilder[TRequest, TResponse]{handler: handler}
}

// NewDefaultCommandBuilder creates a [CommandBuilder] pre-configured with a sensible
// decorator stack: [CommandRecovery] → [CommandLogging] → [CommandValidation] → handler.
//
// You can append additional decorators with [CommandBuilder.Use] before calling Build.
func NewDefaultCommandBuilder[TRequest any, TResponse any](l Logger, handler CommandHandler[TRequest, TResponse]) *CommandBuilder[TRequest, TResponse] {
	return NewCommandBuilder(handler).
		Use(
			CommandRecovery[TRequest, TResponse](l),
			CommandLogging[TRequest, TResponse](l),
			CommandValidation[TRequest, TResponse](),
		)
}

// Use appends one or more decorators to the pipeline.
// Decorators run in the order they are added.
func (b *CommandBuilder[TRequest, TResponse]) Use(decorators ...CommandDecorator[TRequest, TResponse]) *CommandBuilder[TRequest, TResponse] {
	b.decorators = append(b.decorators, decorators...)
	return b
}

// Build compiles the decorator chain and returns a ready [CommandHandler].
// The builder should not be reused after calling Build.
func (b *CommandBuilder[TRequest, TResponse]) Build() CommandHandler[TRequest, TResponse] {
	h := b.handler
	for i := len(b.decorators) - 1; i >= 0; i-- {
		h = b.decorators[i](h)
	}
	return h
}
