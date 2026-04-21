package go_cqrs

import (
	"context"
	"testing"
)

func TestNewQueryBuilder_NoDecorators(t *testing.T) {
	h := NewQueryBuilder(QueryHandler[testQuery, testResult](okQueryHandler)).Build()

	res, err := h(context.Background(), testQuery{ID: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "user-42" {
		t.Fatalf("got %q, want %q", res.Name, "user-42")
	}
}

func TestNewQueryBuilder_DecoratorOrder(t *testing.T) {
	var order []string

	first := func(next QueryHandler[testQuery, testResult]) QueryHandler[testQuery, testResult] {
		return func(ctx context.Context, req testQuery) (testResult, error) {
			order = append(order, "first-before")
			res, err := next(ctx, req)
			order = append(order, "first-after")
			return res, err
		}
	}

	second := func(next QueryHandler[testQuery, testResult]) QueryHandler[testQuery, testResult] {
		return func(ctx context.Context, req testQuery) (testResult, error) {
			order = append(order, "second-before")
			res, err := next(ctx, req)
			order = append(order, "second-after")
			return res, err
		}
	}

	h := NewQueryBuilder(QueryHandler[testQuery, testResult](okQueryHandler)).
		Use(first, second).
		Build()

	_, _ = h(context.Background(), testQuery{ID: 1})

	expected := []string{"first-before", "second-before", "second-after", "first-after"}
	if len(order) != len(expected) {
		t.Fatalf("got %v, want %v", order, expected)
	}
	for i := range expected {
		if order[i] != expected[i] {
			t.Fatalf("order[%d] = %q, want %q", i, order[i], expected[i])
		}
	}
}

func TestNewDefaultQueryBuilder_IncludesAllDecorators(t *testing.T) {
	logger := newSpyLogger()

	h := NewDefaultQueryBuilder(logger, QueryHandler[testQuery, testResult](okQueryHandler)).Build()

	res, err := h(context.Background(), testQuery{ID: 7})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "user-7" {
		t.Fatalf("got %q, want %q", res.Name, "user-7")
	}

	if !logger.hasMsg("query started") {
		t.Fatal("expected 'query started' log message from QueryLogging decorator")
	}
	if !logger.hasMsg("query completed") {
		t.Fatal("expected 'query completed' log message from QueryLogging decorator")
	}
}

func TestNewCommandBuilder_NoDecorators(t *testing.T) {
	h := NewCommandBuilder(CommandHandler[validatedCmd, None](okCommandHandler)).Build()

	_, err := h(context.Background(), validatedCmd{Name: "Alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewDefaultCommandBuilder_IncludesAllDecorators(t *testing.T) {
	logger := newSpyLogger()

	h := NewDefaultCommandBuilder(logger, CommandHandler[validatedCmd, None](okCommandHandler)).Build()

	_, err := h(context.Background(), validatedCmd{Name: "Bob"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !logger.hasMsg("command started") {
		t.Fatal("expected 'command started' log message from CommandLogging decorator")
	}
	if !logger.hasMsg("command completed") {
		t.Fatal("expected 'command completed' log message from CommandLogging decorator")
	}
}

func TestQueryBuilder_Use_Appends(t *testing.T) {
	called := false
	decorator := func(next QueryHandler[testQuery, testResult]) QueryHandler[testQuery, testResult] {
		return func(ctx context.Context, req testQuery) (testResult, error) {
			called = true
			return next(ctx, req)
		}
	}

	h := NewQueryBuilder(QueryHandler[testQuery, testResult](okQueryHandler)).
		Use(decorator).
		Build()

	_, _ = h(context.Background(), testQuery{ID: 1})
	if !called {
		t.Fatal("decorator was not called")
	}
}
