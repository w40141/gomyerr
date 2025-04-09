// Package internal provides a simple error handling package
package internal

import (
	"errors"
)

var (
	_ error = (*Why)(nil)
	_ error = (*Whys)(nil)
)

// Why has a original error and a cause error
type Why struct {
	what error
	why  error
}

// Whys has a original error and multiple causes
type Whys struct {
	why  error
	whys []error
}

// Wrap returns a new Why with the cause error
func (e *What) Wrap(c error) *Why {
	return &Why{what: e, why: c}
}

func (e *Why) Error() string {
	return e.what.Error()
}

func (e *Why) Unwrap() error {
	return e.why
}

// Is reports whether the target error is the same as the current error
func (e *Why) Is(target error) bool {
	if errors.Is(e.what, target) || errors.Is(e.why, target) {
		return true
	}
	return false
}

// Join returns a new Whys with the cause error
func (e *Why) Join(err error) *Whys {
	return &Whys{why: e.what, whys: []error{e.why, err}}
}

func (e *Whys) Error() string {
	return e.why.Error()
}

func (e *Whys) Unwrap() []error {
	return e.whys
}

// Join returns a new Whys with the cause errors
func (e *Whys) Join(errs ...error) *Whys {
	e.whys = append(e.whys, errs...)
	return e
}

// NewWhy returns a new Why
func NewWhy(what, why error) *Why {
	return &Why{what: what, why: why}
}

// FormatMyErr returns a new Why with a formatted message
// func FormatMyErr(why error, format string, a ...any) *Why {
// 	return &Why{what: FormatWhat(format, a...), why: why}
// }
