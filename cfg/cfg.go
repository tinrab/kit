package cfg

import (
	"io/ioutil"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
)

type Config struct {
	Data map[string]interface{}
}

func New() *Config {
	return &Config{
		Data: make(map[string]interface{}),
	}
}

func (c *Config) LoadGlob(pattern string) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, filename := range files {
		err = c.LoadFile(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) LoadFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	ext := filepath.Ext(filename)
	if ext == ".yml" || ext == ".yaml" {
		return c.LoadYAML(data)
	} else if ext == ".json" {
		return c.LoadJSON(data)
	}
	return nil
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
