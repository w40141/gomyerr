// Package gomyerr provides a simple error handling package
// TODO: add stack trace
// TODO: add tags
// TODO: test
package gomyerr

import (
	"fmt"
	"runtime"
)

const maxStackDepth = 10

type why struct {
	msg string
}

// thisErrorWithStack is an error with a message, a cause, and a stack trace
type thisErrorWithStack struct {
	thisError
	stackTrace []uintptr
}

// thisError is an error with a message and a cause
type thisError struct {
	why   error
	cause error
}

var (
	_ error = (*why)(nil)
	_ error = (*thisError)(nil)
	_ error = (*thisErrorWithStack)(nil)
)

func (w why) Error() string {
	return w.msg
}

// new returns a new ThisError
func newWhy(msg string) why {
	return why{msg: msg}
}

func newWhyf(format string, args ...any) why {
	return newWhy(fmt.Sprintf(format, args...))
}

// withCause returns a new ThisError with the cause error
func (w why) withCause(cause error) thisError {
	return thisError{why: w, cause: cause}
}

// Format returns a new origin with the formatted message
func Format(format string, args ...any) error {
	return newWhyf(format, args...)
}

func (e thisError) Error() string {
	return e.why.Error()
}

func (e thisError) Unwrap() error {
	return e.cause
}

// Wrap returns a new ThisError with the cause error
func Wrap(cause error, format string, args ...any) error {
	return newWhyf(format, args...).withCause(cause)
}

func (e thisError) withStack() thisErrorWithStack {
	stack := make([]uintptr, maxStackDepth)
	length := runtime.Callers(3, stack)
	return thisErrorWithStack{thisError: e, stackTrace: stack[:length]}
}

// WrapStack returns a new ThisErrorWithStack with the cause error
func WrapStack(cause error, format string, args ...any) error {
	return newWhyf(format, args...).withCause(cause).withStack()
}

// Stack returns the stack trace
func Stack(err error) []map[string]any {
	if e, ok := err.(*thisErrorWithStack); ok {
		return e.stack()
	}
	return nil
}

func (e thisErrorWithStack) stack() []map[string]any {
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
