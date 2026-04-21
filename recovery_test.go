package go_cqrs

import (
	"context"
	"strings"
	"testing"
)

func TestQueryRecovery_NoPanic(t *testing.T) {
	logger := newSpyLogger()

	h := NewQueryBuilder(QueryHandler[testQuery, testResult](okQueryHandler)).
		Use(QueryRecovery[testQuery, testResult](logger)).
		Build()

	res, err := h(context.Background(), testQuery{ID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "user-1" {
		t.Fatalf("got %q, want %q", res.Name, "user-1")
	}
	if logger.hasLevel("error") {
		t.Fatal("no error should be logged when there is no panic")
	}
}

func TestQueryRecovery_CatchesPanic(t *testing.T) {
	logger := newSpyLogger()

	h := NewQueryBuilder(QueryHandler[testQuery, testResult](panicQueryHandler)).
		Use(QueryRecovery[testQuery, testResult](logger)).
		Build()

	_, err := h(context.Background(), testQuery{ID: 1})
	if err == nil {
		t.Fatal("expected error from recovered panic")
	}
	if !strings.Contains(err.Error(), "panic") {
		t.Fatalf("error should mention panic, got: %v", err)
	}
	if !logger.hasMsg("recovered from panic") {
		t.Fatal("expected panic to be logged")
	}
}

func TestQueryRecovery_ReturnsZeroValue(t *testing.T) {
	logger := newSpyLogger()

	h := NewQueryBuilder(QueryHandler[testQuery, testResult](panicQueryHandler)).
		Use(QueryRecovery[testQuery, testResult](logger)).
		Build()

	res, _ := h(context.Background(), testQuery{ID: 1})
	if res != (testResult{}) {
		t.Fatalf("expected zero value, got %+v", res)
	}
}

func TestCommandRecovery_CatchesPanic(t *testing.T) {
	logger := newSpyLogger()

	h := NewCommandBuilder(CommandHandler[validatedCmd, None](panicCommandHandler)).
		Use(CommandRecovery[validatedCmd, None](logger)).
		Build()

	_, err := h(context.Background(), validatedCmd{Name: "Alice"})
	if err == nil {
		t.Fatal("expected error from recovered panic")
	}
	if !strings.Contains(err.Error(), "panic") {
		t.Fatalf("error should mention panic, got: %v", err)
	}
	if !logger.hasMsg("recovered from panic") {
		t.Fatal("expected panic to be logged")
	}
}
