package pls

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/cobra"
)

const bashHelp = `
Bash:

$ source <(pls add complete bash)

# To load completions for each session, execute once:
Linux:
  $ pls add complete bash > /etc/bash_completion.d/pls
MacOS:
  $ pls add complete bash > /usr/local/etc/bash_completion.d/pls
`

const zshHelp = `
Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ pls add complete zsh > "${fpath[1]}/_pls"

# You will need to start a new shell for this setup to take effect.
`

const fishHelp = `
Fish:

$ pls add complete fish | source

# To load completions for each session, execute once:
$ pls add complete fish > ~/.config/fish/completions/pls.fish
`

var fishInit = exec.Command("pls", "add", "complete", "fish", "|", "source")
var zshInstall = exec.Command("pls", "add", "complete", "zsh", ">", zshSessionCompletionPath)

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
	ValidArgs: []string{"bash", "zsh", "fish"},
	Args:      cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shellChoice := args[0]
		color.HiYellow(fmt.Sprintf("pls will try to install completion for %s...", shellChoice))

		shellChoiceCfgPath, err := configPath(shellChoice)
		if err != nil {
			panic(err)
		}

		correctPath := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("is %s the correct path for your %s configs?", shellChoiceCfgPath, shellChoice),
		}

		survey.AskOne(prompt, &correctPath)

		if !correctPath {
			color.HiRed("oops, then let's not do it!")
			os.Exit(1)
		}

		installItForYou := false
		prompt = &survey.Confirm{
			Message: fmt.Sprintf("do you want me to add the completion script commands to %s for you?", shellChoiceCfgPath),
		}

		survey.AskOne(prompt, &installItForYou)

		if !installItForYou {
			printInstallationInstructionsToStdout(cmd, shellChoice)
			os.Exit(1)
		}

		err = install(shellChoice, shellChoiceCfgPath)
		if err != nil {
			color.HiRed(fmt.Sprintf("[ERROR]: %s", err))
			os.Exit(1)
		}

		color.HiGreen("installation complete!")
	},
}

func install(shellType string, configPath string) error {
	switch shellType {
	case "bash":
		bashInstall := exec.Command("pls", "add", "complete", "bash", ">", configPath)
		bash := utils.ExecuteCommand(bashInstall)
		color.HiBlue(bash)
		break
	case "zsh":
		zshAutoEnabled := false
		prompt := &survey.Confirm{
			Message: "is zsh shell completion enabled?",
		}

		survey.AskOne(prompt, &zshAutoEnabled)
		if !zshAutoEnabled {
			enableAutoload := false
			prompt = &survey.Confirm{
				Message: "can i enable it for you?",
			}

			survey.AskOne(prompt, &enableAutoload)
			if enableAutoload {
				autoLoadCmd := exec.Command("echo", `"autoload -U compinit; compinit"`, ">>", "~/.zshrc")
				autoRes := utils.ExecuteCommand(autoLoadCmd)
				color.HiBlue(autoRes)
			} else {
				color.HiYellow(fmt.Sprintf(`Ok, run: echo "autoload -U compinit; compinit" >> ~/.zshrc"\nafter you've reloaded your shell, come back and re-run the completion command`))
				os.Exit(1)
			}
		}

		zsh := utils.ExecuteCommand(zshInstall)
		color.HiBlue(zsh)
		break
	case "fish":
		resInit := utils.ExecuteCommand(fishInit)
		color.HiYellow(resInit)
		fishInstall := exec.Command("pls", "add", "complete", "fish", ">", configPath)
		fish := utils.ExecuteCommand(fishInstall)
		color.HiBlue(fish)
		break
	default:
		return fmt.Errorf("%s is not supported", shellType)
	}

	return nil
}

func printInstallationInstructionsToStdout(cmd *cobra.Command, shellChoice string) {
	color.HiYellow("\nINSTALLATION:\n")

	switch shellChoice {
	case "bash":
		cmd.Root().GenBashCompletion(os.Stdout)
		color.HiGreen(fmt.Sprintf("\n%s\n", bashHelp))
	case "zsh":
		cmd.Root().GenZshCompletion(os.Stdout)
		color.HiGreen(fmt.Sprintf("\n%s\n", zshHelp))
	case "fish":
		cmd.Root().GenFishCompletion(os.Stdout, true)
		color.HiGreen(fmt.Sprintf("\n%s\n", fishHelp))
	}
}

// completion flags
func init() {
	completionCmd.Flags().BoolVar(&printOutput, "print", false, "print the output of the generate completion script that will be copied to the correct path based off of your shell selection")
}
