package cfg

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/mitchellh/mapstructure"
	"github.com/imdario/mergo"
)

type Config struct {
	Data map[string]interface{}
}

func NewConfig() *Config {
	return &Config{
		Data: make(map[string]interface{}),
	}
}

func (c *Config) LoadYAML(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return c.LoadYAMLString(string(data))
}

func (c *Config) LoadYAMLString(s string) error {
	data := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(s), &data); err != nil {
		return err
	}
	return mergo.Merge(&c.Data, data, mergo.WithOverride)
}

func (c *Config) Decode(out interface{}) error {
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "cfg",
		WeaklyTypedInput: true,
		Result:           out,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
		),
	})
	if err != nil {
		return err
	}
	return d.Decode(c.Data)
}
