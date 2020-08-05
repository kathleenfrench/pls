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

// CreateGitIssuesDropdown creates a GUI dropdown of issues - PRs are considered an issue per the searchservice in the go-github pkg, so we use this for PR search results as well, returns our custom, readable name
func CreateGitIssuesDropdown(issues []*github.Issue) (*github.Issue, *IssueMeta) {
	names := []string{}
	var org string
	var repo string
	nameMap := make(map[string]*github.Issue)
	for _, i := range issues {
		r := i.GetRepositoryURL()
		org, repo = git.ExtractOrganizationAndRepoNameFromRepoURL(r)
		name := fmt.Sprintf("[%s/%s]: %s", org, repo, i.GetTitle())
		names = append(names, name)
		nameMap[name] = i
	}

	choice := gui.SelectPromptWithResponse("select one", names)

	meta := &IssueMeta{
		DisplayName: choice,
		Owner:       org,
		Repo:        repo,
		Number:      nameMap[choice].GetNumber(),
	}

	return nameMap[choice], meta
}

// IssueMeta is a helper for relevant info when making subsequent API calls from the GUI
type IssueMeta struct {
	Owner       string
	Repo        string
	Number      int
	DisplayName string
}

// ChooseWhatToDoWithIssue lets the user decide what to do with their chosen issue
func ChooseWhatToDoWithIssue(issue *github.Issue, meta *IssueMeta, settings config.Settings) error {
	var (
		htmlURL string
		pr      *github.PullRequest
	)

	opts := []string{openInBrowser, readBodyText}
	ctx := context.Background()
	gc := git.NewClient(ctx, settings.GitToken)

	isPullRequest := issue.IsPullRequest()
	if isPullRequest {
		opts = append(opts, openDiff)
		htmlURL = issue.GetPullRequestLinks().GetHTMLURL()
		prFetch, _, err := gc.PullRequests.Get(ctx, meta.Owner, meta.Repo, meta.Number)
		if err != nil {
			return err
		}

		pr = prFetch
	} else {
		htmlURL = issue.GetHTMLURL()
	}

	color.HiBlue("pr: %v\n", pr)

	// add exit option last
	opts = append(opts, exitSelections)
	selected := gui.SelectPromptWithResponse(fmt.Sprintf("what would you like to do with %s?", meta.DisplayName), opts)

	switch selected {
	case readBodyText:
		if isPullRequest {
			color.HiBlue(pr.GetBody())
		} else {
			color.HiBlue(issue.GetBody())
		}
	case openInBrowser:
		if isPullRequest {
			utils.OpenURLInDefaultBrowser(htmlURL)
		} else {
			utils.OpenURLInDefaultBrowser(htmlURL)
		}
	case openDiff:
		utils.OpenURLInDefaultBrowser(fmt.Sprintf("%s/files", issue.PullRequestLinks.GetDiffURL()))
	case exitSelections:
		gui.Exit()
	}

	return nil
}
