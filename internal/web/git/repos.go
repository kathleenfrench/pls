package gitpls

import "github.com/google/go-github/v32/github"

import "github.com/kathleenfrench/pls/pkg/gui"

// CreateGitRepoDropdown prompts the user for selecting from a dropdown of repositories
func CreateGitRepoDropdown(repositories []*github.Repository) *github.Repository {
	names := []string{}
	nameMap := make(map[string]*github.Repository)

	for _, r := range repositories {
		names = append(names, r.GetName())
		nameMap[r.GetName()] = r
	}

	choice := gui.SelectPromptWithResponse("which repository?", names)
	return nameMap[choice]
}

// ChooseWhatToDoWithRepo lets the user decide what to do with their chosen repo
func ChooseWhatToDoWithRepo(gc *github.Client, repository *github.Repository) error {

	// open in browser, clone

	return nil
}
