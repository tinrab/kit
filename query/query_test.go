package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery(t *testing.T) {
	q := NewQuery().
		From("users").
		Columns("name", "created_at", "age").
		Where(Equal("name", "Bob")).
		Take(1)
	stmt, args, err := q.Build(DialectCassandra)
	assert.NoError(t, err)

	assert.Equal(t, "SELECT name, created_at, age FROM users WHERE name = ? LIMIT 1 ALLOW FILTERING", stmt)
	assert.Equal(t, "Bob", args[0])
}
