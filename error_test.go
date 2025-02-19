package terror

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		msg string
	}{
		"正常系": {
			msg: "test",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := New(test.msg)
			if actual.Error() != test.msg {
				t.Errorf("%v: Error() is expected %s, got %s", name, test.msg, actual.Error())
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := map[string]struct {
		msg     string
		args    []any
		wantMsg string
	}{
		"正常系: 引数が0つ": {
			msg:     "test",
			args:    []any{},
			wantMsg: "test",
		},
		"正常系: 引数が1つ": {
			msg:     "test: %s",
			args:    []any{"arg1"},
			wantMsg: "test: arg1",
		},
		"正常系: 引数が2つ": {
			msg:     "test: %s, %d",
			args:    []any{"arg1", 10},
			wantMsg: "test: arg1, 10",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := Format(test.msg, test.args...)
			if actual.Error() != test.wantMsg {
				t.Errorf("%v: Error() is expected %s, got %s", name, test.wantMsg, actual.Error())
			}
		})
	}
}

func TestWhyIs(t *testing.T) {
	tests := map[string]struct {
		origin error
		target error
		want   bool
	}{
		"正常系": {
			origin: New("test"),
			target: New("test"),
			want:   true,
		},
		"異なるエラーは同じにならない": {
			origin: New("test"),
			target: New("test2"),
			want:   false,
		},
		"異なるエラー型は同じにならない": {
			origin: New("test"),
			target: errors.New("test"),
			want:   false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := errors.Is(test.origin, test.target)
			if actual != test.want {
				t.Errorf("%v: Is() is expected %t, got %t", name, test.want, actual)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	tests := map[string]struct {
		origin  error
		cause   error
		wantMsg string
		equal   bool
	}{
		"正常系": {
			origin:  errors.New("test"),
			cause:   New("cause"),
			wantMsg: "test",
			equal:   true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := Wrap(test.cause, test.origin)
			if actual.Error() != test.wantMsg {
				t.Errorf("%v: Error() is expected %s, got %s", name, test.wantMsg, actual.Error())
			}
			unwrapped := errors.Unwrap(actual)
			if !errors.Is(unwrapped, test.cause) {
				t.Errorf("%v: Unwrap() is expected %s, got %s", name, test.cause, unwrapped)
			}
		})
	}
}
