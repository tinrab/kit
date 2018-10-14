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

func TestNumberRangeFloat(t *testing.T) {
	v := New(
		Range(2, 5),
	)
	assert.NoError(t, v.Validate(3))
	assert.Equal(t, ErrNumberOutOfRange, v.Validate(6))
}

func TestTextRegexp(t *testing.T) {
	v := New(
		Regexp("^[a-z]+$"),
	)
	assert.NoError(t, v.Validate("abc"))
	assert.Error(t, ErrInvalidTextPattern, v.Validate("Aabc"))
}

func TestTextMaxLength(t *testing.T) {
	v := New(
		MaxLength(4),
	)
	assert.NoError(t, v.Validate("abc"))
	assert.Error(t, ErrTextTooLong, v.Validate("abcde"))
}

func TestTextMinLength(t *testing.T) {
	v := New(
		MinLength(4),
	)
	assert.NoError(t, v.Validate("abcde"))
	assert.Error(t, ErrTextTooShort, v.Validate("abc"))
}

func TestTextEmail(t *testing.T) {
	v := New(
		Email(),
	)

	assert.NoError(t, v.Validate("user@example.com"))
	assert.Error(t, ErrInvalidEmail, v.Validate("a"))
	assert.Error(t, ErrInvalidEmail, v.Validate("example1234@"))
	assert.Error(t, ErrInvalidEmail, v.Validate("@example1234.com"))
	assert.Error(t, ErrInvalidEmail, v.Validate("user@example"))
}
