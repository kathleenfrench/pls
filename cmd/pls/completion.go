package pls

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "A brief description of your command",
	Long: `To load completions:

Bash:

$ source <(pls completion bash)

# To load completions for each session, execute once:
Linux:
  $ pls completion bash > /etc/bash_completion.d/pls
MacOS:
  $ pls completion bash > /usr/local/etc/bash_completion.d/pls

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ pls completion zsh > "${fpath[1]}/_pls"

# You will need to start a new shell for this setup to take effect.

Fish:

$ pls completion fish | source

# To load completions for each session, execute once:
$ pls completion fish > ~/.config/fish/completions/pls.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}
