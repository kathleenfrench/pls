package pls

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	gitpls "github.com/kathleenfrench/pls/internal/web/git"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
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
		color.HiRed("TODO")
		orgs, err := gitpls.FetchOrganizations(plsCfg.GitUsername, plsCfg.GitToken)
		if err != nil {
			utils.ExitWithError(err)
		}

		choice := gitpls.CreateGitOrganizationsDropdown(orgs)
		_ = gitpls.ChooseWithToDoWithOrganization(choice, plsCfg)
	},
}

var gitMyRepos = &cobra.Command{
	Use:     "repos",
	Aliases: []string{"r", "repositories", "repo", "repository"},
	Short:   "interact with your github repositories",
	Example: color.HiYellowString("pls get my repos"),
	Run: func(cmd *cobra.Command, args []string) {
		repos, err := gitpls.FetchUserRepos("", plsCfg.GitToken)
		if err != nil {
			utils.ExitWithError(err)
		}

		choice := gitpls.CreateGitRepoDropdown(repos)
		_ = gitpls.ChooseWhatToDoWithRepo(choice, plsCfg)
	},
}

var gitMyOrgs = &cobra.Command{
	Use:     "orgs",
	Aliases: []string{"o", "org", "organization", "organizations"},
	Short:   "interact with your github organizations",
	Example: color.HiYellowString("pls get my orgs"),
	Run: func(cmd *cobra.Command, args []string) {
		orgs, err := gitpls.FetchOrganizations("", plsCfg.GitToken)
		if err != nil {
			utils.ExitWithError(err)
		}

		choice := gitpls.CreateGitOrganizationsDropdown(orgs)
		_ = gitpls.ChooseWithToDoWithOrganization(choice, plsCfg)
	},
}

// pls get prs to review (current) | (in <>) | <everywhere>
// pls get prs i merged (current) | (in <>) | <everywhere>
// pls get prs i closed (current) | (in <>) | <everywhere>

var gitMyPRs = &cobra.Command{
	Use:     "prs",
	Aliases: []string{"pulls", "pull", "pr"},
	Short:   "interact with your pull requests",
	Example: color.HiYellowString("\n[PRs in current directory's repository]: pls get my prs\n[PRs in a repository you own]: pls get my prs in myrepo\n[PRs in another's repository]: pls get my prs in organization/repo\n[PRs from all of github]: pls get my prs everywhere"),
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
		case 1:
			// everywhere check
			single := args[0]
			if single != "everywhere" && single != "all" {
				utils.ExitWithError(fmt.Sprintf("%s is not a valid argument", single))
			}

			prs, err := gitpls.FetchUserPullRequestsEverywhere(plsCfg)
			if err != nil {
				utils.ExitWithError(err)
			}

			for _, pr := range prs {
				color.HiGreen(fmt.Sprintf("%v", pr))
			}
			// fetch all
		case 2:
			// pls get my prs in <repo> (owned)
			// pls get my prs in <org>/<repo> (organization/another person's repo)
		default:
			utils.ExitWithError("invalid input, try running `pls get my prs --help`")
		}
	},
}

var repoFetchTypeChecker = map[string]string{
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
		found, ok := repoFetchTypeChecker[args[0]]
		if !ok {
			utils.ExitWithError(fmt.Sprintf("%s is not a valid command", found))
		}

		switch found {
		case "other_user":
			username := args[1]
			otherUserRepos, err := gitpls.FetchUserRepos(username, plsCfg.GitToken)
			gui.Spin.Stop()
			if err != nil {
				utils.ExitWithError(err)
			}

			repos = otherUserRepos
		case "organization":
			organization := args[1]
			color.HiYellow(fmt.Sprintf("fetching repositories in the %s organization...", organization))
			orgRepos, err := gitpls.FetchReposInOrganization(organization, plsCfg.GitToken)
			gui.Spin.Stop()
			if err != nil {
				utils.ExitWithError(err)
			}

			repos = orgRepos
		case "current_user":
			username := ""
			currentUserRepos, err := gitpls.FetchUserRepos(username, plsCfg.GitToken)
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
		_ = gitpls.ChooseWhatToDoWithRepo(choice, plsCfg)
	},
}

//--------------------------------------- INIT

func init() {
	getCmd.AddCommand(myGetSubCmd)

	// get only yours
	myGetSubCmd.AddCommand(gitMyOrgs)
	myGetSubCmd.AddCommand(gitMyRepos)
	myGetSubCmd.AddCommand(gitMyPRs)

	// get someone else's
	getCmd.AddCommand(gitOrgs)
	getCmd.AddCommand(gitRepos)
}
