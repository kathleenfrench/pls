package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
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
	choice := gui.SelectPromptWithResponse("which config value do you want to change?", uiKeys)
	choiceKey := uiKeyMap[choice]
	color.HiYellow(fmt.Sprintf("current value: %v", viperSettings[choiceKey]))

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
		editExisting := gui.ConfirmPrompt("do you want to modify an existing url?", "", false)
		if editExisting {
			shortKeys := utils.GetKeysFromMapString(shorts)
			editWhich := gui.SelectPromptWithResponse("which do you want to change?", shortKeys)
			changedValue = gui.InputPromptWithResponse(fmt.Sprintf("what do you want to change %s to?", editWhich), "")
			v.Set(fmt.Sprintf("webshort.%s", editWhich), changedValue)
		} else {
			target, url := addNewWebShortcut()
			v.Set(target, url)
		}
	default:
		changedValue = gui.InputPromptWithResponse(fmt.Sprintf("what do you want to change %s to?", choice), "")
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