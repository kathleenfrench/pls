package completion

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// PrintInstallationInstructionsToStdout outputs installation instructions for compatible shells shells
func PrintInstallationInstructionsToStdout(cmd *cobra.Command, shellChoice string) {
	color.HiYellow("\nINSTALLATION:\n")

	switch shellChoice {
	case "bash":
		_ = cmd.Root().GenBashCompletion(os.Stdout)
		color.HiGreen(fmt.Sprintf("\n%s\n", BashHelp))
	case "zsh":
		_ = cmd.Root().GenZshCompletion(os.Stdout)
		color.HiGreen(fmt.Sprintf("\n%s\n", ZshHelp))
	case "fish":
		_ = cmd.Root().GenFishCompletion(os.Stdout, true)
		color.HiGreen(fmt.Sprintf("\n%s\n", FishHelp))
	}
}
