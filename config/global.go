package config

import "github.com/tidwall/gjson"

var config = New()

func Get(field string) gjson.Result {
	return config.Get(field)
}

func SetProfile(profile string) {
	config.SetProfile(profile)
}
