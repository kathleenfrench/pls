package gitpls

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/google/go-github/v32/github"
	"github.com/jedib0t/go-pretty/v6/table"
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

	choice := gui.SelectPromptWithResponse("select a repository", names, nil, false)
	return nameMap[choice]
}

// ChooseWhatToDoWithRepo lets the user decide what to do with their chosen repo
func ChooseWhatToDoWithRepo(repository *github.Repository, settings config.Settings, useEnterprise bool) error {
	rows := []table.Row{
		{"full name", repository.GetFullName()},
		{"description", repository.GetDescription()},
		{"default branch", repository.GetDefaultBranch()},
		{"created", humanize.Time(repository.GetCreatedAt().Time)},
		{"updated", humanize.Time(repository.GetUpdatedAt().Time)},
		{"language", repository.GetLanguage()},
	}

	gui.SideBySideTable(rows)

	opts := []string{openInBrowser, cloneRepo, exitSelections}
	selected := gui.SelectPromptWithResponse(fmt.Sprintf("what would you like to do with %s?", repository.GetName()), opts, nil, true)

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

		selected := gui.SelectPromptWithResponse(fmt.Sprintf("where do you want to clone in %s?", repository.GetName()), pathChoices, settings.DefaultCodeDir, false)

		switch selected {
		case settings.DefaultCodeDir:
			clonePath = settings.DefaultCodeDir
		case "Current Directory":
			clonePath = ""
		case "Custom Directory":
			clonePath = gui.InputPromptWithResponse(fmt.Sprintf("what is the *full* path to the directory?"), "", true)
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
