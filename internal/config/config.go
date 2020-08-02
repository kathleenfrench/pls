package config

import (
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Settings represent the default settings for pls
type Settings struct {
	viper         *viper.Viper
	GitToken      string            `yaml:"git_token"`
	GitUsername   string            `yaml:"git_username"`
	Name          string            `yaml:"name"`
	DefaultEditor string            `yaml:"default_editor"`
	WebShortcuts  map[string]string `yaml:"webshort"`
}

func decodeWithYaml(tagName string) viper.DecoderConfigOption {
	return func(c *mapstructure.DecoderConfig) {
		c.TagName = tagName
	}
}

// ParseAndUpdate parses the viper settings as a pls settings struct and updates the config file
func ParseAndUpdate(v *viper.Viper) error {
	s, err := Parse(v)
	if err != nil {
		return err
	}

	return s.UpdateSettings()
}

// Parse unmarshals the viper configs into the pls settings struct
func Parse(v *viper.Viper) (Settings, error) {
	s := Settings{
		viper: v,
	}

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
