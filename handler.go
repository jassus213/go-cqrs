// Package go_cqrs — generic, composable pipelines for Go
// with an explicit split between commands (writes) and queries (reads).
//
// Zero dependencies. Bring your own logger and validator.
package go_cqrs

import "context"

// QueryHandler handles a read operation: accepts a request, returns a result.
//
//	var getUser QueryHandler[GetUserQuery, User] = func(ctx context.Context, req GetUserQuery) (User, error) {
//	    return repo.FindByID(ctx, req.ID)
//	}
type QueryHandler[TRequest any, TResponse any] func(ctx context.Context, req TRequest) (TResponse, error)

// CommandHandler handles a write operation: accepts a request, returns a result.
// Use [None] as TResponse for commands that return no data.
//
//	// Command returning the ID of the created entity:
//	var createUser CommandHandler[CreateUserCmd, uuid.UUID] = func(ctx context.Context, req CreateUserCmd) (uuid.UUID, error) {
//	    return repo.Insert(ctx, req)
//	}
//
//	// Command with no return value:
//	var deleteUser CommandHandler[DeleteUserCmd, None] = func(ctx context.Context, req DeleteUserCmd) (None, error) {
//	    return None{}, repo.Delete(ctx, req.ID)
//	}
type CommandHandler[TRequest any, TResponse any] func(ctx context.Context, req TRequest) (TResponse, error)
