package logging

import "fmt"

// Logger records events to durable storage (e.g. file). It must NOT perform
// interactive terminal output or read user input.
type Logger interface {
	// Raw writes the string as-is (no prefix, no newline added).
	Raw(s string)
	// WriteBlock writes and syncs a complete replay block, reporting failures.
	WriteBlock(s string) error
	Info(msg string, turn int)
	Debug(msg string, turn int)
	Error(msg string, turn int)
	Close() error
}

// NopLogger returns a logger that discards everything.
func NopLogger() Logger { return nopLogger{} }

type nopLogger struct{}

func (nopLogger) WriteBlock(string) error { return fmt.Errorf("file logging is unavailable") }

func (nopLogger) Raw(string)        {}
func (nopLogger) Info(string, int)  {}
func (nopLogger) Debug(string, int) {}
func (nopLogger) Error(string, int) {}
func (nopLogger) Close() error      { return nil }
