// Package internal provides a simple error handling package
package internal

import (
	"fmt"
	"runtime"
)

var _ error = (*Why)(nil)

// Why is an error with a message
type Why struct {
	msg string
}

func (w *Why) Error() string {
	return w.msg
}

// NewWhy returns a new why error
func NewWhy(msg string) *Why {
	return &Why{msg: msg}
}

// FormatWhy returns a new why error with a formatted message
func FormatWhy(format string, a ...any) *Why {
	return NewWhy(fmt.Sprintf(format, a...))
}

// WithStack returns a new ThisErrorWithStack with the current stack trace
func (w *Why) WithStack() *ThisErrorWithStack {
	stack := make([]uintptr, maxStackDepth)
	length := runtime.Callers(3, stack)
	return &ThisErrorWithStack{
		ThisError:  ThisError{why: w},
		stackTrace: stack[:length],
	}
}

// WithCause returns a new ThisError with the cause error
func (w *Why) WithCause(cause error) *ThisError {
	return &ThisError{why: w, cause: cause}
}
