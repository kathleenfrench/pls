package config

import (
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Manager is an interface for managing configs
type Manager interface {
	Get(key string)
	Set(key string, value string)
}

// Settings represent the default settings for pls
type Settings struct {
	GitToken    string `yaml:"git_token"`
	GitUsername string `yaml:"git_username"`
	Name        string `yaml:"name"`
}

func decodeWithYaml(tagName string) viper.DecoderConfigOption {
	return func(c *mapstructure.DecoderConfig) {
		c.TagName = tagName
	}
}

// Parse unmarshals the viper configs into the pls settings struct
func Parse(v *viper.Viper) (Settings, error) {
	s := Settings{}

	dco := decodeWithYaml("yaml")
	err := v.Unmarshal(&s, dco)
	if err != nil {
		return s, err
	}

	return s, nil
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
