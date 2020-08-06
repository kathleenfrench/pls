package pls

import (
	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/kathleenfrench/pls/pkg/web/git"
	"github.com/spf13/cobra"
)

// AllGitRepositories is a flag used to indicate an operation is to be performed on all git repositories found on the user's machine
var AllGitRepositories bool

var cleanCmd = &cobra.Command{
	Use:     "cleanup",
	Short:   "cleanup subcommands",
	Aliases: []string{"c", "clean"},
}

// git clean --------------------------------
var cleanGitSubCmd = &cobra.Command{
	Use:     "git",
	Aliases: []string{"g"},
	Example: "pls cleanup git",
	Short:   "remove git branches that have already been merged into master - defaults to just within the current working directory",
	Run: func(cmd *cobra.Command, args []string) {
		if AllGitRepositories {
			color.HiRed("TO DO")
		} else {
			err := git.CleanupCurrentBranches()
			if err != nil {
				utils.ExitWithError(err)
			}
		}
	},
}

func init() {
	cleanCmd.AddCommand(cleanGitSubCmd)

	// clean branches
	cleanGitSubCmd.Flags().BoolVarP(&AllGitRepositories, "all", "a", false, "cleanup branches in all git repository folders on your machine")
}
