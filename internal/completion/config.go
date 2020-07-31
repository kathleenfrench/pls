package completion

import (
	"fmt"

	"github.com/kathleenfrench/pls/pkg/utils"
)

var (
	fishShell = "fish"
	bashShell = "bash"
	zshShell  = "zsh"
	darwinOS  = "darwin"
	linuxOS   = "linux"
)

// file locations
var (
	fishConfigPath           = "~/.config/fish/completions/pls.fish"
	linuxBashPath            = "/etc/bash_completion.d/pls"
	macBashPath              = "/usr/local/etc/bash_completion.d/pls"
	zshConfigPath            = "~/.zshrc"
	zshSessionCompletionPath = "~/.oh-my-zsh/custom/plugins/zsh-autosuggestions/_pls"
)

// GetShellConfigPath is
func GetShellConfigPath(shellType string) (string, error) {
	currentOS := utils.GetPlatform()
	switch currentOS {
	case "darwin", "linux":
		break
	default:
		return "", fmt.Errorf("%s is not currently supported for shell completion", currentOS)
	}

	switch shellType {
	case fishShell:
		return fishConfigPath, nil
	case bashShell:
		switch currentOS {
		case darwinOS:
			return macBashPath, nil
		case linuxOS:
			return linuxBashPath, nil
		}
	case zshShell:
		return zshSessionCompletionPath, nil
	}

	return "", fmt.Errorf("%s is not currently supported", shellType)
}
