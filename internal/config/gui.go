package config

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/viper"
)

// UpdatePrompt lets the user select from a dropdown of config keys for which value to update
func UpdatePrompt(viperSettings map[string]interface{}) error {
	var changedValue string

	keys := viper.AllKeys()
	uiKeys, uiKeyMap := genGuiKeyMap(keys)
	choice := gui.SelectPromptWithResponse("which config value do you want to change?", uiKeys)
	choiceKey := uiKeyMap[choice]
	color.HiYellow(fmt.Sprintf("current value: %v", viperSettings[choiceKey]))

	if choiceKey == defaultEditorKey {
		changedValue = promptForDefaultEditor()
	} else {
		changedValue = gui.InputPromptWithResponse(fmt.Sprintf("what do you want to change %s to?", choice), "")
	}

	v := viper.GetViper()
	v.Set(choiceKey, changedValue)
	parsed, err := Parse(v)
	if err != nil {
		utils.ExitWithError(err)
	}

	parsed.UpdateSettings()
	color.HiGreen(fmt.Sprintf("successfully updated %s to equal %s!", choice, changedValue))
	return nil
}

func genGuiKeyMap(keys []string) ([]string, map[string]string) {
	m := make(map[string]string)

	for _, k := range keys {
		switch k {
		case githubTokenKey:
			m["Github Token"] = githubTokenKey
		case githubUsernameKey:
			m["Github Username"] = githubUsernameKey
		case nameKey:
			m["Name"] = nameKey
		case defaultEditorKey:
			m["Default Editor"] = defaultEditorKey
		case "useviper":
			// exclude from dropdown
			break
		default:
			m[k] = k
		}
	}

	uiKeys := utils.GetKeysFromMapString(m)
	return uiKeys, m
}

func promptForDefaultEditor() string {
	options := []string{"vim", "emacs", "vscode", "atom", "sublime", "phpstorm"}
	selection := gui.SelectPromptWithResponse("which would you like to set as your default editor?", options)

	switch selection {
	case "phpstorm":
		// check for which executable is installed
		p, err := utils.GetPHPStormExecutable()
		if err != nil {
			utils.ExitWithError(err)
		}

		return p
	default:
		return selection
	}
}
