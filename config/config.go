package config

import (
	"github.com/tidwall/gjson"
	"os"
)

func New() *Config {
	return &Config{}
}

type Config struct {
	content string
}

func (c *Config) Get(field string) gjson.Result {
	return gjson.Get(c.content, field)
}

func (c *Config) SetProfile(profile string) {
	f, err := os.ReadFile(profile + ".json")
	if err != nil {
		panic(err)
	}
	c.content = string(f)
}
