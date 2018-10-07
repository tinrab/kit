package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Database interface {
	GetValue() int
}

type databaseImpl struct {
	Started bool
}

func (d *databaseImpl) GetValue() int {
	return 42
}

func (d *databaseImpl) Open() error {
	d.Started = true
	return nil
}

func (d *databaseImpl) Close() {
}

type service struct {
	Database Database `inject:"db"`
	Started  bool
}

func (s *service) Open() error {
	s.Started = true
	return nil
}

func (s *service) Close() {
}

func TestDependencyInjection(t *testing.T) {
	c := New()
	c.Provide("db", &databaseImpl{})
	c.Provide("s", &service{})

	assert.NoError(t, c.Resolve())

	db, _ := c.GetByName("db").(*databaseImpl)
	assert.True(t, db.Started)
	s, _ := c.GetByName("s").(*service)
	assert.True(t, s.Started)
	assert.Equal(t, 42, s.Database.GetValue())
}

type A struct {
	B *B
}

func (a *A) Value() int {
	return a.B.Value()
}

type B struct {
}

func (b *B) Value() int {
	return 42
}

func TestWithStruct(t *testing.T) {
	c := New()
	c.Provide("a", &A{})
	c.Provide("b", &B{})
	assert.NoError(t, c.Resolve())

	a := c.GetByType(&A{}).(*A)
	assert.NotNil(t, a.B)
	assert.Equal(t, 42, a.Value())
}
