package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := NewError("A")
	assert.Equal(t, err, err.Throw())
}
