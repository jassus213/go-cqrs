// Package main demonstrates a minimal go-cqrs pipeline.
//
// Run with: go run ./examples/basic
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	cqrs "github.com/jassus213/go-cqrs"
)

// --- Logger adapter (slog) ---------------------------------------------------

type SlogAdapter struct{ L *slog.Logger }

func (a SlogAdapter) Info(ctx context.Context, msg string, args ...any) {
	a.L.InfoContext(ctx, msg, args...)
}
func (a SlogAdapter) Warn(ctx context.Context, msg string, args ...any) {
	a.L.WarnContext(ctx, msg, args...)
}
func (a SlogAdapter) Debug(ctx context.Context, msg string, args ...any) {
	a.L.DebugContext(ctx, msg, args...)
}
func (a SlogAdapter) Error(ctx context.Context, msg string, args ...any) {
	a.L.ErrorContext(ctx, msg, args...)
}

// --- Request / Response ------------------------------------------------------

type GetUserQuery struct {
	ID int64
}

type User struct {
	ID   int64
	Name string
}

// --- Request with validation -------------------------------------------------

type CreateUserCmd struct {
	Name string
}

func (c CreateUserCmd) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// --- Handlers ----------------------------------------------------------------

func getUser(_ context.Context, req GetUserQuery) (User, error) {
	return User{ID: req.ID, Name: "Alice"}, nil
}

func createUser(_ context.Context, req CreateUserCmd) (cqrs.None, error) {
	fmt.Printf("created user: %s\n", req.Name)
	return cqrs.None{}, nil
}

// --- Main --------------------------------------------------------------------

func main() {
	logger := SlogAdapter{L: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))}

	// Query: uses the default pipeline (QueryRecovery → QueryLogging → QueryValidation).
	getUserHandler := cqrs.NewDefaultQueryBuilder(logger, cqrs.QueryHandler[GetUserQuery, User](getUser)).Build()

	user, err := getUserHandler(context.Background(), GetUserQuery{ID: 1})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("got user: %+v\n", user)

	// Command: custom pipeline with only logging.
	createUserHandler := cqrs.NewCommandBuilder(cqrs.CommandHandler[CreateUserCmd, cqrs.None](createUser)).
		Use(cqrs.CommandLogging[CreateUserCmd, cqrs.None](logger)).
		Build()

	_, err = createUserHandler(context.Background(), CreateUserCmd{Name: "Bob"})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Validation skipped — no CommandValidation decorator in this pipeline.
	_, err = createUserHandler(context.Background(), CreateUserCmd{Name: ""})
	fmt.Println("validation skipped (no CommandValidation decorator):", err)

	// With full default pipeline — CommandValidation catches the empty name.
	createUserValidated := cqrs.NewDefaultCommandBuilder(logger, cqrs.CommandHandler[CreateUserCmd, cqrs.None](createUser)).Build()

	_, err = createUserValidated(context.Background(), CreateUserCmd{Name: ""})
	fmt.Println("validation caught:", err)
}
