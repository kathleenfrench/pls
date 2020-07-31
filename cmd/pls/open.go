package pls

import (
	"fmt"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:     "open",
	Short:   "open a resource",
	Aliases: []string{"o"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("opening...")
	},
}
