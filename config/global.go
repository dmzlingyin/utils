package config

import "github.com/tidwall/gjson"

var config = New()

func Get(field string) gjson.Result {
	return config.Get(field)
}

func GetString(field string) string {
	return config.Get(field).String()
}

func SetProfile(profile string) {
	config.SetProfile(profile)
}
