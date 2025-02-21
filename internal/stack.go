// Package internal provides a simple error handling package
package internal

import (
	"runtime"
)

const maxStackDepth = 10

var _ error = (*ThisErrorWithStack)(nil)

// ThisErrorWithStack is an error with a message, a cause, and a stack trace
type ThisErrorWithStack struct {
	ThisError
	stackTrace []uintptr
}

// Stack is
func (e *ThisErrorWithStack) Stack() []map[string]any {
	if len(e.stackTrace) == 0 {
		if err, ok := e.cause.(*ThisErrorWithStack); ok {
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
