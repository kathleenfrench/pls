package pls

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/cobra"
)

const bashHelp = `
Bash:

$ source <(pls add completion bash)

# To load completions for each session, execute once:
Linux:
  $ pls add completion bash > /etc/bash_completion.d/pls
MacOS:
  $ pls add completion bash > /usr/local/etc/bash_completion.d/pls
`

const zshHelp = `
Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ pls add completion zsh > "${fpath[1]}/_pls"

# You will need to start a new shell for this setup to take effect.
`

const fishHelp = `
Fish:

$ pls add completion fish | source

# To load completions for each session, execute once:
$ pls add completion fish > ~/.config/fish/completions/pls.fish
`

// flags
var printOutput bool

// file locations
var (
	fishConfigPath           = "~/.config/fish/completions/pls.fish"
	linuxBashPath            = "/etc/bash_completion.d/pls"
	macBashPath              = "/usr/local/etc/bash_completion.d/pls"
	zshConfigPath            = "~/.zshrc"
	zshSessionCompletionPath = "${fpath[1]}/_pls"
)

func configPath(shellType string) (string, error) {
	currentOS := utils.GetPlatform()

	switch currentOS {
	case "darwin", "linux":
		break
	default:
		return "", fmt.Errorf("%s is not currently supported for shell completion", currentOS)
	}

	switch shellType {
	case "fish":
		return fishConfigPath, nil
	case "bash":
		switch currentOS {
		case "darwin":
			return macBashPath, nil
		case "linux":
			return linuxBashPath, nil
		}
	case "zsh":
		return zshConfigPath, nil
	}

	return "", fmt.Errorf("%s is not currently supported", shellType)
}

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish]",
	Short:                 "add shell completion for pls",
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shellChoice := args[0]
		color.HiYellow(fmt.Sprintf("pls will try to install completion for %s...", shellChoice))

		writer := os.Stdout
		if !printOutput {
			writer = nil
		}

		switch shellChoice {
		case "bash":
			color.HiGreen(bashHelp)
			cmd.Root().GenBashCompletion(writer)
		case "zsh":
			color.HiGreen(zshHelp)
			cmd.Root().GenZshCompletion(writer)
		case "fish":
			color.HiGreen(fishHelp)
			cmd.Root().GenFishCompletion(writer, true)
		}
	},
}

// completion flags
func init() {
	completionCmd.Flags().BoolVar(&printOutput, "print", false, "print the output of the generate completion script that will be copied to the correct path based off of your shell selection")
}
