package logging

import (
	"fmt"
	"os"
	"sync"
)

const (
	infoFmt  = "INFO[%d]: %s\n"
	debugFmt = "DEBUG[%d]: %s\n"
	errorFmt = "ERROR[%d]: %s\n"
)

type FileLogger struct {
	mu sync.Mutex
	f  *os.File
}

func NewFileLogger(path string) (*FileLogger, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o666)
	if err != nil {
		return nil, err
	}
	return &FileLogger{f: f}, nil
}

func (l *FileLogger) Raw(s string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f == nil {
		return
	}
	_, _ = l.f.WriteString(s)
}

func (l *FileLogger) Info(msg string, turn int)  { l.writef(infoFmt, msg, turn) }
func (l *FileLogger) Debug(msg string, turn int) { l.writef(debugFmt, msg, turn) }
func (l *FileLogger) Error(msg string, turn int) { l.writef(errorFmt, msg, turn) }

func (l *FileLogger) writef(format, msg string, turn int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f == nil {
		return
	}
	_, _ = l.f.WriteString(fmt.Sprintf(format, turn, msg))
}

func (l *FileLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f == nil {
		return nil
	}
	err := l.f.Close()
	l.f = nil
	return err
}
