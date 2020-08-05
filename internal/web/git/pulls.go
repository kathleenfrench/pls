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

// PullsGetterFlags are evaluated based off of flags/arguments set by the user when searching pull requests of the current user (author:@me)
// see: https://docs.github.com/en/github/searching-for-information-on-github/searching-issues-and-pull-requests
type PullsGetterFlags struct {
	// basics
	Repository   string // repo:USERNAME/REPOSITORY || example: repo:mozilla/shumway matches issues from @mozilla's shumway project
	Organization string // org:ORGNAME | example: org:github matches issues in repositories owned by the GitHub organization
	User         string // user:USERNAME | example: user:defunkt ubuntu matches issues with the word "ubuntu" from repositories owned by @defunkt

	// searching text
	Match       bool   // pls get my prs --match <title|body|comments> [text to search for]
	InTitleText string // in:title | example: warning in:title matches issues with "warning" in their title.
	InBodyText  string // in:body | example: error in:title,body matches issues with "error" in their title or body.
	InComments  string // in:comments | shipit in:comments matches issues mentioning "shipit" in their comments.

	// state

	// inclusions
	ForCurrentBranch bool
	AssignedOnly     bool
	ClosedOnly       bool
	MergedOnly       bool
	DraftsOnly       bool

	// values set given flags used
	State         string
	Assignee      string
	CurrentBranch string // head:HEAD_BRANCH
	Author        string

	// review statuses
	PendingApproval  bool
	Approved         bool
	ChangesRequested bool

	*MetaGetterFlags
}

// FetchUserPullRequestsEverywhere search all of github for user's pull requests
func FetchUserPullRequestsEverywhere(settings config.Settings, getterFlags *PullsGetterFlags) ([]*github.Issue, error) {
	var allPRs []*github.Issue

	opts := github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	// TODO: GETTER FLAGS

	ctx := context.Background()
	gc := git.NewClient(ctx, settings.GitToken)
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

func (g *PullsGetterFlags) constructMyPRSearchQuery() string {
	query := fmt.Sprintf("type:pr state:%s", g.State)

	if g.Author != "" {
		query += fmt.Sprintf(" author:%s", g.Author)
	}

	if g.Organization != "" && g.Repository != "" {
		query += fmt.Sprintf(" repo:%s/%s", g.Organization, g.Repository)
	}

	if g.ForCurrentBranch && g.CurrentBranch != "" {
		query += fmt.Sprintf(" head:%s", g.CurrentBranch)
	}

	if g.DraftsOnly {
		query += " draft:true"
	}

	if g.AssignedOnly {
		query += fmt.Sprintf(" assignee:%s", g.Assignee)
	}

	if g.Approved {
		query += " review:approved"
	}

	if g.ChangesRequested {
		query += " review:changes_requested"
	}

	if g.PendingApproval {
		query += " review:none"
	}

	if g.MergedOnly {
		query += " is:merged"
	}

	color.HiGreen("query: %s", query)
	// query = fmt.Sprintf("author:@me type:pr repo:%s/%s state:%s", getterFlags.Organization, getterFlags.Repository, getterFlags.State)

	return query
}

// FetchPullRequestsFromCWDRepo parses information about your current working directory's git repository and queries git's API for PRs in that repository
func FetchPullRequestsFromCWDRepo(settings config.Settings, getterFlags *PullsGetterFlags) ([]*github.Issue, error) {
	var allPullsInRepo []*github.Issue

	ctx := context.Background()
	gc := git.NewClient(ctx, settings.GitToken)
	query := fmt.Sprintf("author:@me type:pr repo:%s/%s state:%s", getterFlags.Organization, getterFlags.Repository, getterFlags.State)

	// state defined
	// 		endpoint = fmt.Sprintf("/search/issues?q=author:%s+type:pr+repo:%s/%s+state:%s&sort=updated&order=desc", author, org, repo, state)

	// state not defined
	// 		endpoint = fmt.Sprintf("/search/issues?q=author:%s+type:pr+repo:%s/%s&sort=updated&order=desc", author, org, repo)

	opts := github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		prs, resp, err := gc.Search.Issues(ctx, query, &opts)
		if err != nil {
			return nil, err
		}

		allPullsInRepo = append(allPullsInRepo, prs.Issues...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return allPullsInRepo, nil
}

// CreateGitIssuesDropdown creates a GUI dropdown of issues - PRs are considered an issue per the searchservice in the go-github pkg, so we use this for PR search results as well, returns our custom, readable name
func CreateGitIssuesDropdown(issues []*github.Issue) (*github.Issue, string) {
	names := []string{}
	nameMap := make(map[string]*github.Issue)
	for _, i := range issues {
		r := i.GetRepositoryURL()
		org, repo := git.ExtractOrganizationAndRepoNameFromRepoURL(r)
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
