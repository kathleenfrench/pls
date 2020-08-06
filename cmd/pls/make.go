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

var aCmd = &cobra.Command{
	Use:   "a",
	Short: "make a single one of a given resource",
}

var makePRCmd = &cobra.Command{
	Use:     "pullrequest",
	Aliases: []string{"pr", "pull"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("make a PR?")
	},
}

func init() {
	makeCmd.AddCommand(aCmd)
	aCmd.AddCommand(makePRCmd)
}
