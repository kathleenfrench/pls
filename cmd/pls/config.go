package pls

import (
	"fmt"

	"github.com/kathleenfrench/pls/pkg/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// pls default config file is in $HOME/.config/pls/.pls.yaml

func constructConfigPath(home string) string {
	return fmt.Sprintf("%s/.config/pls", home)
}

func constructConfigFilePath(path string) string {
	return fmt.Sprintf("%s/config.yaml", path)
}

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

		cfgPath := constructConfigPath(home)
		// check for whether the directory and config file already exist
		err = utils.CreateDir(cfgPath)
		if err != nil {
			panic(err)
		}

		cfgFilePath := constructConfigFilePath(cfgPath)
		err = utils.CreateFile(cfgFilePath)
		if err != nil {
			panic(err)
		}

		// Search config in home directory with name ".pls" (without extension).
		viper.AddConfigPath(constructConfigPath(home))
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if Verbose {
			fmt.Println(fmt.Sprintf("using config file located at %s", viper.ConfigFileUsed()))
		}
	}
}
