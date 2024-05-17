package config

import "github.com/tidwall/gjson"

var config = New()

func Get(field string) gjson.Result {
	return config.Get(field)
}

func GetString(field string) string {
	return config.Get(field).String()
}

func GetBool(field string) bool {
	return config.Get(field).Bool()
}

func GetUint64(field string) uint64 {
	return config.Get(field).Uint()
}

func SetProfile(profile string) {
	config.SetProfile(profile)
}
