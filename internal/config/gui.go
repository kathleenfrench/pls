package config

import (
	"fmt"

	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
)

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
		case "default_codepath":
			m["Default Code Path"] = defaultCodepathKey
		case "useviper":
			// exclude from dropdown
			break
		default:
			if isWebShortcutKey(k) {
				m["Web Shortcuts"] = webShortcutsKey
			} else {
				m[k] = k
			}
		}
	}

	uiKeys := utils.GetKeysFromMapString(m)
	return uiKeys, m
}

func promptForDefaultEditor() string {
	options := []string{"vim", "emacs", "vscode", "atom", "sublime", "phpstorm"}
	selection := gui.SelectPromptWithResponse("which would you like to set as your default editor?", options, false)

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

func addNewWebShortcut() (string, string) {
	target := gui.InputPromptWithResponse("what do you want to name the shortcut?", "", false)

	url := gui.InputPromptWithResponse(fmt.Sprintf("what is the shortcut url you want to set for %s?", target), "", false)

	if target == "" || url == "" {
		utils.ExitWithError("missing required values")
	}

	return fmt.Sprintf("webshort.%s", target), url
}
