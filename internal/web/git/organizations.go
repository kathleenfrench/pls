package gitpls

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
)

// CreateGitOrganizationsDropdown prompts the user to select an organization from a dropdown
func CreateGitOrganizationsDropdown(organizations []*github.Organization) *github.Organization {
	names := []string{}
	nameMap := make(map[string]*github.Organization)

	for _, o := range organizations {
		color.Green(fmt.Sprintf("org: %v", o))
		names = append(names, o.GetLogin())
		nameMap[o.GetLogin()] = o
	}

	choice := gui.SelectPromptWithResponse("select an organization", names)
	return nameMap[choice]
}

// ChooseWithToDoWithOrganization lets the user decide with to do with their chosen organization
func ChooseWithToDoWithOrganization(gc *github.Client, organization *github.Organization) error {
	opts := []string{openInBrowser, getOrganizationRepos, exitSelections}
	selected := gui.SelectPromptWithResponse(fmt.Sprintf("what would you like to do with %s?", organization.GetName()), opts)

	switch selected {
	case openInBrowser:
		utils.OpenURLInDefaultBrowser(organization.GetHTMLURL())
	case getOrganizationRepos:
		color.HiRed("TODO")
		break
	case exitSelections:
		gui.Exit()
	}

	return nil
}
