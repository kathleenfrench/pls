package gitpls

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	"github.com/kathleenfrench/pls/internal/config"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/kathleenfrench/pls/pkg/web/git"
)

// CreateGitRepoDropdown prompts the user for selecting from a dropdown of repositories
func CreateGitRepoDropdown(repositories []*github.Repository) *github.Repository {
	names := []string{}
	nameMap := make(map[string]*github.Repository)

	for _, r := range repositories {
		names = append(names, r.GetName())
		nameMap[r.GetName()] = r
	}

	choice := gui.SelectPromptWithResponse("select a repository", names, false)
	return nameMap[choice]
}

// ChooseWhatToDoWithRepo lets the user decide what to do with their chosen repo
func ChooseWhatToDoWithRepo(repository *github.Repository, settings config.Settings) error {
	opts := []string{openInBrowser, cloneRepo, exitSelections}
	selected := gui.SelectPromptWithResponse(fmt.Sprintf("what would you like to do with %s?", repository.GetName()), opts, false)

	switch selected {
	case openInBrowser:
		utils.OpenURLInDefaultBrowser(repository.GetHTMLURL())
	case cloneRepo:
		var clonePath string
		// give user choice between current directory, default codebase path, specify a custom path
		pathChoices := []string{
			settings.DefaultCodeDir,
			"Current Directory",
			"Custom Directory",
		}

		selected := gui.SelectPromptWithResponse(fmt.Sprintf("where do you want to clone in %s?", repository.GetName()), pathChoices, false)

		switch selected {
		case settings.DefaultCodeDir:
			clonePath = settings.DefaultCodeDir
		case "Current Directory":
			clonePath = ""
		case "Custom Directory":
			clonePath = gui.InputPromptWithResponse(fmt.Sprintf("what is the *full* path to the directory?"), "")
		}

		err := git.CloneRepository(repository.GetName(), repository.GetCloneURL(), clonePath)
		if err != nil {
			return err
		}
	case exitSelections:
		gui.Exit()
	}

	return nil
}

// FetchReposInOrganization fetches repositories in an organization
func FetchReposInOrganization(organization string, token string) ([]*github.Repository, error) {
	var allOrgRepos []*github.Repository

	ctx := context.Background()
	gc := git.NewClient(ctx, token)
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

// FetchUserRepos fetches repositories by user
func FetchUserRepos(username string, token string) ([]*github.Repository, error) {
	ctx := context.Background()
	gc := git.NewClient(ctx, token)
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
