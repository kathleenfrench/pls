package gitpls

import (
	"context"
	"errors"
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

	choice := gui.SelectPromptWithResponse("select an organization", names, nil, false)
	return nameMap[choice]
}

// ChooseWithToDoWithOrganization lets the user decide with to do with their chosen organization
func ChooseWithToDoWithOrganization(organization *github.Organization, settings config.Settings, useEnterprise bool) error {
	// var gc *github.Client
	// var err error

	// ctx := context.Background()

	// if useEnterprise {
	// 	gc, err = git.NewEnterpriseClient(ctx, settings.GitEnterpriseHostname, settings.GitEnterpriseToken)
	// 	if err != nil {
	// 		return err
	// 	}
	// } else {
	// 	gc = git.NewClient(ctx, settings.GitToken)
	// }

	login := organization.GetLogin()
	if login == "" {
		login = organization.GetName()
		if login == "" {
			return errors.New("no name could be parsed for this organization")
		}
	}

	opts := []string{openInBrowser, getOrganizationRepos, exitSelections}
	selected := gui.SelectPromptWithResponse(fmt.Sprintf("what would you like to do with %s?", organization.GetLogin()), opts, nil, false)

	switch selected {
	case openInBrowser:
		utils.OpenURLInDefaultBrowser(organization.GetHTMLURL())
	case getOrganizationRepos:
		gui.Spin.Start()
		repos, err := FetchReposInOrganization(organization.GetLogin(), settings, useEnterprise)
		gui.Spin.Stop()
		if err != nil {
			return err
		}

		choice := CreateGitRepoDropdown(repos)
		err = ChooseWhatToDoWithRepo(choice, settings, useEnterprise)
		if err != nil {
			return err
		}
	case exitSelections:
		gui.Exit()
	}

	return nil
}

// FetchOrganizations fetches github organizations by user
func FetchOrganizations(username string, settings config.Settings, useEnterprise bool) ([]*github.Organization, error) {
	var gc *github.Client
	var err error

	ctx := context.Background()

	if useEnterprise {
		gc, err = git.NewEnterpriseClient(ctx, settings.GitEnterpriseHostname, settings.GitEnterpriseToken)
		if err != nil {
			return nil, err
		}
	} else {
		gc = git.NewClient(ctx, settings.GitToken)
	}

	opts := github.ListOptions{
		PerPage: 100,
	}

	gui.Spin.Start()
	orgs, _, err := gc.Organizations.List(ctx, "", &opts)
	gui.Spin.Stop()
	if err != nil {
		return nil, err
	}

	return orgs, nil
}
