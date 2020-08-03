package gitpls

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	"github.com/kathleenfrench/pls/internal/config"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/kathleenfrench/pls/pkg/web/git"
)

// CreateGitOrganizationsDropdown prompts the user to select an organization from a dropdown
func CreateGitOrganizationsDropdown(organizations []*github.Organization) *github.Organization {
	names := []string{}
	nameMap := make(map[string]*github.Organization)

	for _, o := range organizations {
		names = append(names, o.GetLogin())
		nameMap[o.GetLogin()] = o
	}

	choice := gui.SelectPromptWithResponse("select an organization", names)
	return nameMap[choice]
}

// ChooseWithToDoWithOrganization lets the user decide with to do with their chosen organization
func ChooseWithToDoWithOrganization(organization *github.Organization, settings config.Settings) error {
	ctx := context.Background()
	_ = git.NewClient(ctx, settings.GitToken)

	opts := []string{openInBrowser, getOrganizationRepos, exitSelections}
	selected := gui.SelectPromptWithResponse(fmt.Sprintf("what would you like to do with %s?", organization.GetName()), opts)

	switch selected {
	case openInBrowser:
		utils.OpenURLInDefaultBrowser(organization.GetHTMLURL())
	case getOrganizationRepos:
		repos, err := FetchReposInOrganization(organization.GetName(), settings.GitToken)
		if err != nil {
			return err
		}

		choice := CreateGitRepoDropdown(repos)
		_ = ChooseWhatToDoWithRepo(choice, settings)
	case exitSelections:
		gui.Exit()
	}

	return nil
}

// FetchOrganizations fetches github organizations by user
func FetchOrganizations(username string, token string) ([]*github.Organization, error) {
	ctx := context.Background()
	gc := git.NewClient(ctx, token)
	opts := github.ListOptions{
		PerPage: 100,
	}

	orgs, _, err := gc.Organizations.List(ctx, "", &opts)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}
