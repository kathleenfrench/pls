package gitpls

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/v32/github"
	"github.com/kathleenfrench/pls/internal/config"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/web/git"
	"github.vimeows.com/Vimeo/dex/pkg/utils"
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
	Locked           bool
	MentionedMe      bool

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

// FetchPullRequests search all of github for user's pull requests based on search criteria
func FetchPullRequests(settings config.Settings, getterFlags *PullsGetterFlags) (*github.Client, []*github.Issue, error) {
	var allPRs []*github.Issue
	var gc *github.Client
	var err error

	opts := github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	ctx := context.Background()

	if getterFlags.UseEnterpriseAccount {
		gc, err = git.NewEnterpriseClient(ctx, settings.GitEnterpriseHostname, settings.GitEnterpriseToken)
		if err != nil {
			return nil, nil, err
		}
	} else {
		gc = git.NewClient(ctx, settings.GitToken)
	}

	query := getterFlags.constructMyPRSearchQuery()

	for {
		prs, resp, err := gc.Search.Issues(ctx, query, &opts)
		if err != nil {
			return nil, nil, err
		}

		allPRs = append(allPRs, prs.Issues...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return gc, allPRs, nil
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

	if g.Locked {
		query += " is:locked"
	}

	if g.MentionedMe {
		query += " mentions:@me"
	}

	return query
}

// CreatePullRequestFromCWD creates a PR for the branch in your current working directory
func CreatePullRequestFromCWD(settings config.Settings) error {
	var gc *github.Client

	remoteRefExists, err := git.RemoteRefOfCurrentBranchExists()
	if err != nil {
		return err
	}

	if !remoteRefExists {
		return errors.New("no remote ref for your current branch exists - make sure to push it before attempting to create a PR")
	}

	ctx := context.Background()

	org, err := git.CurrentRepositoryOrganization()
	if err != nil {
		return err
	}

	repo, err := git.CurrentRepositoryName()
	if err != nil {
		return err
	}

	isEnterprise, err := git.IsEnterpriseGit()
	if err != nil {
		return err
	}

	newPR, err := collectPullRequestResponses(settings, isEnterprise)
	if err != nil {
		return err
	}

	if isEnterprise {
		gc, err = git.NewEnterpriseClient(ctx, settings.GitEnterpriseHostname, settings.GitEnterpriseToken)
		if err != nil {
			return err
		}
	} else {
		gc = git.NewClient(ctx, settings.GitToken)
	}

	fmt.Println(newPR, org, repo)
	pr, _, err := gc.PullRequests.Create(ctx, org, repo, newPR)
	if err != nil {
		return err
	}

	gui.Log(":party:", fmt.Sprintf("success! your PR %q is now open for business", pr.GetTitle()), fmt.Sprintf("#%d", pr.GetNumber()))

	openNow := gui.ConfirmPrompt("do you want to open your PR in the browser?", "", true, true)
	if openNow {
		utils.LaunchInBrowser(pr.GetHTMLURL())
	}

	return nil
}
