// Package terror provides a simple error handling package
// TODO: add tags
// TODO: test
package terror

import (
	"github.com/w40141/gomyerr/internal"
)

// New returns a new error
func New(msg string) error {
	return internal.NewWhat(msg)
}

// Format returns a new origin with the formatted message
func Format(format string, a ...any) error {
	return internal.FormatWhat(format, a...)
}

// Wrap returns a new ThisError with the cause error
func Wrap(cause error, origin error) error {
	return internal.NewWhy(origin, cause)
}

// Wrapf returns a new ThisError with the cause error
func Wrapf(cause error, format string, a ...any) error {
	return internal.FormatMyErr(cause, format, a...)
}

// WithStack returns a new ThisErrorWithStack with the cause error
func WithStack(e error) error {
	if e == nil {
		return nil
	}
	if _, ok := e.(*internal.MyErr); ok {
		return e
	}
	if _, ok := e.(*internal.Why); ok {
		return e.(*internal.Why).WithStack()
	}
	if _, ok := e.(*internal.What); ok {
		return e.(*internal.What).WithStack()
	}
	newErr := internal.FormatWhat("%v", e).WithStack()
	return newErr
}

// WrapStack returns a new ThisErrorWithStack with the cause error
func WrapStack(cause error, origin error) error {
	return internal.NewWhy(origin, cause).WithStack()
}

// WrapStackf returns a new ThisErrorWithStack with the cause error
func WrapStackf(cause error, format string, a ...any) error {
	return internal.FormatMyErr(cause, format, a...).WithStack()
}

// Stack returns the stack trace
func Stack(err error) []map[string]any {
	if e, ok := err.(*internal.MyErr); ok {
		return e.Stack()
	}
	return nil
}
