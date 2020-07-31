package completion

import (
	"fmt"
	"os/exec"

	"github.com/kathleenfrench/pls/pkg/utils"
)

// FishHelp prints to users who elect to manually install fish
const FishHelp = `
Fish:

$ pls add complete fish | source

# To load completions for each session, execute once:
$ pls add complete fish > ~/.config/fish/completions/pls.fish
`

var fishInit = exec.Command("pls", "add", "complete", "fish", "|", "source")

// FishInstall handles installing the fish completion scripts on the user's host machine
func FishInstall(path string) error {
	srcPlsCompletion := "pls add complete fish | source"
	_, err := utils.BashExec(srcPlsCompletion)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	fishInstallCmd := fmt.Sprintf("pls add complete fish > %s", path)
	_, err = utils.BashExec(fishInstallCmd)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}
