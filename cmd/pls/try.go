package pls

import "github.com/spf13/cobra"

var tryCmd = &cobra.Command{
	Use:     "try",
	Short:   "try to do something",
	Aliases: []string{"t"},
}
