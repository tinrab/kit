package cfg

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"reflect"
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

func TestConfig(t *testing.T) {
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

func TestJSONConfig(t *testing.T) {
	c := New()
	assert.NoError(t, c.LoadFile("./0_default.yml"))
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
			Port:    4200,
			Timeout: 5 * time.Second,
			IDs:     []uint64{1, 2, 42},
			Factor:  3.14,
		},
	}
	assert.True(t, reflect.DeepEqual(expected, d))
}
