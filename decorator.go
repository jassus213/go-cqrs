package go_cqrs

// QueryDecorator wraps a [QueryHandler] and returns a new one with added behavior.
// Decorators are applied in order: the first in the list is the outermost layer of the pipeline.
type QueryDecorator[TRequest any, TResponse any] func(next QueryHandler[TRequest, TResponse]) QueryHandler[TRequest, TResponse]

// CommandDecorator wraps a [CommandHandler] and returns a new one with added behavior.
// Decorators are applied in order: the first in the list is the outermost layer of the pipeline.
type CommandDecorator[TRequest any, TResponse any] func(next CommandHandler[TRequest, TResponse]) CommandHandler[TRequest, TResponse]
