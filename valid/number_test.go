package valid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberRangeFloat(t *testing.T) {
	v := New(
		Range(2, 5),
	)
	assert.NoError(t, v.Validate(3))
	assert.Equal(t, ErrNumberOutOfRange, v.Validate(6))
}
