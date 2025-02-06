// Package gomyerr provides a simple error handling package
// TODO: add stack trace
// TODO: add tags
// TODO: test
package gomyerr

import "errors"

// MyErr is a simple error struct
type MyErr struct {
	origin error
	cause  error
}

type origin struct {
	msg string
}

func (e *origin) Error() string {
	return e.msg
}

// New returns a new MyErr
func New(msg string) *MyErr {
	return &MyErr{origin: &origin{msg: msg}}
}

// From returns a new MyErr with the origin error
func From(err error) *MyErr {
	return &MyErr{origin: err}
}

// Error returns the error message
func (e *MyErr) Error() string {
	return e.origin.Error()
}

// WithCause returns a new MyErr with the cause error
func (e *MyErr) WithCause(cause error) *MyErr {
	if e.cause != nil {
		e.cause = From(cause)
	}
	if _, ok := e.cause.(*MyErr); ok {
		e = e.cause.(*MyErr).WithCause(cause)
	} else {
		e.cause = From(cause)
	}
	return e
}

// Join returns a new MyErr with the origin error and the cause error
func Join(origin error, err ...error) *MyErr {
	e := From(origin)
	if len(err) == 0 {
		return e
	}
	return e.WithCause(Join(err[0], err[1:]...))
}

func (e *MyErr) Unwrap() error {
	return e.cause
}

// Is returns true if the target error is the origin or the cause error
func (e *MyErr) Is(target error) bool {
	if errors.Is(e.origin, target) {
		return true
	}
	if errors.Is(e.cause, target) {
		return true
	}
	return false
}

// As returns true if the target is the origin or the cause error
func (e *MyErr) As(target interface{}) bool {
	if errors.As(e.origin, target) {
		return true
	}
	if errors.As(e.cause, target) {
		return true
	}
	return false
}
