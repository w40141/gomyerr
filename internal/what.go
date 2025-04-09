// Package internal provides a simple error handling package
package internal

import (
	"fmt"
)

var _ error = (*What)(nil)

// What is an error with a message
type What struct {
	msg string
}

func (e *What) Error() string {
	return e.msg
}

// NewWhat returns a new why error
func NewWhat(msg string) *What {
	return &What{msg: msg}
}

// FormatWhat returns a new why error with a formatted message
func FormatWhat(format string, a ...any) *What {
	return NewWhat(fmt.Sprintf(format, a...))
}
