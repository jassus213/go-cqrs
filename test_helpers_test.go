package go_cqrs

import (
	"context"
	"fmt"
	"sync"
)

// spyLogger captures log calls for assertions in tests.
type spyLogger struct {
	mu      sync.Mutex
	entries []logEntry
}

type logEntry struct {
	level string
	msg   string
	args  []any
}

func newSpyLogger() *spyLogger { return &spyLogger{} }

func (s *spyLogger) Info(_ context.Context, msg string, args ...any) {
	s.record("info", msg, args)
}

func (s *spyLogger) Debug(_ context.Context, msg string, args ...any) {
	s.record("debug", msg, args)
}

func (s *spyLogger) Error(_ context.Context, msg string, args ...any) {
	s.record("error", msg, args)
}

func (s *spyLogger) record(level, msg string, args []any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, logEntry{level: level, msg: msg, args: args})
}

func (s *spyLogger) all() []logEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make([]logEntry, len(s.entries))
	copy(cp, s.entries)
	return cp
}

func (s *spyLogger) hasLevel(level string) bool {
	for _, e := range s.all() {
		if e.level == level {
			return true
		}
	}
	return false
}

func (s *spyLogger) hasMsg(msg string) bool {
	for _, e := range s.all() {
		if e.msg == msg {
			return true
		}
	}
	return false
}

// --- common test types -------------------------------------------------------

type testQuery struct {
	ID int64
}

type testResult struct {
	Name string
}

type validatedCmd struct {
	Name string
}

func (c validatedCmd) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// handler that always succeeds.
func okHandler(_ context.Context, req testQuery) (testResult, error) {
	return testResult{Name: fmt.Sprintf("user-%d", req.ID)}, nil
}

// handler that always fails.
func errHandler(_ context.Context, _ testQuery) (testResult, error) {
	return testResult{}, fmt.Errorf("db connection failed")
}

// handler that panics.
func panicHandler(_ context.Context, _ testQuery) (testResult, error) {
	panic("something went terribly wrong")
}
