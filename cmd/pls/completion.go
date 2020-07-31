package pls

import (
	"os"

	"github.com/fatih/color"
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

var printOutput bool

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish|powershell]",
	Short:                 "add shell completion for pls",
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		writer := os.Stdout
		if !printOutput {
			writer = nil
		}
		switch args[0] {
		case "bash":
			color.HiGreen(bashHelp)
			cmd.Root().GenBashCompletion(writer)
		case "zsh":
			color.HiGreen(zshHelp)
			cmd.Root().GenZshCompletion(writer)
		case "fish":
			color.HiGreen(fishHelp)
			cmd.Root().GenFishCompletion(writer, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(writer)
		}
	},
}

// completion flags
func init() {
	completionCmd.Flags().BoolVar(&printOutput, "print", false, "print the output of the generate completion script that will be copied to the correct path based off of your shell selection")
}
