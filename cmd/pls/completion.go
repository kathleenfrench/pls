package pls

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/internal/completion"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/cobra"
)

var completeMethodCmd = &cobra.Command{
	Use:                   "complete [bash|zsh|fish]",
	Hidden:                true,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	Run: func(cmd *cobra.Command, args []string) {
		shellChoice := args[0]
		switch shellChoice {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		}
	},
}

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:       "completion [bash|zsh|fish]",
	Short:     "add shell completion for pls",
	Long:      "writes completion scripts to your local machine based off of your preferred shell. `pls` currently has support for: `bash`, `zsh`, and `fish`.",
	ValidArgs: []string{"bash", "zsh", "fish"},
	Args:      cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shellChoice := args[0]

		// check if it's ok for pls to install the scripts
		canInstall := completion.PermissionToInstallPrompt(shellChoice)
		if !canInstall {
			completion.PrintInstallationInstructionsToStdout(cmd, shellChoice)
			os.Exit(1)
		}

		// determine the correct path for installing completion scripts based off of shell selection
		shellChoiceCfgPath, err := completion.GetShellConfigPath(shellChoice)
		if err != nil {
			panic(err)
		}

		// confirm that this is the right path
		correctPath := completion.ConfirmInstallationPath(shellChoiceCfgPath, shellChoice)
		if !correctPath {
			utils.PrintError("hmm, i'll need to add support for custom paths then...")
			os.Exit(1)
		}

		// install the completion scripts on the host machine
		err = install(shellChoice, shellChoiceCfgPath)
		if err != nil {
			utils.PrintError(err)
			os.Exit(1)
		}

		color.HiGreen("installation complete! reload your shell for the settings to take effect!")
	},
}

func install(shellType string, configPath string) error {
	switch shellType {
	case "bash":
		err := completion.BashInstall(configPath)
		if err != nil {
			os.Exit(1)
		}
	case "zsh":
		err := completion.ZshInstall(configPath)
		if err != nil {
			os.Exit(1)
		}
	case "fish":
		err := completion.FishInstall(configPath)
		if err != nil {
			os.Exit(1)
		}
	default:
		return fmt.Errorf("%s is not supported", shellType)
	}

	return nil
}
