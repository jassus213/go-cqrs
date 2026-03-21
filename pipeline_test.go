package go_cqrs

import (
	"context"
	"strings"
	"testing"
)

// Integration tests — full pipeline with multiple decorators.

func TestPipeline_DefaultBuilder_Success(t *testing.T) {
	logger := newSpyLogger()

	uc := NewDefaultBuilder(logger, okHandler).Build()

	res, err := uc(context.Background(), testQuery{ID: 99})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "user-99" {
		t.Fatalf("got %q, want %q", res.Name, "user-99")
	}
}

func TestPipeline_DefaultBuilder_HandlerError(t *testing.T) {
	logger := newSpyLogger()

	uc := NewDefaultBuilder(logger, errHandler).Build()

	_, err := uc(context.Background(), testQuery{ID: 1})
	if err == nil {
		t.Fatal("expected error")
	}
	if !logger.hasLevel("error") {
		t.Fatal("expected error to be logged")
	}
}

func TestPipeline_DefaultBuilder_Panic(t *testing.T) {
	logger := newSpyLogger()

	uc := NewDefaultBuilder(logger, panicHandler).Build()

	_, err := uc(context.Background(), testQuery{ID: 1})
	if err == nil {
		t.Fatal("expected error from panic")
	}
	if !strings.Contains(err.Error(), "panic") {
		t.Fatalf("error should mention panic, got: %v", err)
	}
}

func TestPipeline_DefaultBuilder_ValidationFails(t *testing.T) {
	logger := newSpyLogger()

	called := false
	handler := func(_ context.Context, req validatedCmd) (None, error) {
		called = true
		return None{}, nil
	}

	uc := NewDefaultBuilder(logger, handler).Build()

	_, err := uc(context.Background(), validatedCmd{Name: ""})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if called {
		t.Fatal("handler should not be called when validation fails")
	}
}

func TestPipeline_DefaultBuilder_ValidationPasses(t *testing.T) {
	logger := newSpyLogger()

	handler := func(_ context.Context, req validatedCmd) (None, error) {
		return None{}, nil
	}

	uc := NewDefaultBuilder(logger, handler).Build()

	_, err := uc(context.Background(), validatedCmd{Name: "Bob"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPipeline_CustomDecorator(t *testing.T) {
	var trace []string

	custom := func(next UseCase[testQuery, testResult]) UseCase[testQuery, testResult] {
		return func(ctx context.Context, req testQuery) (testResult, error) {
			trace = append(trace, "custom")
			return next(ctx, req)
		}
	}

	uc := NewBuilder(okHandler).
		Use(custom).
		Build()

	_, _ = uc(context.Background(), testQuery{ID: 1})
	if len(trace) != 1 || trace[0] != "custom" {
		t.Fatalf("custom decorator not called, trace: %v", trace)
	}
}

func TestPipeline_NoneResponse(t *testing.T) {
	handler := func(_ context.Context, _ testQuery) (None, error) {
		return None{}, nil
	}

	uc := NewBuilder(handler).Build()

	res, err := uc(context.Background(), testQuery{ID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != (None{}) {
		t.Fatalf("expected None{}, got %+v", res)
	}
}
