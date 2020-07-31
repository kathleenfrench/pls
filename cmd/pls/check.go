package pls

import (
	"fmt"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:     "check",
	Short:   "run a check on a given resource",
	Aliases: []string{"ck"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("checking...")
	},
}
