package internal

import (
	"errors"
	"testing"
)

func TestNewWhy(t *testing.T) {
	msg := "test"
	err := NewWhy(msg)
	if err.Error() != msg {
		t.Errorf("expected %s, got %s", msg, err.Error())
	}
}

func TestFormatWhy(t *testing.T) {
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
			actual := FormatWhy(tt.format, tt.args...)
			if actual.Error() != tt.want {
				t.Errorf("expected %s, got %s", tt.want, actual.Error())
			}
		})
	}
}

func TestWithCause(t *testing.T) {
	tests := map[string]struct {
		why   *Why
		cause error
		want  string
	}{
		"正常系": {
			why:   NewWhy("origin"),
			cause: NewWhy("cause"),
			want:  "cause",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.why.WithCause(tt.cause)
			if actual.cause.Error() != tt.want {
				t.Errorf("expected %s, got %s", tt.want, actual.Error())
			}
		})
	}
}

func TestWhyIs(t *testing.T) {
	origin := NewWhy("origin")
	tests := map[string]struct {
		why    error
		target error
		want   bool
	}{
		"同じインスタンスは同一になる": {
			why:    origin,
			target: origin,
			want:   true,
		},
		"異なるインタンスは同一ではない": {
			why:    NewWhy("origin"),
			target: NewWhy("origin"),
			want:   false,
		},
		"異なるメッセージは同一でない": {
			why:    NewWhy("origin"),
			target: NewWhy("cause"),
			want:   false,
		},
		"異なるエラーは同一でない": {
			why:    NewWhy("origin"),
			target: errors.New("cause"),
			want:   false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if errors.Is(tt.why, tt.target) != tt.want {
				t.Errorf("expected %t, got %t", tt.want, errors.Is(tt.why, tt.target))
			}
		})
	}
}
