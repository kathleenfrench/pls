package completion

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// PermissionToInstallPrompt prompts the user for whether pls can install the completion scripts on their host machine
func PermissionToInstallPrompt(shell string) bool {
	granted := false
	msg := fmt.Sprintf("do you want me to install the %s completion scripts on your machine?", shell)

	prompt := &survey.Confirm{
		Message: msg,
	}

	survey.AskOne(prompt, &granted)
	return granted
}

func ConfirmInstallationPath(path string, shell string) bool {
	correct := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("is %s the correct path for installing %s completions?", path, shell),
	}

	survey.AskOne(prompt, &correct)
	return correct
}
