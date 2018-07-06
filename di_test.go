package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type Database interface {
	Dependency
	GetValue() int
}

type databaseImpl struct {
	Started bool
}

func (d *databaseImpl) GetValue() int {
	return 42
}

func (d *databaseImpl) Init() error {
	d.Started = true
	return nil
}

type service struct {
	Database Database `inject:"db"`
	Started  bool
}

func (s *service) Init() error {
	s.Started = true
	return nil
}

func TestDependencyInjection(t *testing.T) {
	di := NewDependencyInjection()
	di.Provide("db", &databaseImpl{})
	di.Provide("s", &service{})

	assert.NoError(t, di.Resolve())

	db := di.Get("db").(*databaseImpl)
	assert.True(t, db.Started)
	s := di.Get("s").(*service)
	assert.True(t, s.Started)
	assert.Equal(t, 42, s.Database.GetValue())
}
