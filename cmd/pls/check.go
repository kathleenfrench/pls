package pls

import "github.com/spf13/cobra"

var checkCmd = &cobra.Command{
	Use:     "check",
	Short:   "run a check on a given resource",
	Aliases: []string{"ck"},
}
