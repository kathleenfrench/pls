package pls

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	gitpls "github.com/kathleenfrench/pls/internal/web/git"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/kathleenfrench/pls/pkg/web/git"
	"github.com/spf13/cobra"
)

var gitUsername string

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"git"},
	Short:   "shorthand for `git` in most cases, but can also get you other stuff",
}

var gitOrgs = &cobra.Command{
	Use:     "orgs",
	Aliases: []string{"o", "org", "organization", "organizations"},
	Short:   "interact with github organizations",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		gc := git.NewClient(ctx, plsCfg.GitToken)
		orgs, _, err := gc.Organizations.List(ctx, plsCfg.GitUsername, nil)
		if err != nil {
			utils.ExitWithError(err)
		}

		for i, o := range orgs {
			color.Green(fmt.Sprintf("%d) %s\n", i+1, o.GetLogin()))
		}
	},
}

var gitRepos = &cobra.Command{
	Use:     "repos",
	Aliases: []string{"r", "repositories", "repo", "repository"},
	Short:   "interact with github repositories",
	Run: func(cmd *cobra.Command, args []string) {
		username := plsCfg.GitUsername
		if gitUsername != "" {
			username = gitUsername
		}

		ctx := context.Background()
		gc := git.NewClient(ctx, plsCfg.GitToken)
		repos, _, err := gc.Repositories.List(ctx, username, nil)
		if err != nil {
			utils.ExitWithError(err)
		}

		choice := gitpls.CreateGitRepoDropdown(repos)
		color.HiGreen("chosen repo: %v", choice)
	},
}

func init() {
	getCmd.PersistentFlags().StringVarP(&gitUsername, "username", "u", "", "when you want to specify a github username other than your own")
	getCmd.AddCommand(gitOrgs)
	getCmd.AddCommand(gitRepos)
}
