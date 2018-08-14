package cfg

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"reflect"
)

type testConfig struct {
	A        testAConfig   `cfg:"a"`
	D        string        `cfg:"d"`
	E        float64       `cfg:"e"`
	LiveTime time.Duration `cfg:"live_time"`
}

type testAConfig struct {
	B uint32 `cfg:"b"`
	C []int  `cfg:"c"`
}

func TestConfig(t *testing.T) {
	c := NewConfig()
	assert.NoError(t, c.LoadYAMLString(`
a:
  b: 42
  c: [1, 2, 3]
d: "Hi"
`))
	assert.NoError(t, c.LoadYAMLString(`
a:
  b: 13
e: 3.14
live_time: 5s
`))

	d := testConfig{}
	assert.NoError(t, c.Decode(&d))

	expected := testConfig{
		A: testAConfig{
			B: 13,
			C: []int{1, 2, 3},
		},
		D:        "Hi",
		E:        3.14,
		LiveTime: 5 * time.Second,
	}
	assert.True(t, reflect.DeepEqual(expected, d))
}
