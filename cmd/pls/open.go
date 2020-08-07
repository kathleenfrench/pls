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
	Example: fmt.Sprintf("pls open [opens dropdown GUI of your url shortcuts]\npls open https://google.com\npls open google.com"),
	Long:    "`pls` comes with a few pre-baked url shortcuts, but the rest is up to you. view your configs (`pls show configs`) to see what shortcuts have already been set to the `webshort` property. if you ever want to update these values - whether that be changing an existing url or adding a new shortcut - simply run `pls update configs` and follow the onscreen prompts.",
	Aliases: []string{"o"},
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			faveKeys := utils.GetKeysFromMapString(plsCfg.WebShortcuts)
			choice := gui.SelectPromptWithResponse("select a target URL from your shortcuts", faveKeys, nil, false)
			color.Blue(fmt.Sprintf("opening %s...", choice))
			utils.OpenURLInDefaultBrowser(plsCfg.WebShortcuts[choice])
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
