package pls

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"

	gitpls "github.com/kathleenfrench/pls/internal/web/git"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
)

var (
	fetchAll bool
)

// ------------------------------------------------------

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"git"},
	Short:   "shorthand for `git` in most cases, but can also get you other stuff",
}

// --------------------------- ORGANIZATIONS

var gitOrgs = &cobra.Command{
	Use:     "orgs",
	Aliases: []string{"o", "org", "organization", "organizations"},
	Short:   "interact with someone else's github organizations",
	Run: func(cmd *cobra.Command, args []string) {
		color.HiRed("REALLY SORRY BUT I'M NOT AVAILABLE YET")
	},
}

// --------------------------- REPOS

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
	Example: "pls get repos for <username>\npls get repos by <username>\npls get repos in <organization>",
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
			otherUserRepos, err := gitpls.FetchUserRepos(username, plsCfg, work)
			gui.Spin.Stop()
			if err != nil {
				utils.ExitWithError(err)
			}

			repos = otherUserRepos
		case "organization":
			organization := args[1]
			color.HiYellow(fmt.Sprintf("fetching repositories in the %s organization...", organization))
			orgRepos, err := gitpls.FetchReposInOrganization(organization, plsCfg, work)
			gui.Spin.Stop()
			if err != nil {
				utils.ExitWithError(err)
			}

			repos = orgRepos
		case "current_user":
			username := ""
			currentUserRepos, err := gitpls.FetchUserRepos(username, plsCfg, work)
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
		err := gitpls.ChooseWhatToDoWithRepo(choice, plsCfg, work)
		if err != nil {
			utils.ExitWithError(err)
		}
	},
}

// ------------------------------------------------------
// INIT
// ------------------------------------------------------

func init() {
	getCmd.PersistentFlags().BoolVar(&fetchAll, "all", false, "search all of github")
	getCmd.AddCommand(gitOrgs)
	getCmd.AddCommand(gitRepos)
}
