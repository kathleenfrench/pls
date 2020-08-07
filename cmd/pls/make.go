package pls

import (
	gitpls "github.com/kathleenfrench/pls/internal/web/git"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
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

var makeDocsCmd = &cobra.Command{
	Use:    "docs",
	Hidden: true,
	Short:  "generate markdown documentation for pls commands",
	Run: func(cmd *cobra.Command, args []string) {
		topCmd := cmd.Root()
		err := doc.GenMarkdownTree(topCmd, "./docs")
		if err != nil {
			utils.ExitWithError(err)
		}

		gui.Log(":thumbs_up:", "docs generated in ./docs!", nil)
	},
}

func init() {
	makeCmd.AddCommand(aCmd)
	aCmd.AddCommand(makePRCmd)
	makeCmd.AddCommand(makeDocsCmd)
}
