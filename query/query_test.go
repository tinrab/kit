package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	stmt, args, err := NewQuery().
		From("users").
		Columns("name", "created_at", "age").
		Where(Equal("name", "Bob")).
		Take(1).
		Build(DialectKindCassandra)
	assert.NoError(t, err)

	assert.Equal(t, "SELECT name, created_at, age FROM users WHERE name = ? LIMIT 1 ALLOW FILTERING", stmt)
	assert.Equal(t, "Bob", args[0])
}

func TestInQuery(t *testing.T) {
	stmt, args, err := NewQuery().
		From("users").
		Columns("name").
		Where(In("id", []interface{}{10, 20, 50})).
		Build(DialectKindCassandra)
	assert.NoError(t, err)

	assert.Equal(t, "SELECT name FROM users WHERE id IN (?, ?, ?) LIMIT 0 ALLOW FILTERING", stmt)
	assert.EqualValues(t, []interface{}{10, 20, 50}, args)
}
