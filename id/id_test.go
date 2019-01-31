package id

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Cat struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
}

func TestParse(t *testing.T) {
	i, err := Parse("6604873748002701312")
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

func TestGenerate(t *testing.T) {
	g := NewGenerator(42)

	var ids []ID

	for i := 0; i < 10; i++ {
		id := g.Generate()
		for _, x := range ids {
			if x == id {
				t.Fail()
			}
		}
		ids = append(ids, id)
	}
}

func TestGenerateList(t *testing.T) {
	g := NewGenerator(42)
	const n = 1026

	ids := g.GenerateList(n)

	assert.Equal(t, n, len(ids))

	for i, x := range ids {
		if i == 0 {
			continue
		}

		for _, y := range ids[:i] {
			if y == x {
				t.Fatal("generated duplicate IDs")
			}
		}
	}
}
