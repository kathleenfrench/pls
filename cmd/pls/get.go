package pls

import (
	"context"
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	gitpls "github.com/kathleenfrench/pls/internal/web/git"
	"github.com/kathleenfrench/pls/pkg/gui"
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

var fetchTypeChecker = map[string]string{
	"by":     "other_user",
	"for":    "other_user",
	"in":     "organization",
	"my":     "current_user",
	"called": "search",
}

var gitRepos = &cobra.Command{
	Use:     "repos",
	Aliases: []string{"r", "repositories", "repo", "repository"},
	Short:   "interact with someone else's github repositories",
	Example: color.HiGreenString(fmt.Sprintf("\npls get repos for <username>\npls get repos by <username>\npls get repos in <organization>")),
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		gui.Spin.Start()

		var repos []*github.Repository
		ctx := context.Background()
		found, ok := fetchTypeChecker[args[0]]
		if !ok {
			utils.ExitWithError(fmt.Sprintf("%s is not a valid command", found))
		}

		switch found {
		case "other_user":
			username := args[1]
			otherUserRepos, err := fetchRepos(ctx, username)
			gui.Spin.Stop()
			if err != nil {
				utils.ExitWithError(err)
			}

			repos = otherUserRepos
		case "organization":
			organization := args[1]
			color.HiYellow(fmt.Sprintf("fetching repositories in the %s organization...", organization))
			orgRepos, err := fetchReposInOrganization(ctx, organization)
			gui.Spin.Stop()
			if err != nil {
				utils.ExitWithError(err)
			}

			repos = orgRepos
		case "current_user":
			username := ""
			currentUserRepos, err := fetchRepos(ctx, username)
			gui.Spin.Stop()
			if err != nil {
				utils.ExitWithError(err)
			}

			repos = currentUserRepos
		case "search":
			color.HiRed("TODO")
		default:
			utils.ExitWithError("sorry, but i don't understand what you want")
		}

		color.HiYellow(fmt.Sprintf("%d repositories returned", len(repos)))
		choice := gitpls.CreateGitRepoDropdown(repos)
		_ = gitpls.ChooseWhatToDoWithRepo(choice)
		gui.Exit()
	},
}

//--------------------------------------- HELPERS

func fetchReposInOrganization(ctx context.Context, organization string) ([]*github.Repository, error) {
	var allOrgRepos []*github.Repository
	gc := git.NewClient(ctx, plsCfg.GitToken)
	opts := github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		repos, resp, err := gc.Repositories.ListByOrg(ctx, organization, &opts)
		if err != nil {
			return nil, err
		}

		allOrgRepos = append(allOrgRepos, repos...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return allOrgRepos, nil
}

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
		if args[0] == "in" {
			return "", errors.New("organization")
		}

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
