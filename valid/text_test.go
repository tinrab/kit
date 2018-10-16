package valid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextRegexp(t *testing.T) {
	v := New(
		Regexp("^[a-z]+$"),
	)
	assert.NoError(t, v.Validate("abc"))
	assert.Equal(t, ErrInvalidTextPattern, v.Validate("Aabc"))
}

func TestTextMaxLength(t *testing.T) {
	v := New(
		MaxLength(4),
	)
	assert.NoError(t, v.Validate("abc"))
	assert.Equal(t, ErrTextTooLong, v.Validate("abcde"))
}

func TestTextMinLength(t *testing.T) {
	v := New(
		MinLength(4),
	)
	assert.NoError(t, v.Validate("abcde"))
	assert.Error(t, ErrTextTooShort, v.Validate("abc"))
}

func TestTextMinMaxLength(t *testing.T) {
	v := New(
		MinMaxLength(3, 5),
	)
	assert.NoError(t, v.Validate("123"))
	assert.NoError(t, v.Validate("1234"))
	assert.NoError(t, v.Validate("12345"))
	assert.Equal(t, ErrTextTooShort, v.Validate("12"))
	assert.Equal(t, ErrTextTooLong, v.Validate("123456"))
}

func TestTextHasUpperCase(t *testing.T) {
	v := New(
		HasUpperCase(),
	)

	assert.NoError(t, v.Validate("Abc"))
	assert.Equal(t, ErrTextMissingUpper, v.Validate(""))
	assert.Equal(t, ErrTextMissingUpper, v.Validate("abc"))
}

func TestTextHasLowerCase(t *testing.T) {
	v := New(
		HasLowerCase(),
	)

	assert.NoError(t, v.Validate("ABc"))
	assert.Equal(t, ErrTextMissingLower, v.Validate(""))
	assert.Equal(t, ErrTextMissingLower, v.Validate("ABC"))
}

func TestTextHasSpecial(t *testing.T) {
	v := New(
		HasSpecial(),
	)

	assert.NoError(t, v.Validate("abc!"))
	assert.Equal(t, ErrTextMissingSpecial, v.Validate(""))
	assert.Equal(t, ErrTextMissingSpecial, v.Validate("abc4"))
}

func TestTextHasNumber(t *testing.T) {
	v := New(
		HasNumber(),
	)

	assert.NoError(t, v.Validate("abc42"))
	assert.Equal(t, ErrTextMissingNumber, v.Validate(""))
	assert.Equal(t, ErrTextMissingNumber, v.Validate("abc_"))
}

func TestTextEmail(t *testing.T) {
	v := New(
		Email(),
	)

	assert.NoError(t, v.Validate("user@example.com"))
	assert.Equal(t, ErrInvalidEmail, v.Validate("a"))
	assert.Equal(t, ErrInvalidEmail, v.Validate("example1234@"))
	assert.Equal(t, ErrInvalidEmail, v.Validate("@example1234.com"))
	assert.Equal(t, ErrInvalidEmail, v.Validate("user@example"))
}

func TestTextSlug(t *testing.T) {
	v := New(
		Slug(),
	)

	assert.NoError(t, v.Validate("hello-world"))

	invalidSlugs := []string{
		"",
		"a-",
		"-a",
		"a-$-a",
		"-",
		"a---",
		"A",
	}
	for _, slug := range invalidSlugs {
		assert.Equal(t, ErrTextInvalidSlug, v.Validate(slug))
	}
}
