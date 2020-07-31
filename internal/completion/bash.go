package completion

import (
	"fmt"

	"github.com/kathleenfrench/pls/pkg/utils"
)

// BashHelp prints for users who elect to install their bash completion scripts manually
const BashHelp = `
Bash:

$ source <(pls add complete bash)

# To load completions for each session, execute once:
Linux:
  $ pls add complete bash > /etc/bash_completion.d/pls
MacOS:
  $ pls add complete bash > /usr/local/etc/bash_completion.d/pls
`

// BashInstall adds the bash completion scripts to the user's host machine
func BashInstall(path string) error {
	bashInstallCmd := fmt.Sprintf("pls add complete bash > %s", path)
	_, err := utils.BashExec(bashInstallCmd)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}
