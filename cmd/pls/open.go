package pls

import "github.com/spf13/cobra"

var openCmd = &cobra.Command{
	Use:     "open",
	Short:   "open a resource",
	Aliases: []string{"o"},
}
