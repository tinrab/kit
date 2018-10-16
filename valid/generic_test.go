package valid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	v := New(
		Equal(42),
	)
	assert.NoError(t, v.Validate(42))
	assert.Error(t, ErrNotEqual, v.Validate(42.0))
}

func TestNil(t *testing.T) {
	v := New(
		Nil(),
	)
	assert.NoError(t, v.Validate(nil))
	assert.Error(t, ErrNotNil, v.Validate(42))
}

func TestNotNil(t *testing.T) {
	v := New(
		NotNil(),
	)
	assert.NoError(t, v.Validate(42))
	assert.Error(t, ErrNil, v.Validate(nil))
}
