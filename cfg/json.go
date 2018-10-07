package cfg

import (
	"encoding/json"

	"github.com/imdario/mergo"
)

func (c *Config) LoadJSONString(s string) error {
	return c.LoadJSON([]byte(s))
}

func (c *Config) LoadJSON(data []byte) error {
	res := make(map[string]interface{})
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	return mergo.Merge(&c.Data, res, mergo.WithOverride)
}
