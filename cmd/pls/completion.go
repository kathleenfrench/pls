package pls

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/internal/completion"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/cobra"
)

// flags
var printOutput bool

// install scripts
var fishInit = exec.Command("pls", "add", "complete", "fish", "|", "source")

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

		// check if it's ok for pls to install the scripts
		canInstall := completion.PermissionToInstallPrompt(shellChoice)
		if !canInstall {
			completion.PrintInstallationInstructionsToStdout(cmd, shellChoice)
			os.Exit(1)
		}

		shellChoiceCfgPath, err := completion.GetShellConfigPath(shellChoice)
		if err != nil {
			panic(err)
		}

		correctPath := completion.ConfirmInstallationPath(shellChoiceCfgPath, shellChoice)
		if !correctPath {
			utils.PrintError("hmm, i'll need to add support for custom paths then...")
			os.Exit(1)
		}

		err = install(shellChoice, shellChoiceCfgPath)
		if err != nil {
			utils.PrintError(err)
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
				_, err := utils.BashExec(`echo "autoload -U compinit; compinit" >> ~/.zshrc`)
				if err != nil {
					utils.PrintError(err)
					os.Exit(1)
				}
			} else {
				color.HiYellow(fmt.Sprintf(`Ok, run: echo "autoload -U compinit; compinit" >> ~/.zshrc"\nafter you've reloaded your shell, come back and re-run the completion command`))
				os.Exit(1)
			}
		}

		zshCmd := fmt.Sprintf("pls add complete zsh > %s", fmt.Sprintf("%s/_pls", configPath))
		color.HiYellow(fmt.Sprintf("[RUNNING]: %s", zshCmd))
		_, err := utils.BashExec(zshCmd)
		if err != nil {
			utils.PrintError(err)
			os.Exit(1)
		}

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

// completion flags
func init() {
	completionCmd.Flags().BoolVar(&printOutput, "print", false, "print the output of the generate completion script that will be copied to the correct path based off of your shell selection")
}
