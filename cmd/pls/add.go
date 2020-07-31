package pls

import "github.com/spf13/cobra"

var addSubCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "add various resources",
}

func init() {
	addSubCmd.AddCommand(completionCmd)
}
