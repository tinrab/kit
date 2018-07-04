package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependencyInjection(t *testing.T) {
	di := NewDependencyInjection()
	di.Provide("x", 42)
	di.ProvideWith("y", func(di *DependencyInjection) interface{} {
		return di.Get("x").(int) * 3
	})

	assert.Equal(t, 42, di.Get("x").(int))
	assert.Equal(t, 42*3, di.Get("y").(int))
}
