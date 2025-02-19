// Package terror provides a simple error handling package
// TODO: add tags
// TODO: test
package terror

import (
	"errors"
	"fmt"
	"runtime"
)

const maxStackDepth = 10

type why struct {
	msg string
}

// thisError is an error with a message and a cause
type thisError struct {
	why   error
	cause error
}

// thisErrorWithStack is an error with a message, a cause, and a stack trace
type thisErrorWithStack struct {
	thisError
	stackTrace []uintptr
}

var (
	_ error = (*why)(nil)
	_ error = (*thisError)(nil)
	_ error = (*thisErrorWithStack)(nil)
)

func (w *why) Error() string {
	return w.msg
}

func (w *why) Is(target error) bool {
	t, ok := target.(*why)
	return ok && t.msg == w.msg
}

// New returns a new error
func New(msg string) error {
	return newWhy(msg)
}

// new returns a new why
func newWhy(msg string) *why {
	return &why{msg: msg}
}

func newWhyf(format string, a ...any) *why {
	return newWhy(fmt.Sprintf(format, a...))
}

// withCause returns a new ThisError with the cause error
func (w *why) withCause(cause error) *thisError {
	return &thisError{why: w, cause: cause}
}

// Format returns a new origin with the formatted message
func Format(format string, a ...any) error {
	return newWhyf(format, a...)
}

func (e *thisError) Error() string {
	return e.why.Error()
}

func (e *thisError) Unwrap() error {
	return e.cause
}

func (e *thisError) Is(target error) bool {
	if target == nil {
		return false
	}
	if e == target {
		return true
	}
	if errors.Is(e.why, target) || errors.Is(e.cause, target) {
		return true
	}
	return false
}

// Wrap returns a new ThisError with the cause error
func Wrap(cause error, origin error) error {
	return &thisError{why: origin, cause: cause}
}

// Wrapf returns a new ThisError with the cause error
func Wrapf(cause error, format string, a ...any) error {
	return newWhyf(format, a...).withCause(cause)
}

func (e *thisError) withStack() *thisErrorWithStack {
	stack := make([]uintptr, maxStackDepth)
	length := runtime.Callers(3, stack)
	return &thisErrorWithStack{thisError: *e, stackTrace: stack[:length]}
}

func (w *why) withStack() *thisErrorWithStack {
	stack := make([]uintptr, maxStackDepth)
	length := runtime.Callers(3, stack)
	return &thisErrorWithStack{thisError: thisError{why: w}, stackTrace: stack[:length]}
}

// WithStack returns a new ThisErrorWithStack with the cause error
func WithStack(e error) error {
	if e == nil {
		return nil
	}
	if _, ok := e.(*thisErrorWithStack); ok {
		return e
	}
	if _, ok := e.(*thisError); ok {
		return e.(*thisError).withStack()
	}
	if _, ok := e.(*why); ok {
		return e.(*why).withStack()
	}
	newErr := newWhyf("%v", e).withStack()
	return newErr
}

// WrapStack returns a new ThisErrorWithStack with the cause error
func WrapStack(cause error, origin error) error {
	e := &thisError{why: origin, cause: cause}
	return e.withStack()
}

// WrapStackf returns a new ThisErrorWithStack with the cause error
func WrapStackf(cause error, format string, a ...any) error {
	return newWhyf(format, a...).withCause(cause).withStack()
}

// Stack returns the stack trace
func Stack(err error) []map[string]any {
	if e, ok := err.(*thisErrorWithStack); ok {
		return e.stack()
	}
	return nil
}

func (e *thisErrorWithStack) stack() []map[string]any {
	if len(e.stackTrace) == 0 {
		if err, ok := e.cause.(*thisErrorWithStack); ok {
			return err.stack()
		}
		return nil
	}

	frames := runtime.CallersFrames(e.stackTrace)
	stack := make([]map[string]any, 0, maxStackDepth)
	for {
		frame, more := frames.Next()
		stack = append(stack, map[string]any{
			"file":     frame.File,
			"line":     frame.Line,
			"function": frame.Function,
		})
		if !more {
			break
		}
	}
	return stack
}
