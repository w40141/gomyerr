package internal

import (
	"errors"
	"testing"
)

func TestNewErr(t *testing.T) {
	msg := "test"
	err := NewErr(msg)
	if err.Error() != msg {
		t.Errorf("expected %s, got %s", msg, err.Error())
	}
}

func TestFormatErr(t *testing.T) {
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
			actual := FormatErr(tt.format, tt.args...)
			if actual.Error() != tt.want {
				t.Errorf("expected %s, got %s", tt.want, actual.Error())
			}
		})
	}
}

func TestErrIs(t *testing.T) {
	origin := NewErr("origin")
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
			why:    NewErr("origin"),
			target: NewErr("origin"),
			want:   false,
		},
		"異なるメッセージは同一でない": {
			why:    NewErr("origin"),
			target: NewErr("cause"),
			want:   false,
		},
		"異なるエラーは同一でない": {
			why:    NewErr("origin"),
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

// TODO: implement
func TestWithStack(_ *testing.T) {
}
