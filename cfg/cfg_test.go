package cfg

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type config struct {
	Name    string        `cfg:"name"`
	Service serviceConfig `cfg:"service"`
}

type serviceConfig struct {
	Port    uint16        `cfg:"port"`
	Timeout time.Duration `cfg:"timeout"`
	IDs     []uint64      `cfg:"ids"`
	Factor  float32       `cfg:"factor"`
}

func TestGlob(t *testing.T) {
	c := New()
	assert.NoError(t, c.LoadGlob("./*"))

	d := config{}
	assert.NoError(t, c.Decode(&d))
	expected := config{
		Name: "App",
		Service: serviceConfig{
			Port:    4200,
			Timeout: 5 * time.Second,
			IDs:     []uint64{1, 2, 42},
			Factor:  3.14,
		},
	}
	assert.True(t, reflect.DeepEqual(expected, d))
}

func TestJSON(t *testing.T) {
	c := New()
	assert.NoError(t, c.LoadJSONString(`{
		"name": "App",
		"service": {
			"port": 4200,
			"ids": [1, 2, 42],
			"factor": 3.14
		}
	}`))

	d := config{}
	assert.NoError(t, c.Decode(&d))
	expected := config{
		Name: "App",
		Service: serviceConfig{
			Port:   4200,
			IDs:    []uint64{1, 2, 42},
			Factor: 3.14,
		},
	}
	assert.True(t, reflect.DeepEqual(expected, d))
}

func TestYAML(t *testing.T) {
	c := New()
	assert.NoError(t, c.LoadYAMLString(`
name: "App"
service:
  port: 4200
  ids: [1, 2, 42]
  factor: 3.14
`))

	d := config{}
	assert.NoError(t, c.Decode(&d))
	expected := config{
		Name: "App",
		Service: serviceConfig{
			Port:   4200,
			IDs:    []uint64{1, 2, 42},
			Factor: 3.14,
		},
	}
	assert.True(t, reflect.DeepEqual(expected, d))
}
