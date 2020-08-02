package config

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
)

// UpdatePrompt lets the user select from a dropdown of config keys for which value to update
func UpdatePrompt(viperSettings map[string]interface{}) error {
	keys := utils.GetKeysFromMapStringInterface(viperSettings)
	uiKeys, uiKeyMap := genGuiKeyMap(keys)

	choice := gui.SelectPromptWithResponse("which config value do you want to change?", uiKeys)
	choiceKey := uiKeyMap[choice]

	color.HiBlue(fmt.Sprintf("wants to change: %s", choiceKey))

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
		case "userviper":
			// don't let users change this one
			break
		default:
			m[k] = k
		}
	}

	uiKeys := utils.GetKeysFromMapString(m)
	return uiKeys, m
}
