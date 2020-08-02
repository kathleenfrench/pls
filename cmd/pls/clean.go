package pls

import (
	"github.com/kathleenfrench/pls/pkg/web/git"
	"github.com/spf13/cobra"
)

// AllGitRepositories is a flag used to indicate an operation is to be performed on all git repositories found on the user's machine
var AllGitRepositories bool

var cleanCmd = &cobra.Command{
	Use:     "clean",
	Short:   "cleanup subcommands",
	Aliases: []string{"c"},
}

// git clean --------------------------------
var cleanGitSubCmd = &cobra.Command{
	Use:     "git",
	Short:   "git cleanup subcommands",
	Aliases: []string{"g"},
}

var cleanGitBranchesSubCmd = &cobra.Command{
	Use:     "branches",
	Short:   "remove git branches that have already been merged into master - defaults to just within the current working directory",
	Aliases: []string{"b"},
	Run: func(cmd *cobra.Command, args []string) {
		git.CleanupCurrentBranches()
	},
}

func init() {
	cleanCmd.AddCommand(cleanGitSubCmd)

	// clean branches
	cleanGitSubCmd.AddCommand(cleanGitBranchesSubCmd)
	cleanGitBranchesSubCmd.Flags().BoolVarP(&AllGitRepositories, "all", "a", false, "cleanup branches in all git repository folders on your machine")
}
