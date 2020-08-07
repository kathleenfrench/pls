package pls

import (
	"fmt"
	"path"
	"strings"

	internalutils "github.com/kathleenfrench/pls/internal/utils"
	gitpls "github.com/kathleenfrench/pls/internal/web/git"
	"github.com/kathleenfrench/pls/pkg/gui"
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
	Short:   "let `pls` open a pull request for you based off of the branch in your current working directory",
	Long:    "`pls` will open a pull request for you, but only after verifying that your current local branch is synced with its remote ref. if it isn't synced, then `pls` will confirm whether or not you want to let `pls` handle adding, committing, and/or pushing the branch for you. once that's done, `pls` will prompt you for various values:\n- title\n- PR description\n- whether to link it to an existing issue\n- whether to create it as a draft (if using an enterprise account)\n\n`pls` will even spawn a temporary file in your preferred text editor for composing the body of the pull request. once you finish adding values, `pls` handles creating the PR!",
	Example: "pls make a pr",
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
		gui.Log(":popcorn:", "generating pls documentation...", nil)

		err := utils.CreateDir(internalutils.PublishDocsDirectory)
		if err != nil {
			utils.ExitWithError(err)
		}

		fp := internalutils.FrontMatter
		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return fmt.Sprintf("/pls/%s", strings.ToLower(base))
		}

		err = internalutils.GenMarkdownDocumentation(cmd.Root(), fmt.Sprintf("./%s", internalutils.PublishDocsDirectory), fp, linkHandler)
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
