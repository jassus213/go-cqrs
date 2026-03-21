// Package go_cqrs provides a generic, composable use-case pipeline for Go applications.
//
// It implements the decorator (middleware) pattern around a simple function signature,
// allowing you to build request pipelines with cross-cutting concerns like logging,
// validation, recovery and timeouts — all without framework lock-in.
//
// Zero dependencies. Bring your own logger and validator.
package go_cqrs

import "context"

// UseCase is a function that handles a request and returns a response.
// This is the core abstraction — your business logic implements this signature.
//
//	type GetUserQuery struct { ID int64 }
//
//	var getUser UseCase[GetUserQuery, User] = func(ctx context.Context, req GetUserQuery) (User, error) {
//	    return repo.FindByID(ctx, req.ID)
//	}
type UseCase[TRequest any, TResponse any] func(ctx context.Context, req TRequest) (TResponse, error)
