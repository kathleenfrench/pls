package pls

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"V"},
	Short:   "print the current version of pls",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("pls version: %s\n", Version)
	},
}
