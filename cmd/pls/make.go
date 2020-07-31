package pls

import (
	"fmt"

	"github.com/spf13/cobra"
)

var makeCmd = &cobra.Command{
	Use:     "make",
	Aliases: []string{"mk"},
	Short:   "let pls create something for you",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("making...")
	},
}
