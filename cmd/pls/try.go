package pls

import (
	"fmt"

	"github.com/spf13/cobra"
)

var tryCmd = &cobra.Command{
	Use:     "try",
	Short:   "try to do something",
	Aliases: []string{"t"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("trying...")
	},
}
