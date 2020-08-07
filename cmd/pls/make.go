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
			return fmt.Sprintf("/pls/%s/", strings.ToLower(base))
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
