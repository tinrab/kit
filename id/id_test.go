package id

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Cat struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
}

func TestParse(t *testing.T) {
	i, err := ParseID("6604873748002701312")
	assert.NoError(t, err)
	assert.Equal(t, ID(6604873748002701312), i)
}

func TestJSON(t *testing.T) {
	g := NewGenerator(1)

	i := g.Generate()
	cat := Cat{
		ID:   i,
		Name: "John",
	}
	js, err := json.Marshal(cat)
	assert.NoError(t, err)

	decodedCat := &Cat{}
	assert.NoError(t, json.Unmarshal(js, decodedCat))

	assert.Equal(t, cat.Name, decodedCat.Name)
	assert.Equal(t, i, decodedCat.ID)
}
