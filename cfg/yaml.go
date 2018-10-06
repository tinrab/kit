package cfg

import (
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

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
