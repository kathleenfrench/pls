package config

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/viper"
)

// default keys
const (
	githubUsernameKey = "git_username"
	githubTokenKey    = "git_token"
	nameKey           = "name"
)

// checkForUnsetDefaults prompts the user for unset defaults and sets them
func checkForUnsetDefaults() {
	var (
		gt     string
		gu     string
		nm     string
		prompt *survey.Input
	)

	gitUsername := viper.Get(githubUsernameKey)
	if gitUsername == nil {
		prompt = &survey.Input{
			Message: "what is your github username?",
		}

		survey.AskOne(prompt, &gu)
	}

	gitToken := viper.Get(githubTokenKey)
	if gitToken == nil {
		prompt = &survey.Input{
			Message: "what is your github token?",
		}

		survey.AskOne(prompt, &gt)
		viper.Set(githubTokenKey, gt)
	}

	name := viper.Get(nameKey)
	if name == nil {
		prompt = &survey.Input{
			Message: "what's your name?",
		}

		survey.AskOne(prompt, &nm)
		viper.Set(nameKey, nm)
	}
}

// Initialize creates the directory and/or file with defaults for the application's configuration settings
func Initialize() {
	// set defaults
	checkForUnsetDefaults()

	// set fs properties
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
