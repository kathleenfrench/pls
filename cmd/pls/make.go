package pls

import "github.com/spf13/cobra"

var makeCmd = &cobra.Command{
	Use:     "make",
	Aliases: []string{"mk"},
	Short:   "let pls create something for you",
}
