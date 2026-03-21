package go_cqrs

import (
	"context"
	"testing"
)

func TestNewBuilder_NoDecorators(t *testing.T) {
	uc := NewBuilder(okHandler).Build()

	res, err := uc(context.Background(), testQuery{ID: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "user-42" {
		t.Fatalf("got %q, want %q", res.Name, "user-42")
	}
}

func TestNewBuilder_DecoratorOrder(t *testing.T) {
	var order []string

	first := func(next UseCase[testQuery, testResult]) UseCase[testQuery, testResult] {
		return func(ctx context.Context, req testQuery) (testResult, error) {
			order = append(order, "first-before")
			res, err := next(ctx, req)
			order = append(order, "first-after")
			return res, err
		}
	}

	second := func(next UseCase[testQuery, testResult]) UseCase[testQuery, testResult] {
		return func(ctx context.Context, req testQuery) (testResult, error) {
			order = append(order, "second-before")
			res, err := next(ctx, req)
			order = append(order, "second-after")
			return res, err
		}
	}

	uc := NewBuilder(okHandler).
		Use(first, second).
		Build()

	_, _ = uc(context.Background(), testQuery{ID: 1})

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

func TestNewDefaultBuilder_IncludesAllDecorators(t *testing.T) {
	logger := newSpyLogger()

	uc := NewDefaultBuilder(logger, okHandler).Build()

	res, err := uc(context.Background(), testQuery{ID: 7})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "user-7" {
		t.Fatalf("got %q, want %q", res.Name, "user-7")
	}

	// Logging decorator should have produced debug entries.
	if !logger.hasMsg("executing") {
		t.Fatal("expected 'executing' log message from Logging decorator")
	}
	if !logger.hasMsg("completed") {
		t.Fatal("expected 'completed' log message from Logging decorator")
	}
}

func TestBuilder_Use_Appends(t *testing.T) {
	called := false
	decorator := func(next UseCase[testQuery, testResult]) UseCase[testQuery, testResult] {
		return func(ctx context.Context, req testQuery) (testResult, error) {
			called = true
			return next(ctx, req)
		}
	}

	uc := NewBuilder(okHandler).
		Use(decorator).
		Build()

	_, _ = uc(context.Background(), testQuery{ID: 1})
	if !called {
		t.Fatal("decorator was not called")
	}
}
