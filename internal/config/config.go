package config

import (
	"fmt"

	"github.com/kathleenfrench/pls/pkg/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Initialize creates the directory and/or file with defaults for the application's configuration settings
func Initialize() {
	// set defaults
	viper.SetDefault("yeet", "lol")

	viper.AddConfigPath(constructConfigPath())
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.SetConfigFile(fmt.Sprintf("%s/%s.%s", constructConfigPath(), configFileName, configFileType))

	// check for whether the directory and config file already exist
	err := utils.CreateDir(constructConfigPath())
	if err != nil {
		utils.ExitWithError(err)
	}

	err = utils.CreateFile(viper.ConfigFileUsed())
	if err != nil {
		utils.ExitWithError(err)
	}

	viper.AutomaticEnv()

	_ = viper.SafeWriteConfig()
	err = viper.ReadInConfig()
	if err != nil {
		err = viper.WriteConfig()
		if err != nil {
			utils.ExitWithError(err)
		}
	}

	viper.WatchConfig()
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