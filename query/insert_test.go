package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	i := NewInsert().
		Into("users").
		Columns("name", "age").
		Row("Bob", 20).
		Row("John", 24)

	stmt, args, err := i.Build(DialectKindCassandra)
	assert.NoError(t, err)

	assert.Equal(t, "INSERT INTO users(name, age) VALUES (?, ?), (?, ?)", stmt)
	assert.Equal(t, "Bob", args[0])
	assert.Equal(t, 20, args[1])
	assert.Equal(t, "John", args[2])
	assert.Equal(t, 24, args[3])
}
