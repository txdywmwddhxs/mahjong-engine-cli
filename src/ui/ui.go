package ui

// UI is the interactive surface for the user. It must NOT write to log files.
type UI interface {
	// Info prints a user-facing informational line (with INFO prefix).
	Info(msg string)

	// Plain prints raw text without adding prefix or newline.
	Plain(msg string)

	// Plainln prints raw text with a trailing newline.
	Plainln(msg string)

	// PromptInfo prints an INFO-prefixed prompt and reads one line of input.
	PromptInfo(msg string) (string, error)

	// PromptPlain prints a raw prompt and reads one line of input.
	PromptPlain(msg string) (string, error)
}
