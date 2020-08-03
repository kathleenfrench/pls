package pls

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	gitpls "github.com/kathleenfrench/pls/internal/web/git"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/kathleenfrench/pls/pkg/web/git"
	"github.com/spf13/cobra"
)

//--------------------------------------- COMMANDS

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"git"},
	Short:   "shorthand for `git` in most cases, but can also get you other stuff",
}

var myGetSubCmd = &cobra.Command{
	Use:     "my",
	Aliases: []string{"m"},
	Short:   "fetch your stuff specifically",
}

var gitOrgs = &cobra.Command{
	Use:     "orgs",
	Aliases: []string{"o", "org", "organization", "organizations"},
	Short:   "interact with someone else's github organizations",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		gc := git.NewClient(ctx, plsCfg.GitToken)
		opts := github.ListOptions{
			PerPage: 100,
		}

		orgs, _, err := gc.Organizations.List(ctx, plsCfg.GitUsername, &opts)
		if err != nil {
			utils.ExitWithError(err)
		}

		choice := gitpls.CreateGitOrganizationsDropdown(orgs)
		_ = gitpls.ChooseWithToDoWithOrganization(gc, choice)
	},
}

var gitMyRepos = &cobra.Command{
	Use:     "repos",
	Aliases: []string{"r", "repositories", "repo", "repository"},
	Short:   "interact with your github repositories",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		repos, err := fetchRepos(ctx, "")
		if err != nil {
			utils.ExitWithError(err)
		}

		choice := gitpls.CreateGitRepoDropdown(repos)
		_ = gitpls.ChooseWhatToDoWithRepo(choice)
	},
}

var gitMyOrgs = &cobra.Command{
	Use:     "orgs",
	Aliases: []string{"o", "org", "organization", "organizations"},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		gc := git.NewClient(ctx, plsCfg.GitToken)
		opts := github.ListOptions{
			PerPage: 100,
		}

		orgs, _, err := gc.Organizations.List(ctx, "", &opts)
		if err != nil {
			utils.ExitWithError(err)
		}

		choice := gitpls.CreateGitOrganizationsDropdown(orgs)
		_ = gitpls.ChooseWithToDoWithOrganization(gc, choice)
	},
}

var gitRepos = &cobra.Command{
	Use:     "repos",
	Aliases: []string{"r", "repositories", "repo", "repository"},
	Short:   "interact with someone else's github repositories",
	Example: color.HiGreenString(fmt.Sprintf("\npls get repos for <username>\npls get repos by <username>")),
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		byUser, err := forWhoCheck(args)
		if err != nil {
			utils.ExitWithError(err)
		}

		ctx := context.Background()
		repos, err := fetchRepos(ctx, byUser)
		if err != nil {
			utils.ExitWithError(err)
		}

		choice := gitpls.CreateGitRepoDropdown(repos)
		_ = gitpls.ChooseWhatToDoWithRepo(choice)
	},
}

//--------------------------------------- HELPERS

func fetchRepos(ctx context.Context, username string) ([]*github.Repository, error) {
	gc := git.NewClient(ctx, plsCfg.GitToken)
	opts := &github.RepositoryListOptions{
		Affiliation: "owner",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allRepos []*github.Repository

	for {
		repos, resp, err := gc.Repositories.List(ctx, username, opts)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)
		color.HiRed(fmt.Sprintf("next page: %d", resp.NextPage))
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	if username == "" {
		username = "you"
	}

	color.HiGreen(fmt.Sprintf("%d repositories found for %s", len(allRepos), username))

	return allRepos, nil
}

func forWhoCheck(args []string) (username string, err error) {
	if args[0] != "for" && args[0] != "by" {
		return "", fmt.Errorf("%s is not a valid subcommand", args[0])
	}

	return args[1], nil
}

//--------------------------------------- INIT

func init() {
	// getCmd.PersistentFlags().StringVarP(&gitUsername, "username", "u", "", "when you want to specify a github username other than your own")

	getCmd.AddCommand(myGetSubCmd)

	// get only yours
	myGetSubCmd.AddCommand(gitMyOrgs)
	myGetSubCmd.AddCommand(gitMyRepos)

	// get someone else's
	getCmd.AddCommand(gitOrgs)
	getCmd.AddCommand(gitRepos)
}
