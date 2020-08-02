package pls

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"V"},
	Short:   "print the current version of pls",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
	Hidden: true,
}

func printVersion() {
	color.HiYellow(fmt.Sprintf("VERSION %s\n", Version))
}
