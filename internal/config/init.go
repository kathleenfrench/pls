package config

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/kathleenfrench/pls/pkg/web/git"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const notApplicable = "NA"

// default keys
const (
	githubUsernameKey  = "git_username"
	githubTokenKey     = "git_token"
	nameKey            = "name"
	defaultEditorKey   = "default_editor"
	webShortcutsKey    = "webshort"
	defaultCodepathKey = "default_codepath"

	githubEnterpriseUsernameKey = "ghe_username"
	githubEnterpriseTokenKey    = "ghe_token"
	githubEnterpriseHostKey     = "ghe_hostname"
)

func unset(val interface{}) bool {
	if val == nil || val == "" {
		return true
	}

	return false
}

// checkForUnsetRequiredDefaults prompts the user for unset defaults and sets them
func checkForUnsetRequiredDefaults() bool {
	var unsetFound bool

	if unset(viper.Get(githubUsernameKey)) {
		// check if we can find it first using the git pkg
		usernameCheck, err := git.CheckForGitUsername()
		if err == nil {
			confirmUsername := gui.ConfirmPrompt(fmt.Sprintf("is %s the correct username?", usernameCheck), "", true, true)
			if confirmUsername {
				viper.Set(githubUsernameKey, usernameCheck)
			} else {
				correctUsername := gui.InputPromptWithResponse("input the correct username", "", true)
				viper.Set(githubUsernameKey, correctUsername)
			}
		} else {
			unsetFound = true
			gu := gui.InputPromptWithResponse("what is your github username?", "", true)
			viper.Set(githubUsernameKey, gu)
		}
	}

	if unset(viper.Get(githubTokenKey)) {
		unsetFound = true
		gt := gui.InputPromptWithResponse("what is your github token?", "", true)
		viper.Set(githubTokenKey, gt)
	}

	if unset(viper.Get(nameKey)) {
		whoami, _ := utils.BashExec("whoami")
		unsetFound = true
		nm := gui.InputPromptWithResponse("what's your name?", whoami, true)
		viper.Set(nameKey, nm)
	}

	if unset(viper.Get(defaultEditorKey)) {
		unsetFound = true
		de := promptForDefaultEditor()
		viper.Set(defaultEditorKey, de)
	}

	if unset(viper.Get(defaultCodepathKey)) {
		unsetFound = true
		home, err := homedir.Dir()
		if err != nil {
			utils.ExitWithError(err)
		}

		home = fmt.Sprintf("%s/", home)
		codePath := gui.InputPromptWithResponse(fmt.Sprintf("what is the relative path from %s you want repos cloned?", home), "", true)
		codePath = fmt.Sprintf("%s%s", home, codePath)
		color.Red(fmt.Sprintf("codepath: %s", codePath))
		viper.Set(defaultCodepathKey, codePath)
	}

	if viper.Get(webShortcutsKey) == nil {
		unsetFound = true
		viper.Set(webShortcutsKey, defaultWebShortcuts)
	}

	if viper.Get(githubEnterpriseHostKey) != notApplicable {
		unsetFound = true
		useGitEnterprise := gui.ConfirmPrompt("do you want to configure pls to work with github enterprise?", "", false, true)
		if useGitEnterprise {
			host := gui.InputPromptWithResponse("what is the hostname of your git enterprise?", "github.[company].com", true)
			viper.Set(githubEnterpriseHostKey, host)
			gheUser := gui.InputPromptWithResponse("what is your git enterprise username?", "", true)
			viper.Set(githubEnterpriseUsernameKey, gheUser)
			gheToken := gui.InputPromptWithResponse("what is your git enterprise token?", "", true)
			viper.Set(githubEnterpriseTokenKey, gheToken)
		} else {
			viper.Set(githubEnterpriseHostKey, notApplicable)
			viper.Set(githubEnterpriseUsernameKey, notApplicable)
			viper.Set(githubEnterpriseTokenKey, notApplicable)
		}
	}

	return unsetFound
}

var defaultWebShortcuts = map[string]string{
	"git":   "https://github.com/",
	"gmail": "https://mail.google.com/mail/u/0/#inbox",
}

// Initialize creates the directory and/or file with defaults for the application's configuration settings
func Initialize() {
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
		utils.PrintError(fmt.Sprintf("ReadInConfig: %s", err))
		err = viper.WriteConfig()
		if err != nil {
			utils.PrintError(fmt.Sprintf("WriteConfig: %s", err))
			utils.ExitWithError(err)
		}
	}

	// set defaults
	unsetValuesFound := checkForUnsetRequiredDefaults()
	if unsetValuesFound {
		err = viper.WriteConfig()
		if err != nil {
			utils.ExitWithError(err)
		}
	}

	viper.WatchConfig()
}
