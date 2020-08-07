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
		outDir := "docs/content/commands"
		err := utils.CreateDir(outDir)
		if err != nil {
			utils.ExitWithError(err)
		}

		fp := internalutils.FrontMatter
		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return "/commands/" + strings.ToLower(base) + "/"
		}

		topCmd := cmd.Root()
		topCmd.DisableAutoGenTag = true

		err = internalutils.GenMarkdownDocumentation(topCmd, fmt.Sprintf("./%s", outDir), fp, linkHandler)
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
