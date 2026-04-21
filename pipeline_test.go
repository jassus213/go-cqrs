package go_cqrs

import (
	"context"
	"strings"
	"testing"
)

// Integration tests — full pipeline with multiple decorators.

func TestPipeline_DefaultQueryBuilder_Success(t *testing.T) {
	logger := newSpyLogger()

	h := NewDefaultQueryBuilder(logger, QueryHandler[testQuery, testResult](okQueryHandler)).Build()

	res, err := h(context.Background(), testQuery{ID: 99})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "user-99" {
		t.Fatalf("got %q, want %q", res.Name, "user-99")
	}
}

func TestPipeline_DefaultQueryBuilder_HandlerError(t *testing.T) {
	logger := newSpyLogger()

	h := NewDefaultQueryBuilder(logger, QueryHandler[testQuery, testResult](errQueryHandler)).Build()

	_, err := h(context.Background(), testQuery{ID: 1})
	if err == nil {
		t.Fatal("expected error")
	}
	if !logger.hasLevel("error") {
		t.Fatal("expected error to be logged")
	}
}

func TestPipeline_DefaultQueryBuilder_Panic(t *testing.T) {
	logger := newSpyLogger()

	h := NewDefaultQueryBuilder(logger, QueryHandler[testQuery, testResult](panicQueryHandler)).Build()

	_, err := h(context.Background(), testQuery{ID: 1})
	if err == nil {
		t.Fatal("expected error from panic")
	}
	if !strings.Contains(err.Error(), "panic") {
		t.Fatalf("error should mention panic, got: %v", err)
	}
}

func TestPipeline_DefaultCommandBuilder_ValidationFails(t *testing.T) {
	logger := newSpyLogger()

	called := false
	handler := func(_ context.Context, req validatedCmd) (None, error) {
		called = true
		return None{}, nil
	}

	h := NewDefaultCommandBuilder(logger, CommandHandler[validatedCmd, None](handler)).Build()

	_, err := h(context.Background(), validatedCmd{Name: ""})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if called {
		t.Fatal("handler should not be called when validation fails")
	}
}

func TestPipeline_DefaultCommandBuilder_ValidationPasses(t *testing.T) {
	logger := newSpyLogger()

	handler := func(_ context.Context, req validatedCmd) (None, error) {
		return None{}, nil
	}

	h := NewDefaultCommandBuilder(logger, CommandHandler[validatedCmd, None](handler)).Build()

	_, err := h(context.Background(), validatedCmd{Name: "Bob"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPipeline_CustomQueryDecorator(t *testing.T) {
	var trace []string

	custom := func(next QueryHandler[testQuery, testResult]) QueryHandler[testQuery, testResult] {
		return func(ctx context.Context, req testQuery) (testResult, error) {
			trace = append(trace, "custom")
			return next(ctx, req)
		}
	}

	h := NewQueryBuilder(QueryHandler[testQuery, testResult](okQueryHandler)).
		Use(custom).
		Build()

	_, _ = h(context.Background(), testQuery{ID: 1})
	if len(trace) != 1 || trace[0] != "custom" {
		t.Fatalf("custom decorator not called, trace: %v", trace)
	}
}

func TestPipeline_NoneResponse(t *testing.T) {
	handler := func(_ context.Context, _ validatedCmd) (None, error) {
		return None{}, nil
	}

	h := NewCommandBuilder(CommandHandler[validatedCmd, None](handler)).Build()

	res, err := h(context.Background(), validatedCmd{Name: "Alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != (None{}) {
		t.Fatalf("expected None{}, got %+v", res)
	}
}
