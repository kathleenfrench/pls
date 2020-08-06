package gitpls

import (
	"context"
	"errors"
	"fmt"
	"strconv"

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
	metas := make(map[string]*IssueMeta)
	for _, i := range issues {
		r := i.GetRepositoryURL()
		org, repo = git.ExtractOrganizationAndRepoNameFromRepoURL(r)
		name := fmt.Sprintf("[%s/%s]: %s", org, repo, i.GetTitle())
		names = append(names, name)
		nameMap[name] = i
		metas[name] = &IssueMeta{
			DisplayName: name,
			Owner:       org,
			Repo:        repo,
			Number:      i.GetNumber(),
		}
	}

	choice := gui.SelectPromptWithResponse("select one", names, nil, false)
	return nameMap[choice], metas[choice]
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
		state   string
	)

	opts := []string{openInBrowser, readBodyText, editSelection}
	ctx := context.Background()

	isPullRequest := issue.IsPullRequest()
	if isPullRequest {
		opts = append(opts, openDiff, mergeSelection)
		htmlURL = issue.GetPullRequestLinks().GetHTMLURL()
		prFetch, _, err := gc.PullRequests.Get(context.Background(), meta.Owner, meta.Repo, meta.Number)
		if err != nil {
			return err
		}

		pr = prFetch
		body = pr.GetBody()
		title = pr.GetTitle()
		state = pr.GetState()
	} else {
		htmlURL = issue.GetHTMLURL()
		body = issue.GetBody()
		title = issue.GetTitle()
		state = issue.GetState()
	}

	if state == stateClosed {
		opts = append(opts, openSelection)
	} else if state == stateOpen {
		opts = append(opts, closeSelection)
	}

	// add exit and close options last
	opts = append(opts, closeSelection, exitSelections)
	selected := gui.SelectPromptWithResponse(fmt.Sprintf("what would you like to do with %s?", meta.DisplayName), opts, nil, true)

	switch selected {
	case readBodyText:
		render := fmt.Sprintf("# %s\n\n%s", title, body)
		fmt.Println(gui.RenderMarkdown(render))
		return nextOpts(gc, issue, meta, settings)
	case editSelection:
		if isPullRequest {
			updatedPR, _, err := gc.PullRequests.Edit(ctx, meta.Owner, meta.Repo, pr.GetNumber(), pr)

			if err != nil {
				return err
			}

			pr = updatedPR
		} else {
			color.HiRed("TODO")
		}
	case mergeSelection:
		if !pr.GetMergeable() || pr.GetMergeableState() != "clean" {
			return errors.New("this PR is currently not in a mergeable state")
		}

		methods := []string{mergeStraight, mergeSquash, mergeRebase}
		mergeMethod := gui.SelectPromptWithResponse("what type of merge do you want to perform?", methods, mergeStraight, true)
		opts := github.PullRequestOptions{
			MergeMethod: mergeMethod,
		}

		result, _, err := gc.PullRequests.Merge(ctx, meta.Owner, meta.Repo, pr.GetNumber(), "", &opts)
		if err != nil {
			return err
		}

		if !result.GetMerged() {
			return errors.New(result.GetMessage())
		}

		gui.Log(":balloon:", result.GetMessage(), result.GetSHA())
	case openSelection:
		color.HiRed("TODO")
	case closeSelection:
		color.HiRed("TODO")
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
	selected := gui.SelectPromptWithResponse("what now?", opts, nil, true)

	switch selected {
	case returnToMenu:
		return ChooseWhatToDoWithIssue(gc, issue, meta, settings)
	case exitSelections:
		gui.Exit()
	}

	return nil
}

func collectPullRequestResponses(settings config.Settings, isEnterprise bool) (*github.NewPullRequest, error) {
	var draft bool

	pr := &github.NewPullRequest{}

	currentBranch, err := git.CurrentBranch()
	if err != nil {
		return nil, err
	}

	base := "master"
	title := gui.InputPromptWithResponse("what do you want to call this PR?", "", true)

	// drafts pull requests are available in public repositories with GitHub Free and GitHub Free for organizations, GitHub Pro, and legacy per-repository billing plans, and in public and private repositories with GitHub Team and GitHub Enterprise Cloud
	if isEnterprise {
		draft = gui.ConfirmPrompt("do you want to create this as a draft?", "", true, true)
	}

	issueLinkCheck := gui.ConfirmPrompt("do you want to link this to an existing issue?", "", false, true)
	if issueLinkCheck {
		numAsString := gui.InputPromptWithResponse("what is the issue number?", "do not include #", true)
		num, err := strconv.Atoi(numAsString)
		if err != nil {
			return nil, err
		}

		pr.Issue = &num
	}

	editorCmd := utils.EditorLaunchCommands[settings.DefaultEditor]
	body := gui.TextEditorInputAndSave("enter a description of this PR", "", editorCmd)

	body += fmt.Sprintf("\n\n---\n<sub>:balloon: i opened this PR by saying [`pls`](https://github.com/kathleenfrench/pls)</sub>\n")

	// set the values
	pr.Title = &title
	pr.Base = &base
	pr.Head = &currentBranch
	pr.Draft = &draft
	pr.Body = &body

	return pr, nil
}
