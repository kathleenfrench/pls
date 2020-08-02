package config

import (
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/viper"
)

// Manager is an interface for managing configs
type Manager interface {
	Get(key string)
	Set(key string, value string)
}

// Settings represent the default settings for pls
type Settings struct {
	GithubToken    string
	GithubUsername string
	Name           string
}

// Get fetches a config value by key
func Get(key string) interface{} {
	return viper.Get(key)
}

// Set sets a config key and value and saves it to the config file
func Set(key string, value string) {
	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		utils.ExitWithError(err)
	}
}
