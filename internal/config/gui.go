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
	keys := viper.AllKeys()
	uiKeys, uiKeyMap := genGuiKeyMap(keys)
	choice := gui.SelectPromptWithResponse("which config value do you want to change?", uiKeys)
	choiceKey := uiKeyMap[choice]
	color.HiYellow(fmt.Sprintf("current value: %v", viperSettings[choiceKey]))
	changedValue := gui.InputPromptWithResponse(fmt.Sprintf("what do you want to change %s to?", choice))

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
