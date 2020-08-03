package gitpls

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
)

// CreateGitRepoDropdown prompts the user for selecting from a dropdown of repositories
func CreateGitRepoDropdown(repositories []*github.Repository) *github.Repository {
	names := []string{}
	nameMap := make(map[string]*github.Repository)

	for _, r := range repositories {
		names = append(names, r.GetName())
		nameMap[r.GetName()] = r
	}

	choice := gui.SelectPromptWithResponse("select a repository", names)
	return nameMap[choice]
}

// ChooseWhatToDoWithRepo lets the user decide what to do with their chosen repo
func ChooseWhatToDoWithRepo(repository *github.Repository) error {
	opts := []string{openInBrowser, cloneRepo, exitSelections}
	// open in browser, clone
	selected := gui.SelectPromptWithResponse(fmt.Sprintf("what would you like to do with %s?", repository.GetName()), opts)

	switch selected {
	case openInBrowser:
		utils.OpenURLInDefaultBrowser(repository.GetHTMLURL())
	case cloneRepo:
		color.HiRed("TODO")
		break
	case exitSelections:
		gui.Exit()
	}

	return nil
}
