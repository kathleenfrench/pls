package pls

import (
	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/kathleenfrench/pls/pkg/web/git"
	"github.com/spf13/cobra"
)

// CURRENTLY UNUSED

var tryCmd = &cobra.Command{
	Use:     "try",
	Short:   "try to do something",
	Aliases: []string{"t"},
	Run: func(cmd *cobra.Command, args []string) {
		unpushed, err := git.HasUnpushedChangesOrCommits()
		if err != nil {
			utils.ExitWithError(err)
		}

		color.HiGreen("has unpushed commits? %v", unpushed)
	},
	Hidden: true,
}
