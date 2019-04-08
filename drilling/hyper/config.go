package main

import "github.com/spf13/viper"

func GetString(key, defaultValue string) (value string) {
	if value = viper.GetString(key); value == "" {
		value = defaultValue
	}
	return
}
