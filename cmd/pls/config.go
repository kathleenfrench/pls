package pls

import (
	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		// Search config in home directory with name ".pls" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pls")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
