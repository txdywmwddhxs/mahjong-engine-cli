package ui

import (
	"bufio"
	"fmt"
	"io"
)

type Console struct {
	out     io.Writer
	scanner *bufio.Scanner
}

func NewConsole(in io.Reader, out io.Writer) *Console {
	return &Console{
		out:     out,
		scanner: bufio.NewScanner(in),
	}
}

func (c *Console) Info(msg string) {
	_, _ = fmt.Fprintf(c.out, "INFO: %s\n", msg)
}

func (c *Console) Plain(msg string) {
	_, _ = fmt.Fprint(c.out, msg)
}

func (c *Console) Plainln(msg string) {
	_, _ = fmt.Fprintln(c.out, msg)
}

func (c *Console) PromptInfo(msg string) (string, error) {
	_, _ = fmt.Fprintf(c.out, "INFO: %s: ", msg)
	return c.readLine()
}

func (c *Console) PromptPlain(msg string) (string, error) {
	_, _ = fmt.Fprint(c.out, msg)
	return c.readLine()
}

func (c *Console) readLine() (string, error) {
	if c.scanner.Scan() {
		return c.scanner.Text(), nil
	}
	if err := c.scanner.Err(); err != nil {
		return "", err
	}
	return "", io.EOF
}
