package security

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScrypt(t *testing.T) {
	data := []byte("1234")
	hash, err := ScryptToBase64(data, 1<<15, 8, 1, 32)
	assert.NoError(t, err)

	assert.True(t, ScryptEqualsBase64([]byte("1234"), hash))
	assert.False(t, ScryptEqualsBase64([]byte("12345"), hash))
}
