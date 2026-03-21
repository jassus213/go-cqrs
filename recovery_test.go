package go_cqrs

import (
	"context"
	"strings"
	"testing"
)

func TestRecovery_NoPanic(t *testing.T) {
	logger := newSpyLogger()

	uc := NewBuilder(okHandler).
		Use(Recovery[testQuery, testResult](logger)).
		Build()

	res, err := uc(context.Background(), testQuery{ID: 1})
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

func TestRecovery_CatchesPanic(t *testing.T) {
	logger := newSpyLogger()

	uc := NewBuilder(panicHandler).
		Use(Recovery[testQuery, testResult](logger)).
		Build()

	_, err := uc(context.Background(), testQuery{ID: 1})
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

func TestRecovery_ReturnsZeroValue(t *testing.T) {
	logger := newSpyLogger()

	uc := NewBuilder(panicHandler).
		Use(Recovery[testQuery, testResult](logger)).
		Build()

	res, _ := uc(context.Background(), testQuery{ID: 1})
	if res != (testResult{}) {
		t.Fatalf("expected zero value, got %+v", res)
	}
}
