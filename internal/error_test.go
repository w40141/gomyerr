package internal

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewThisError(t *testing.T) {
	origin := NewWhy("origin")
	cause := NewWhy("cause")
	err := NewThisError(origin, cause)
	if err.Error() != "origin" {
		t.Errorf("expected origin, got %s", err.Error())
	}
	if err.Unwrap().Error() != "cause" {
		t.Errorf("expected cause, got %s", err.Unwrap().Error())
	}
}

func TestFormatThisError(t *testing.T) {
	cause := NewWhy("cause")
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
			actual := FormatThisError(cause, tt.format, tt.args...)
			if actual.Error() != tt.want {
				t.Errorf("expected %s, got %s", tt.want, actual.Error())
			}
		})
	}
}

func TestThisErrorIs(t *testing.T) {
	origin := NewWhy("origin")
	origin2 := fmt.Errorf("origin")
	cause := NewWhy("cause")
	cause2 := fmt.Errorf("cause")

	tests := map[string]struct {
		err    *ThisError
		target error
		want   bool
	}{
		"whyがイコール": {
			err:    NewThisError(origin, NewWhy("cause")),
			target: origin,
			want:   true,
		},
		"whyがイコール2": {
			err:    NewThisError(origin2, NewWhy("cause")),
			target: origin2,
			want:   true,
		},
		"causeがイコール": {
			err:    NewThisError(NewWhy("origin"), cause),
			target: cause,
			want:   true,
		},
		"causeがイコール2": {
			err:    NewThisError(NewWhy("origin"), cause2),
			target: cause2,
			want:   true,
		},
		"ノットイコール": {
			err:    NewThisError(NewWhy("origin"), NewWhy("cause")),
			target: NewWhy("not equal"),
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
