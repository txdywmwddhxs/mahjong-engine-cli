package logging

// Logger records events to durable storage (e.g. file). It must NOT perform
// interactive terminal output or read user input.
type Logger interface {
	// Raw writes the string as-is (no prefix, no newline added).
	Raw(s string)
	Info(msg string, turn int)
	Debug(msg string, turn int)
	Error(msg string, turn int)
	Close() error
}

// NopLogger returns a logger that discards everything.
func NopLogger() Logger { return nopLogger{} }

type nopLogger struct{}

func (nopLogger) Raw(string)        {}
func (nopLogger) Info(string, int)  {}
func (nopLogger) Debug(string, int) {}
func (nopLogger) Error(string, int) {}
func (nopLogger) Close() error      { return nil }
