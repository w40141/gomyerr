// Package internal provides a simple error handling package
package internal

import (
	"errors"
	"runtime"
)

var _ error = (*ThisError)(nil)

// ThisError is an error with a message and a cause
type ThisError struct {
	why   error
	cause error
}

func (e *ThisError) Error() string {
	return e.why.Error()
}

func (e *ThisError) Unwrap() error {
	return e.cause
}

// Is reports whether the target error is the same as the current error
func (e *ThisError) Is(target error) bool {
	if errors.Is(e.why, target) || errors.Is(e.cause, target) {
		return true
	}
	return false
}

// WithStack returns a new ThisErrorWithStack with the current stack trace
func (e *ThisError) WithStack() *ThisErrorWithStack {
	stack := make([]uintptr, maxStackDepth)
	length := runtime.Callers(3, stack)
	return &ThisErrorWithStack{ThisError: *e, stackTrace: stack[:length]}
}

// NewThisError returns a new ThisError
func NewThisError(origin, cause error) *ThisError {
	return &ThisError{why: origin, cause: cause}
}

// FormatThisError returns a new ThisError with a formatted message
func FormatThisError(causer error, format string, a ...any) *ThisError {
	return &ThisError{why: FormatWhy(format, a...), cause: causer}
}
