package internal

import (
	"errors"
	"fmt"
	"testing"
)

func TestWithCause(t *testing.T) {
	cause := NewErr("cause")
	tests := map[string]struct {
		why   *Err
		cause error
		want  string
	}{
		"正常系": {
			why:   NewErr("origin"),
			cause: cause,
			want:  "cause",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.why.WithCause(tt.cause)
			if errors.Is(actual, tt.cause) {
				t.Errorf("expected %s, got %s", tt.want, actual.Error())
			}
		})
	}
}

func TestNewMyErr(t *testing.T) {
	origin := NewErr("origin")
	cause := NewErr("cause")
	actual := NewMyErr(origin, cause)
	if errors.Is(actual.why, origin) {
		t.Errorf("expected origin, got %s", actual.Error())
	}
	if errors.Is(actual.cause, cause) {
		t.Errorf("expected cause, got %s", actual.Unwrap().Error())
	}
}

func TestFormatMyErr(t *testing.T) {
	cause := NewErr("cause")
	tests := map[string]struct {
		format string
		args   []any
		want   string
	}{
		"0個": {
			format: "test",
			args:   []any{},
			want:   "test",
		},
		"1個": {
			format: "test %s",
			args:   []any{"test"},
			want:   "test test",
		},
		"2個": {
			format: "test %s %d",
			args:   []any{"test", 1},
			want:   "test test 1",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := FormatMyErr(cause, tt.format, tt.args...)
			if actual.Error() != tt.want {
				t.Errorf("expected %s, got %s", tt.want, actual.Error())
			}
		})
	}
}

func TestMyErrIs(t *testing.T) {
	origin := NewErr("origin")
	origin2 := fmt.Errorf("origin")
	cause := NewErr("cause")
	cause2 := fmt.Errorf("cause")

	tests := map[string]struct {
		err    *MyErr
		target error
		want   bool
	}{
		"whyがイコール": {
			err:    NewMyErr(origin, NewErr("cause")),
			target: origin,
			want:   true,
		},
		"whyがイコール2": {
			err:    NewMyErr(origin2, NewErr("cause")),
			target: origin2,
			want:   true,
		},
		"causeがイコール": {
			err:    NewMyErr(NewErr("origin"), cause),
			target: cause,
			want:   true,
		},
		"causeがイコール2": {
			err:    NewMyErr(NewErr("origin"), cause2),
			target: cause2,
			want:   true,
		},
		"ノットイコール": {
			err:    NewMyErr(NewErr("origin"), NewErr("cause")),
			target: NewErr("not equal"),
			want:   false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if errors.Is(tt.err, tt.target) != tt.want {
				t.Errorf("expected %t, got %t", tt.want, errors.Is(tt.err, tt.target))
			}
		})
	}
}
