// Package terror provides a simple error handling package
// TODO: add tags
// TODO: test
package terror

import (
	"github.com/w40141/gomyerr/internal"
)

// Error returns a new error
func Error(msg string) error {
	return internal.NewWhy(msg)
}

// Errorf returns a new origin with the formatted message
func Errorf(format string, a ...any) error {
	return internal.FormatWhy(format, a...)
}

// Wrap returns a new ThisError with the cause error
func Wrap(cause error, origin error) error {
	return internal.NewThisError(origin, cause)
}

// Wrapf returns a new ThisError with the cause error
func Wrapf(cause error, format string, a ...any) error {
	return internal.FormatThisError(cause, format, a...)
}

// WithStack returns a new ThisErrorWithStack with the cause error
func WithStack(e error) error {
	if e == nil {
		return nil
	}
	if _, ok := e.(*internal.ThisErrorWithStack); ok {
		return e
	}
	if _, ok := e.(*internal.ThisError); ok {
		return e.(*internal.ThisError).WithStack()
	}
	if _, ok := e.(*internal.Why); ok {
		return e.(*internal.Why).WithStack()
	}
	newErr := internal.FormatWhy("%v", e).WithStack()
	return newErr
}

// WrapStack returns a new ThisErrorWithStack with the cause error
func WrapStack(cause error, origin error) error {
	return internal.NewThisError(origin, cause).WithStack()
}

// WrapStackf returns a new ThisErrorWithStack with the cause error
func WrapStackf(cause error, format string, a ...any) error {
	return internal.FormatThisError(cause, format, a...).WithStack()
}

// Stack returns the stack trace
func Stack(err error) []map[string]any {
	if e, ok := err.(*internal.ThisErrorWithStack); ok {
		return e.Stack()
	}
	return nil
}
