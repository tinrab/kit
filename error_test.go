package kit

import (
	"testing"
)

var (
	ErrTestError = NewCodeError(42, "test error")
)

func A() error {
	return B()
}

func B() error {
	return ErrTestError.Throw()
}

func TestWrap(t *testing.T) {
	if err := A(); err != nil {
		t.Log(err)
	}
}
