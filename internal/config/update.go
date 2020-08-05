package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
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

	if s.DefaultEditor != "" {
		s.viper.Set(defaultEditorKey, strings.TrimSpace(s.DefaultEditor))
	}

	if s.WebShortcuts != nil {
		s.viper.Set(webShortcutsKey, s.WebShortcuts)
	}

	if s.DefaultCodeDir != "" {
		s.viper.Set(defaultCodepathKey, s.DefaultCodeDir)
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

// UpdatePrompt lets the user select from a dropdown of config keys for which value to update
func UpdatePrompt(viperSettings map[string]interface{}) error {
	var changedValue string

	v := viper.GetViper()
	keys := viper.AllKeys()
	uiKeys, uiKeyMap := genGuiKeyMap(keys)
	choice := gui.SelectPromptWithResponse("which config value do you want to change?", uiKeys, false)
	choiceKey := uiKeyMap[choice]

	switch choiceKey {
	case defaultEditorKey:
		changedValue = promptForDefaultEditor()
		v.Set(choiceKey, changedValue)
	case webShortcutsKey:
		shorts := make(map[string]string)
		err := viper.UnmarshalKey(webShortcutsKey, &shorts)
		if err != nil {
			utils.ExitWithError(err)
		}

		// check if they want to edit an existing value
		editExisting := gui.ConfirmPrompt("do you want to modify an existing url?", "", false, true)
		if editExisting {
			shortKeys := utils.GetKeysFromMapString(shorts)
			editWhich := gui.SelectPromptWithResponse("which do you want to change?", shortKeys, false)
			changedValue = gui.InputPromptWithResponse(fmt.Sprintf("what do you want to change %s to?", editWhich), "", true)
			v.Set(fmt.Sprintf("webshort.%s", editWhich), changedValue)
		} else {
			target, url := addNewWebShortcut()
			v.Set(target, url)
			changedValue = url
		}
	case defaultCodepathKey:
		rel := gui.InputPromptWithResponse(fmt.Sprintf("what do you want to change %s to?", choice), "", true)
		home, err := homedir.Dir()
		if err != nil {
			utils.ExitWithError(err)
		}

		changedValue := fmt.Sprintf("%s/%s", home, rel)
		confirmChange := gui.ConfirmPrompt(fmt.Sprintf("is %s the correct path?", changedValue), "", true, true)
		if !confirmChange {
			retry := gui.InputPromptWithResponse(fmt.Sprintf("what do you want to change %s to?", choice), "", true)
			changedValue = fmt.Sprintf("%s/%s", home, retry)
		}

		v.Set(choiceKey, changedValue)
	default:
		changedValue = gui.InputPromptWithResponse(fmt.Sprintf("what do you want to change %s to?", choice), "", true)
		v.Set(choiceKey, changedValue)
	}

	parsed, err := Parse(v)
	if err != nil {
		utils.ExitWithError(err)
	}

	parsed.UpdateSettings()
	color.HiGreen(fmt.Sprintf("successfully updated %s to equal %s!", choice, changedValue))
	return nil
}

func isWebShortcutKey(key string) bool {
	return strings.Contains(key, "webshort.")
}
