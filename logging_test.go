package go_cqrs

import (
	"context"
	"testing"
)

func TestQueryLogging_Success(t *testing.T) {
	logger := newSpyLogger()

	h := NewQueryBuilder(QueryHandler[testQuery, testResult](okQueryHandler)).
		Use(QueryLogging[testQuery, testResult](logger)).
		Build()

	_, err := h(context.Background(), testQuery{ID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !logger.hasMsg("query started") {
		t.Fatal("expected 'query started' log message")
	}
	if !logger.hasMsg("query completed") {
		t.Fatal("expected 'query completed' log message")
	}

	for _, e := range logger.all() {
		if e.level != "debug" {
			t.Fatalf("expected debug level, got %q for msg %q", e.level, e.msg)
		}
	}
}

func TestQueryLogging_Error(t *testing.T) {
	logger := newSpyLogger()

	h := NewQueryBuilder(QueryHandler[testQuery, testResult](errQueryHandler)).
		Use(QueryLogging[testQuery, testResult](logger)).
		Build()

	_, err := h(context.Background(), testQuery{ID: 1})
	if err == nil {
		t.Fatal("expected error")
	}

	if !logger.hasLevel("error") {
		t.Fatal("expected error-level log when handler returns error")
	}
}

func TestQueryLogging_ContainsOperationName(t *testing.T) {
	logger := newSpyLogger()

	h := NewQueryBuilder(QueryHandler[testQuery, testResult](okQueryHandler)).
		Use(QueryLogging[testQuery, testResult](logger)).
		Build()

	_, _ = h(context.Background(), testQuery{ID: 1})

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

func TestQueryLogging_ContainsElapsed(t *testing.T) {
	logger := newSpyLogger()

	h := NewQueryBuilder(QueryHandler[testQuery, testResult](okQueryHandler)).
		Use(QueryLogging[testQuery, testResult](logger)).
		Build()

	_, _ = h(context.Background(), testQuery{ID: 1})

	found := false
	for _, e := range logger.all() {
		if e.msg == "query completed" {
			for i := 0; i+1 < len(e.args); i += 2 {
				if e.args[i] == "elapsed" {
					found = true
					break
				}
			}
		}
	}
	if !found {
		t.Fatal("expected 'elapsed' key in query completed log args")
	}
}

func TestCommandLogging_Success(t *testing.T) {
	logger := newSpyLogger()

	h := NewCommandBuilder(CommandHandler[validatedCmd, None](okCommandHandler)).
		Use(CommandLogging[validatedCmd, None](logger)).
		Build()

	_, err := h(context.Background(), validatedCmd{Name: "Alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !logger.hasMsg("command started") {
		t.Fatal("expected 'command started' log message")
	}
	if !logger.hasMsg("command completed") {
		t.Fatal("expected 'command completed' log message")
	}
}

func TestCommandLogging_Error(t *testing.T) {
	logger := newSpyLogger()

	h := NewCommandBuilder(CommandHandler[validatedCmd, None](errCommandHandler)).
		Use(CommandLogging[validatedCmd, None](logger)).
		Build()

	_, err := h(context.Background(), validatedCmd{Name: "Alice"})
	if err == nil {
		t.Fatal("expected error")
	}

	if !logger.hasLevel("error") {
		t.Fatal("expected error-level log when handler returns error")
	}
}
