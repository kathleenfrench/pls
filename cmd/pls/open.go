package pls

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:     "open",
	Short:   "open any url in your default browser from the command line, or select from a set of common favorites",
	Aliases: []string{"o"},
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// dropdown defaults
			gui.Exit()
		}

		var (
			validatedURL string
			err          error
			arg          = args[0]
		)

		// first check if it's set as a favorite/shortcut
		shortcuts := plsCfg.WebShortcuts
		if _, ok := shortcuts[arg]; ok {
			validatedURL = shortcuts[arg]
		} else {
			// if it's not, then see if we can validate it as a url
			validatedURL, err = utils.ValidateURL(arg)
			if err != nil {
				utils.ExitWithError(fmt.Sprintf("%s - %s is not a valid URL", err, arg))
			}

		}

		color.Blue(fmt.Sprintf("opening %s...", validatedURL))
		utils.OpenURLInDefaultBrowser(validatedURL)

	},
}
