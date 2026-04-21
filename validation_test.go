package go_cqrs

import (
	"context"
	"strings"
	"testing"
)

func TestCommandValidation_PassesValid(t *testing.T) {
	handler := func(_ context.Context, req validatedCmd) (None, error) {
		return None{}, nil
	}

	h := NewCommandBuilder(CommandHandler[validatedCmd, None](handler)).
		Use(CommandValidation[validatedCmd, None]()).
		Build()

	_, err := h(context.Background(), validatedCmd{Name: "Alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommandValidation_RejectsInvalid(t *testing.T) {
	called := false
	handler := func(_ context.Context, req validatedCmd) (None, error) {
		called = true
		return None{}, nil
	}

	h := NewCommandBuilder(CommandHandler[validatedCmd, None](handler)).
		Use(CommandValidation[validatedCmd, None]()).
		Build()

	_, err := h(context.Background(), validatedCmd{Name: ""})
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

func TestQueryValidation_SkipsNonValidator(t *testing.T) {
	h := NewQueryBuilder(QueryHandler[testQuery, testResult](okQueryHandler)).
		Use(QueryValidation[testQuery, testResult]()).
		Build()

	res, err := h(context.Background(), testQuery{ID: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "user-5" {
		t.Fatalf("got %q, want %q", res.Name, "user-5")
	}
}
