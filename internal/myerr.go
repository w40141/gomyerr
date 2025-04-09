// Package internal provides a simple error handling package
package internal

import (
	"runtime"
)

const maxStackDepth = 10

var _ error = (*MyErr)(nil)

// MyErr is an error with a message, a cause, and a stack trace
type MyErr struct {
	Why
	stackTrace []uintptr
}

// Stack return the stack trace of the error
func (e *MyErr) Stack() []map[string]any {
	if len(e.stackTrace) == 0 {
		if err, ok := e.why.(*MyErr); ok {
			return err.Stack()
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

// WithStack returns a new MyErr with the current stack trace
func (e *What) WithStack() *MyErr {
	stack := make([]uintptr, maxStackDepth)
	length := runtime.Callers(3, stack)
	return &MyErr{
		Why:        Why{what: e},
		stackTrace: stack[:length],
	}
}

// WithStack returns a new ThisErrorWithStack with the current stack trace
func (e *Why) WithStack() *MyErr {
	stack := make([]uintptr, maxStackDepth)
	length := runtime.Callers(3, stack)
	return &MyErr{Why: *e, stackTrace: stack[:length]}
}
