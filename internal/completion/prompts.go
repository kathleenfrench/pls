package completion

import (
	"fmt"

	"github.com/kathleenfrench/pls/pkg/gui"
)

// PermissionToInstallPrompt prompts the user for whether pls can install the completion scripts on their host machine
func PermissionToInstallPrompt(shell string) bool {
	granted := gui.ConfirmPrompt(fmt.Sprintf("do you want me to install the %s completion scripts on your machine?", shell), "", true, true)
	return granted
}

// ConfirmInstallationPath verifies whether the given installation path is correct
func ConfirmInstallationPath(path string, shell string) bool {
	correct := gui.ConfirmPrompt(fmt.Sprintf("is %s the correct path for installing %s completions?", path, shell), "", true, true)
	return correct
}
