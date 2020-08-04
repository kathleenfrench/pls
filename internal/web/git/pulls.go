package gitpls

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v32/github"
	"github.com/kathleenfrench/pls/internal/config"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/kathleenfrench/pls/pkg/web/git"
)

// PullRequestSearchFlags are optional flags that can be set on a PR search
type PullRequestSearchFlags struct {
	Closed   bool
	Language string
	Draft    bool
}

// FetchUserPullRequestsEverywhere search all of github for user's pull requests
func FetchUserPullRequestsEverywhere(settings config.Settings) ([]*github.Issue, error) {
	var allPRs []*github.Issue

	opts := github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	ctx := context.Background()
	gc := git.NewClient(ctx, settings.GitToken)
	// TODO: add support for closed?
	query := "author:@me type:pr state:open"

	for {
		prs, resp, err := gc.Search.Issues(ctx, query, &opts)
		if err != nil {
			return nil, err
		}

		allPRs = append(allPRs, prs.Issues...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return allPRs, nil
}

func extractOrganizationAndRepoNameFromRepoURL(url string) (organization string, repo string) {
	// example RepositoryURL: "https://api.github.com/repos/counterThreat/chess_app",
	splitAfter := strings.SplitAfter(url, "repos/")
	orgSplit := strings.Split(splitAfter[1], "/")
	return orgSplit[0], orgSplit[1]
}

// CreateGitIssuesDropdown creates a GUI dropdown of issues - PRs are considered an issue per the searchservice in the go-github pkg, so we use this for PR search results as well, returns our custom, readable name
func CreateGitIssuesDropdown(issues []*github.Issue) (*github.Issue, string) {
	names := []string{}
	nameMap := make(map[string]*github.Issue)
	for _, i := range issues {
		r := i.GetRepositoryURL()
		org, repo := extractOrganizationAndRepoNameFromRepoURL(r)
		name := fmt.Sprintf("[%s/%s]: %s", org, repo, i.GetTitle())
		names = append(names, name)
		nameMap[name] = i
	}

	choice := gui.SelectPromptWithResponse("select one", names)
	return nameMap[choice], choice
}

// ChooseWhatToDoWithIssue lets the user decide what to do with their chosen issue
func ChooseWhatToDoWithIssue(issue *github.Issue, issueName string, prSearchResult bool, settings config.Settings) error {
	opts := []string{openInBrowser}

	if issue.IsPullRequest() {
		opts = append(opts, openDiff)
	}

	opts = append(opts, exitSelections)
	selected := gui.SelectPromptWithResponse(fmt.Sprintf("what would you like to do with %s?", issueName), opts)

	switch selected {
	case openInBrowser:
		if issue.IsPullRequest() {
			utils.OpenURLInDefaultBrowser(issue.GetPullRequestLinks().GetHTMLURL())
		} else {
			utils.OpenURLInDefaultBrowser(issue.GetHTMLURL())
		}
	case openDiff:
		utils.OpenURLInDefaultBrowser(fmt.Sprintf("%s/files", issue.PullRequestLinks.GetDiffURL()))
	case exitSelections:
		gui.Exit()
	}

	return nil
}
