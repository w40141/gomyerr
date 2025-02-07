package gomyerr

import (
	"fmt"
	"log/slog"
	"runtime"
)

// MyError はエラーの構造体です。
type MyError struct {
	message string
	stack   []uintptr
	err     error
}

var hasStack = true

const maxDepth = 50

// NewMyError はエラーを生成します。
func NewMyError(msg string) *MyError {
	if !hasStack {
		return &MyError{message: msg}
	}
	stack := make([]uintptr, maxDepth)
	length := runtime.Callers(2, stack)
	return &MyError{message: msg, stack: stack[:length]}
}

// Wrap はエラーをラップします。
func (e *MyError) Wrap(err error) *MyError {
	if myErr, ok := e.err.(*MyError); ok {
		newErr := NewMyError(myErr.Error())
		newErr.err = err
		e.err = newErr
		return e
	}
	e.err = err
	return e
}

// Format はエラーをラップします。
func Format(msg string, errs ...error) *MyError {
	if len(errs) == 0 {
		return NewMyError(msg)
	}

	head, tail := errs[0], errs[1:]
	if len(tail) == 0 {
		return NewMyError(msg).Wrap(head)
	}
	return Format(msg, tail...)
}

// Log はエラーログを出力します。
func (e *MyError) Log(instance string) {
	slog.Error(
		e.Error(),
		"detail",
		e.Cause(),
		"instance",
		instance,
		"stacktrace",
		e.Stack(),
		"logName",
		"almex-abd-dev/abd-functions/post-error-log",
	)
}

// Error はエラーメッセージを取得します。
func (e *MyError) Error() string {
	return e.message
}

// Unwrap はエラーを取得します。
func (e *MyError) Unwrap() error {
	return e.err
}

// Cause はエラーの原因を取得します。
func (e *MyError) Cause() string {
	if e.err == nil {
		return e.message
	}
	if myErr, ok := e.err.(*MyError); ok {
		return fmt.Sprintf("%s: %s", e.message, myErr.Cause())
	}
	return fmt.Sprintf("%s: %s", e.message, e.err.Error())
}

// Stack はスタックトレースを取得します。
func (e *MyError) Stack() []map[string]any {
	if !hasStack {
		return nil
	}

	if e.stack == nil {
		if myErr, ok := e.err.(*MyError); ok {
			return myErr.Stack()
		}
		return nil
	}

	frames := runtime.CallersFrames(e.stack)
	stack := make([]map[string]any, 0, maxDepth)
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

// Replacer はログの属性を置換します。
func Replacer(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.MessageKey:
		a.Key = "message"
	case slog.LevelKey:
		a.Key = "severity"
	case slog.SourceKey:
		a.Key = "logging.googleapis.com/sourceLocation"
	default:
	}

	return a
}
