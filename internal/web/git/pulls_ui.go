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

	choice := gui.SelectPromptWithResponse("select one", names, false)

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
func ChooseWhatToDoWithIssue(gc *github.Client, issue *github.Issue, meta *IssueMeta, settings config.Settings) error {
	var (
		htmlURL string
		pr      *github.PullRequest
		body    string
		title   string
	)

	opts := []string{openInBrowser, readBodyText}

	isPullRequest := issue.IsPullRequest()
	if isPullRequest {
		opts = append(opts, openDiff)
		htmlURL = issue.GetPullRequestLinks().GetHTMLURL()
		prFetch, _, err := gc.PullRequests.Get(context.Background(), meta.Owner, meta.Repo, meta.Number)
		if err != nil {
			return err
		}

		pr = prFetch
		body = pr.GetBody()
		title = pr.GetTitle()
	} else {
		htmlURL = issue.GetHTMLURL()
		body = issue.GetBody()
		title = issue.GetTitle()
	}

	// add exit option last
	opts = append(opts, exitSelections)
	selected := gui.SelectPromptWithResponse(fmt.Sprintf("what would you like to do with %s?", meta.DisplayName), opts, false)

	switch selected {
	case readBodyText:
		render := fmt.Sprintf("# %s\n\n%s", title, body)
		fmt.Println(gui.RenderMarkdown(render))
		return nextOpts(gc, issue, meta, settings)
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

func nextOpts(gc *github.Client, issue *github.Issue, meta *IssueMeta, settings config.Settings) error {
	opts := []string{returnToMenu, exitSelections}
	selected := gui.SelectPromptWithResponse("what now?", opts, true)

	switch selected {
	case returnToMenu:
		return ChooseWhatToDoWithIssue(gc, issue, meta, settings)
	case exitSelections:
		gui.Exit()
	}

	return nil
}
