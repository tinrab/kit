package cfg

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/mitchellh/mapstructure"
	"github.com/imdario/mergo"
	"path/filepath"
)

type Config struct {
	Data map[string]interface{}
}

func NewConfig() *Config {
	return &Config{
		Data: make(map[string]interface{}),
	}
}

func (c *Config) LoadYAMLFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return c.LoadYAML(data)
}

func (c *Config) LoadYAMLGlob(pattern string) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, filename := range files {
		ext := filepath.Ext(filename)
		if ext == ".yml" || ext == ".yaml" {
			c.LoadYAMLFile(filename)
		}
	}
	return nil
}

func (c *Config) LoadYAMLString(s string) error {
	return c.LoadYAML([]byte(s))
}

func (c *Config) LoadYAML(data []byte) error {
	res := make(map[string]interface{})
	if err := yaml.Unmarshal(data, &res); err != nil {
		return err
	}
	return mergo.Merge(&c.Data, res, mergo.WithOverride)
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
