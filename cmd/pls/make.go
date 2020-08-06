package pls

import (
	gitpls "github.com/kathleenfrench/pls/internal/web/git"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/cobra"
)

var makeCmd = &cobra.Command{
	Use:     "make",
	Aliases: []string{"mk"},
	Short:   "let pls create something for you",
}

var aCmd = &cobra.Command{
	Use:   "a",
	Short: "make a single one of a given resource",
}

var makePRCmd = &cobra.Command{
	Use:     "pullrequest",
	Aliases: []string{"pr", "pull"},
	Run: func(cmd *cobra.Command, args []string) {
		err := gitpls.CreatePullRequestFromCWD(plsCfg)
		if err != nil {
			utils.ExitWithError(err)
		}
	},
}

func init() {
	makeCmd.AddCommand(aCmd)
	aCmd.AddCommand(makePRCmd)
}
