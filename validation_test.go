package go_cqrs

import (
	"context"
	"strings"
	"testing"
)

func TestValidation_PassesValid(t *testing.T) {
	handler := func(_ context.Context, req validatedCmd) (None, error) {
		return None{}, nil
	}

	uc := NewBuilder(handler).
		Use(Validation[validatedCmd, None]()).
		Build()

	_, err := uc(context.Background(), validatedCmd{Name: "Alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidation_RejectsInvalid(t *testing.T) {
	called := false
	handler := func(_ context.Context, req validatedCmd) (None, error) {
		called = true
		return None{}, nil
	}

	uc := NewBuilder(handler).
		Use(Validation[validatedCmd, None]()).
		Build()

	_, err := uc(context.Background(), validatedCmd{Name: ""})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "name is required") {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("handler should not be called when validation fails")
	}
}

func TestValidation_SkipsNonValidator(t *testing.T) {
	// testQuery does not implement Validator — should pass through.
	uc := NewBuilder(okHandler).
		Use(Validation[testQuery, testResult]()).
		Build()

	res, err := uc(context.Background(), testQuery{ID: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "user-5" {
		t.Fatalf("got %q, want %q", res.Name, "user-5")
	}
}
