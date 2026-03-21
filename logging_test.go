package go_cqrs

import (
	"context"
	"testing"
)

func TestLogging_Success(t *testing.T) {
	logger := newSpyLogger()

	uc := NewBuilder(UseCase[testQuery, testResult](okHandler)).
		Use(Logging[testQuery, testResult](logger)).
		Build()

	_, err := uc(context.Background(), testQuery{ID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !logger.hasMsg("executing") {
		t.Fatal("expected 'executing' log message")
	}
	if !logger.hasMsg("completed") {
		t.Fatal("expected 'completed' log message")
	}

	// Both should be debug level on success.
	for _, e := range logger.all() {
		if e.level != "debug" {
			t.Fatalf("expected debug level, got %q for msg %q", e.level, e.msg)
		}
	}
}

func TestLogging_Error(t *testing.T) {
	logger := newSpyLogger()

	uc := NewBuilder(UseCase[testQuery, testResult](errHandler)).
		Use(Logging[testQuery, testResult](logger)).
		Build()

	_, err := uc(context.Background(), testQuery{ID: 1})
	if err == nil {
		t.Fatal("expected error")
	}

	if !logger.hasLevel("error") {
		t.Fatal("expected error-level log when handler returns error")
	}
}

func TestLogging_ContainsOperationName(t *testing.T) {
	logger := newSpyLogger()

	uc := NewBuilder(UseCase[testQuery, testResult](okHandler)).
		Use(Logging[testQuery, testResult](logger)).
		Build()

	_, _ = uc(context.Background(), testQuery{ID: 1})

	found := false
	for _, e := range logger.all() {
		for i := 0; i+1 < len(e.args); i += 2 {
			if e.args[i] == "operation" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Fatal("expected 'operation' key in log args")
	}
}

func TestLogging_ContainsElapsed(t *testing.T) {
	logger := newSpyLogger()

	uc := NewBuilder(okHandler).
		Use(Logging[testQuery, testResult](logger)).
		Build()

	_, _ = uc(context.Background(), testQuery{ID: 1})

	found := false
	for _, e := range logger.all() {
		if e.msg == "completed" {
			for i := 0; i+1 < len(e.args); i += 2 {
				if e.args[i] == "elapsed" {
					found = true
					break
				}
			}
		}
	}
	if !found {
		t.Fatal("expected 'elapsed' key in completed log args")
	}
}
