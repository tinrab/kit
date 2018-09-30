package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdate(t *testing.T) {
	stmt, args, err := NewUpdate().
		Table("users").
		Assign("name", "Bob").
		Where(Equal("id", 42)).
		Build(DialectCassandra)

	assert.NoError(t, err)
	assert.Equal(t, "UPDATE users SET name = ? WHERE id = ?", stmt)
	assert.Equal(t, "Bob", args[0])
	assert.Equal(t, 42, args[1])
}
