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
	// Simulate a DB lookup.
	return User{ID: req.ID, Name: "Alice"}, nil
}

func createUser(_ context.Context, req CreateUserCmd) (cqrs.None, error) {
	fmt.Printf("created user: %s\n", req.Name)
	return cqrs.None{}, nil
}

// --- Main --------------------------------------------------------------------

func main() {
	logger := SlogAdapter{L: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))}

	// Query: uses the default pipeline (Recovery → Logging → Validation).
	getUserUC := cqrs.NewDefaultBuilder(logger, getUser).Build()

	user, err := getUserUC(context.Background(), GetUserQuery{ID: 1})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("got user: %+v\n", user)

	// Command: custom pipeline with only logging.
	createUserUC := cqrs.NewBuilder(createUser).
		Use(cqrs.Logging[CreateUserCmd, cqrs.None](logger)).
		Build()

	_, err = createUserUC(context.Background(), CreateUserCmd{Name: "Bob"})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Validation error example.
	_, err = createUserUC(context.Background(), CreateUserCmd{Name: ""})
	fmt.Println("validation skipped (no Validation decorator):", err)

	// With validation.
	createUserValidated := cqrs.NewDefaultBuilder(logger, createUser).Build()

	_, err = createUserValidated(context.Background(), CreateUserCmd{Name: ""})
	fmt.Println("validation caught:", err)
}
