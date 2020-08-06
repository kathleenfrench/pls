package gitpls

import (
	"context"
	"errors"
	"fmt"
	"strconv"

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

	// add exit option last
	opts = append(opts, exitSelections)
	selected := gui.SelectPromptWithResponse(fmt.Sprintf("what would you like to do with %s?", meta.DisplayName), opts, nil, true)

	switch selected {
	case readBodyText:
		render := fmt.Sprintf("# %s\n\n%s", title, body)
		fmt.Println(gui.RenderMarkdown(render))
		return nextOpts(gc, issue, meta, settings)
	case editSelection:
		editOpts := []string{editTitle, editBody, editState}
		switch isPullRequest {
		case true:
			if pr.GetDraft() {
				editOpts = append(editOpts, editReadyForReview)
			}

			prEditTarget := gui.SelectPromptWithResponse("what do you want to change?", editOpts, nil, true)
			switch prEditTarget {
			case editTitle:
				updatedTitle := gui.InputPromptWithResponse("what do you want to call this PR?", title, true)
				pr.Title = &updatedTitle
			case editReadyForReview:
				noDraft := false
				pr.Draft = &noDraft
			case editBody:
				editorCmd := utils.EditorLaunchCommands[settings.DefaultEditor]
				updatedBody := gui.TextEditorInputAndSave("make updates to your PR body", body, editorCmd)
				pr.Body = &updatedBody
			case editState:
				changedState := changeState(state, issue.GetTitle())
				if changedState == nil {
					return nextOpts(gc, issue, meta, settings)
				}

				pr.State = changedState
			}

			gui.Spin.Start()
			updatedPR, _, err := gc.PullRequests.Edit(ctx, meta.Owner, meta.Repo, pr.GetNumber(), pr)
			gui.Spin.Stop()
			if err != nil {
				return err
			}

			gui.Log(":+1:", fmt.Sprintf("successfully updated your PR %q", updatedPR.GetTitle()), updatedPR.GetNumber())
			pr = updatedPR
		default:
			ir := &github.IssueRequest{}
			// edit issue
			editTarget := gui.SelectPromptWithResponse("what do you want to change?", editOpts, nil, true)
			switch editTarget {
			case editTitle:
				updatedTitle := gui.InputPromptWithResponse("what do you want to call this issue?", title, true)
				ir.Title = &updatedTitle
			case editBody:
				editorCmd := utils.EditorLaunchCommands[settings.DefaultEditor]
				updatedBody := gui.TextEditorInputAndSave("make updates to your issue body", body, editorCmd)
				ir.Body = &updatedBody
			case editState:
				changedState := changeState(state, issue.GetTitle())
				if changedState == nil {
					return nextOpts(gc, issue, meta, settings)
				}

				ir.State = changedState
			}

			updatedIssue, _, err := gc.Issues.Edit(ctx, meta.Owner, meta.Repo, issue.GetNumber(), ir)
			if err != nil {
				return err
			}

			gui.Log(":+1:", fmt.Sprintf("successfully updated your issue %q", updatedIssue.GetTitle()), updatedIssue.GetNumber())
			issue = updatedIssue
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
		gui.PleaseHold("checking you back into master and pulling the latest code", nil)
		err = git.CheckoutMasterAndPull()
		if err != nil {
			return fmt.Errorf("could not checkout and pull latest master - %s", err)
		}

		gui.PleaseHold("removing merged branches", nil)
		err = git.CleanupCurrentBranches()
		if err != nil {
			return err
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

func changeState(state string, title string) *string {
	var (
		untypedClosed = "closed"
		untypedOpen   = "open"
	)

	switch state {
	case stateOpen:
		closeIt := gui.ConfirmPrompt(fmt.Sprintf("are you sure you want to close %q?", title), "", false, true)
		if closeIt {
			return &untypedClosed
		}
	case stateClosed:
		reOpenIt := gui.ConfirmPrompt(fmt.Sprintf("are you sure you want to re-open %q?", title), "", false, true)
		if reOpenIt {
			return &untypedOpen
		}
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
