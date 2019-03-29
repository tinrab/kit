package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSnakeCase(t *testing.T) {
	cases := [][]string{
		{"aB", "a_b"},
		{"a42", "a42"},
		{"AAb", "a_ab"},
		{"A_b", "a_b"},
	}

	for _, c := range cases {
		assert.Equal(t, c[1], ToSnakeCase(c[0]))
	}
}
