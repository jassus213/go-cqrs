package go_cqrs

// UseCaseDecorator wraps a [UseCase] and returns a new one with added behavior.
// Decorators are applied in order: the first decorator in the slice wraps the outermost layer.
//
// Example — a simple timing decorator:
//
//	func Timer[Req, Res any]() UseCaseDecorator[Req, Res] {
//	    return func(next UseCase[Req, Res]) UseCase[Req, Res] {
//	        return func(ctx context.Context, req Req) (Res, error) {
//	            start := time.Now()
//	            res, err := next(ctx, req)
//	            fmt.Println("elapsed:", time.Since(start))
//	            return res, err
//	        }
//	    }
//	}
type UseCaseDecorator[TRequest any, TResponse any] func(next UseCase[TRequest, TResponse]) UseCase[TRequest, TResponse]
