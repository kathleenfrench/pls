package pls

import "github.com/spf13/cobra"

var cleanCmd = &cobra.Command{
	Use:     "clean",
	Short:   "cleanup subcommands",
	Aliases: []string{"c"},
}
