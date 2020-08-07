package pls

import (
	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/clean"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/cobra"
)

// AllGitRepositories is a flag used to indicate an operation is to be performed on all git repositories found on the user's machine
var AllGitRepositories bool

var cleanCmd = &cobra.Command{
	Use:     "cleanup",
	Short:   "cleanup subcommands",
	Long:    "cleanup is used for auditing various local resources and determining what artifacts, if any, are eliglble for removal so as to free up more space on your machine",
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
			err := clean.CurrentRepoGitBranches()
			if err != nil {
				utils.ExitWithError(err)
			}
		}
	},
}

var cleanDockerSubCmd = &cobra.Command{
	Use:     "docker",
	Aliases: []string{"d", "dk"},
	Example: "pls cleanup docker",
	Short:   "prune local docker resources to free up space",
	Run: func(cmd *cobra.Command, args []string) {
		err := clean.SystemPrune()
		if err != nil {
			utils.ExitWithError(err)
		}
	},
}

func init() {
	cleanCmd.AddCommand(cleanGitSubCmd)
	cleanCmd.AddCommand(cleanDockerSubCmd)

	// clean branches
	cleanGitSubCmd.Flags().BoolVarP(&AllGitRepositories, "all", "a", false, "cleanup branches in all git repository folders on your machine")
}
