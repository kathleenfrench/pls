package completion

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
)

// ZshHelp is printed for users who elect to manually install the completion scripts
const ZshHelp = `
Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ pls add complete zsh > "${fpath[1]}/_pls"

# You will need to start a new shell for this setup to take effect.
`

// ZshInstall handles installing zsh completion scripts to the user's host machine
func ZshInstall(path string) error {
	// verify auto update is enabled
	if !confirmZshShellCompletionEnabled() {
		color.HiYellow(`before adding zsh completion scripts, you must enable the feature, run:\necho "autoload -U compinit; compinit" >> ~/.zshrc\n\nnow, try installing through pls again!`)
		os.Exit(1)
	}

	zshCmd := fmt.Sprintf("pls add complete zsh > %s", path)
	color.HiYellow(fmt.Sprintf("[RUNNING]: %s", zshCmd))
	_, err := utils.BashExec(zshCmd)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}

func confirmZshShellCompletionEnabled() bool {
	zshAutoEnabled := gui.ConfirmPrompt("is zsh shell completion enabled?", `you can check this by searching for echo "autoload -U compinit; compinit" in your .zshrc file`, true, true)
	return zshAutoEnabled
}
