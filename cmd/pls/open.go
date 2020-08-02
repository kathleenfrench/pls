package pls

import (
	"fmt"

	"github.com/fatih/color"
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
		} else {
			validatedURL, err := utils.ValidateURL(args[0])
			if err != nil {
				utils.ExitWithError(fmt.Sprintf("%s - %s is not a valid URL", err, args[0]))
			}
			color.Blue(fmt.Sprintf("opening %s...", validatedURL))
			utils.OpenURLInDefaultBrowser(validatedURL)
		}
	},
}
