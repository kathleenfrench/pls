package config

import (
	"path/filepath"
	"strings"
)

// UpdateSettings checks for pls config values that have already been set and ensures they're preserved when updating configs
func (s *Settings) UpdateSettings() error {
	cfgFile := s.viper.ConfigFileUsed()

	if s.GitToken != "" {
		s.viper.Set(githubTokenKey, strings.TrimSpace(s.GitToken))
	}

	if s.GitUsername != "" {
		s.viper.Set(githubUsernameKey, strings.TrimSpace(s.GitUsername))
	}

	if s.Name != "" {
		s.viper.Set(nameKey, strings.TrimSpace(s.Name))
	}

	s.viper.MergeInConfig()
	s.viper.SetConfigFile(cfgFile)
	s.viper.SetConfigType(filepath.Ext(cfgFile))
	err := s.viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}
