package cfg

import (
	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
)

func (c *Config) LoadTOMLString(s string) error {
	return c.LoadTOML([]byte(s))
}

func (c *Config) LoadTOML(data []byte) error {
	res := make(map[string]interface{})
	if err := toml.Unmarshal(data, &res); err != nil {
		return err
	}
	return mergo.Merge(&c.Data, res, mergo.WithOverride)
}
