package pls

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	gitpls "github.com/kathleenfrench/pls/internal/web/git"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/kathleenfrench/pls/pkg/web/git"
	"github.com/spf13/cobra"
)

// FOCUS ON IDIOMATIC COMMANDS

// commands assume context of: (no arg: current repo) | (in <repo>) | <all of github>

// pls get prs to review (current) | (in <org/repo>) | <everywhere>

// pls get --all my prs
// pls get --all my --work prs

// indicate enterprise by using the --work flag
// pls get my --work prs

// pls get my prs [default: all open prs, includes drafts]
// pls get my --merged prs
// pls get my --closed prs
// pls get my --draft prs
// pls get my --merged prs
// pls get my --pending prs
// pls get my --approved prs
// pls get my --changesneeded prs

// searching text
// pls get my prs where <text to search for> --isin|isinthe|inthe|in <title|body|comments>

// TODO: add --by flag for sorting

// search flag variables
var (
	mergedOnly       bool
	closedOnly       bool
	draftOnly        bool
	pending          bool
	approved         bool
	changesNeeded    bool
	assignedOnly     bool
	forCurrentBranch bool
	locked           bool
	mentionedMe      bool
	searchTarget     []string // <title|body|comments>
	work             bool
)

// ------------------------------------------------------

var myGetSubCmd = &cobra.Command{
	Use:     "my",
	Aliases: []string{"m"},
	Short:   "fetch your stuff specifically",
}

// --------------------------- ORGANIZATIONS

var gitMyOrgs = &cobra.Command{
	Use:     "orgs",
	Aliases: []string{"o", "org", "organization", "organizations"},
	Short:   "interact with your github organizations",
	Example: "pls get my orgs",
	Run: func(cmd *cobra.Command, args []string) {
		orgs, err := gitpls.FetchOrganizations("", plsCfg, work)
		if err != nil {
			utils.ExitWithError(err)
		}

		choice := gitpls.CreateGitOrganizationsDropdown(orgs)
		err = gitpls.ChooseWithToDoWithOrganization(choice, plsCfg, work)
		if err != nil {
			utils.ExitWithError(err)
		}
	},
}

// --------------------------- REPOS
var gitMyRepos = &cobra.Command{
	Use:     "repos",
	Aliases: []string{"r", "repositories", "repo", "repository"},
	Short:   "interact with your github repositories",
	Example: "pls get my repos",
	Run: func(cmd *cobra.Command, args []string) {
		repos, err := gitpls.FetchUserRepos("", plsCfg, work)
		if err != nil {
			utils.ExitWithError(err)
		}

		choice := gitpls.CreateGitRepoDropdown(repos)
		err = gitpls.ChooseWhatToDoWithRepo(choice, plsCfg, work)
		if err != nil {
			utils.ExitWithError(err)
		}
	},
}

// --------------------------- ISSUES
var gitMyIssues = &cobra.Command{
	Use:     "issues",
	Aliases: []string{"i", "issue"},
	Short:   "interact with your issues",
	Example: fmt.Sprintf("[issues in current directory's repository]: pls get my issues\n[issues in a repository you own]: pls get my issues in myrepo\n[issues in another's repository]: pls get my issues in organization/repo\n[issues from all of github]: pls get --all my issues\n[issues on my work account]: pls get my --work issues"),
	Run: func(cmd *cobra.Command, args []string) {
		getterFlags := &gitpls.IssueGetterFlags{
			ClosedOnly:      closedOnly,
			AssignedOnly:    assignedOnly,
			Locked:          locked,
			MetaGetterFlags: &gitpls.MetaGetterFlags{},
			PROnly:          false,
		}

		if work {
			getterFlags.MetaGetterFlags.UseEnterpriseAccount = true
			getterFlags.Author = plsCfg.GitEnterpriseUsername
		} else {
			getterFlags.Author = "@me"
		}

		if !getterFlags.ClosedOnly && !getterFlags.MergedOnly {
			getterFlags.State = "open"
		} else {
			getterFlags.State = "closed"
		}

		if getterFlags.AssignedOnly {
			if getterFlags.MetaGetterFlags.UseEnterpriseAccount {
				getterFlags.Assignee = plsCfg.GitEnterpriseUsername
			} else {
				getterFlags.Assignee = "@me"
			}
		}

		switch len(args) {
		case 0:
			// fetch all check
			if !fetchAll {
				// get for whatever is in the current working directory's repo
				org, err := git.CurrentRepositoryOrganization()
				if err != nil {
					utils.ExitWithError(err)
				}

				repo, err := git.CurrentRepositoryName()
				if err != nil {
					utils.ExitWithError(err)
				}

				getterFlags.Organization = org
				getterFlags.Repository = repo

				isEnterprise, err := git.IsEnterpriseGit()
				if err != nil {
					utils.ExitWithError(err)
				}

				src := "github"
				if isEnterprise {
					gui.PleaseHold("github enterprise repository detected...", nil)
					getterFlags.MetaGetterFlags.UseEnterpriseAccount = true
					getterFlags.Author = plsCfg.GitEnterpriseUsername

					if getterFlags.AssignedOnly {
						getterFlags.Assignee = plsCfg.GitEnterpriseUsername
					}

					src = "github enterprise"
				}

				gui.PleaseHold(fmt.Sprintf("searching %s/%s", org, repo), src)
			}

			gui.Spin.Start()
			gc, issues, err := gitpls.SearchIssues(plsCfg, getterFlags)
			gui.Spin.Stop()
			if err != nil {
				utils.ExitWithError(err)
			}

			if len(issues) == 0 {
				gui.OhNo("no issues found matching that crtieria")
				gui.Exit()
			}

			issue, issueMeta := gitpls.CreateGitIssuesDropdown(issues)
			err = gitpls.ChooseWhatToDoWithIssue(gc, issue, issueMeta, plsCfg)
			if err != nil {
				utils.ExitWithError(err)
			}

			break
		case 1:
			utils.ExitWithError(fmt.Sprintf("%s is not a valid argument", args[0]))
			break
		case 2:
			// pls get my issues in <repo> (owned)
			// pls get my issues in <org>/<repo> (organization/another person's repo)
			in := args[0]
			target := args[1]

			if in != "in" {
				utils.ExitWithError(fmt.Sprintf("%s %s is not a valid input", in, target))
			}

			// determine whether we're searching their @username's repo or another
			if strings.Contains(target, "/") {
				sp := strings.Split(target, "/")
				getterFlags.Organization = sp[0]
				getterFlags.Repository = sp[1]
			} else {
				if work {
					utils.ExitWithError("when working with git enterprise resources, an organization value must be specified")
				}

				getterFlags.Organization = plsCfg.GitUsername
				getterFlags.Repository = target
			}

			src := "github"
			if work {
				src = "github enterprise"
			}

			gui.PleaseHold(fmt.Sprintf("searching %s/%s", getterFlags.Organization, getterFlags.Repository), src)
			gui.Spin.Start()
			gc, issues, err := gitpls.SearchIssues(plsCfg, getterFlags)
			gui.Spin.Stop()
			if err != nil {
				utils.ExitWithError(err)
			}

			if len(issues) == 0 {
				color.HiYellow("no issues found matching that criteria!")
				gui.Exit()
			}

			issue, issueMeta := gitpls.CreateGitIssuesDropdown(issues)
			err = gitpls.ChooseWhatToDoWithIssue(gc, issue, issueMeta, plsCfg)
			if err != nil {
				utils.ExitWithError(err)
			}

			break
		default:
			utils.ExitWithError("invalid input, try running `pls get my issues --help`")
		}
	},
}

// --------------------------- PULL REQUESTS
var gitMyPRs = &cobra.Command{
	Use:     "prs",
	Aliases: []string{"pulls", "pull", "pr"},
	Short:   "interact with your pull requests",
	Example: fmt.Sprintf("[PRs in current directory's repository]: pls get my prs\n[PRs in a repository you own]: pls get my prs in myrepo\n[PRs in another's repository]: pls get my prs in organization/repo\n[PRs from all of github]: pls get --all my prs\n[PRs on my work account]: pls get my --work prs"),
	Run: func(cmd *cobra.Command, args []string) {
		getterFlags := &gitpls.IssueGetterFlags{
			MergedOnly:       mergedOnly,
			ClosedOnly:       closedOnly,
			DraftsOnly:       draftOnly,
			PendingApproval:  pending,
			Approved:         approved,
			ChangesRequested: changesNeeded,
			AssignedOnly:     assignedOnly,
			ForCurrentBranch: forCurrentBranch,
			Locked:           locked,
			MetaGetterFlags:  &gitpls.MetaGetterFlags{},
			PROnly:           true,
		}

		if work {
			getterFlags.MetaGetterFlags.UseEnterpriseAccount = true
			getterFlags.Author = plsCfg.GitEnterpriseUsername
		} else {
			getterFlags.Author = "@me"
		}

		if !getterFlags.ClosedOnly && !getterFlags.MergedOnly {
			getterFlags.State = "open"
		} else {
			getterFlags.State = "closed"
		}

		if getterFlags.AssignedOnly {
			if getterFlags.MetaGetterFlags.UseEnterpriseAccount {
				getterFlags.Assignee = plsCfg.GitEnterpriseUsername
			} else {
				getterFlags.Assignee = "@me"
			}
		}

		if getterFlags.ForCurrentBranch {
			cb, err := git.CurrentBranch()
			if err != nil {
				utils.ExitWithError(err)
			}

			// verify a remote ref exists first
			branchHasPR := git.RemoteRefExists(cb)
			if !branchHasPR {
				// push the branch
				err = git.PushBranchToOrigin(cb)
				if err != nil {
					utils.ExitWithError(fmt.Sprintf("your %s branch does not have an open PR, and there was an error attempting to push it - %s", cb, err))
				}
			}

			getterFlags.CurrentBranch = cb
		}

		switch len(args) {
		case 0:
			// fetch all check
			if !fetchAll {
				// get for whatever is in the current working directory's repo
				org, err := git.CurrentRepositoryOrganization()
				if err != nil {
					utils.ExitWithError(err)
				}

				repo, err := git.CurrentRepositoryName()
				if err != nil {
					utils.ExitWithError(err)
				}

				getterFlags.Organization = org
				getterFlags.Repository = repo

				isEnterprise, err := git.IsEnterpriseGit()
				if err != nil {
					utils.ExitWithError(err)
				}

				src := "github"
				if isEnterprise {
					gui.PleaseHold("github enterprise repository detected...", nil)
					getterFlags.MetaGetterFlags.UseEnterpriseAccount = true
					getterFlags.Author = plsCfg.GitEnterpriseUsername

					if getterFlags.AssignedOnly {
						getterFlags.Assignee = plsCfg.GitEnterpriseUsername
					}

					src = "github enterprise"
				}

				gui.PleaseHold(fmt.Sprintf("searching %s/%s", org, repo), src)
			}

			gui.Spin.Start()
			gc, prs, err := gitpls.SearchIssues(plsCfg, getterFlags)
			gui.Spin.Stop()
			if err != nil {
				utils.ExitWithError(err)
			}

			if len(prs) == 0 {
				gui.OhNo("no PRs found matching that crtieria")
				gui.Exit()
			}

			pr, prMeta := gitpls.CreateGitIssuesDropdown(prs)
			err = gitpls.ChooseWhatToDoWithIssue(gc, pr, prMeta, plsCfg)
			if err != nil {
				utils.ExitWithError(err)
			}

			break
		case 1:
			utils.ExitWithError(fmt.Sprintf("%s is not a valid argument", args[0]))
			break
		case 2:
			// pls get my prs in <repo> (owned)
			// pls get my prs in <org>/<repo> (organization/another person's repo)
			in := args[0]
			target := args[1]

			if in != "in" {
				utils.ExitWithError(fmt.Sprintf("%s %s is not a valid input", in, target))
			}

			// determine whether we're searching their @username's repo or another
			if strings.Contains(target, "/") {
				sp := strings.Split(target, "/")
				getterFlags.Organization = sp[0]
				getterFlags.Repository = sp[1]
			} else {
				if work {
					utils.ExitWithError("when working with git enterprise resources, an organization value must be specified")
				}

				getterFlags.Organization = plsCfg.GitUsername
				getterFlags.Repository = target
			}

			src := "github"
			if work {
				src = "github enterprise"
			}
			gui.PleaseHold(fmt.Sprintf("searching %s/%s", getterFlags.Organization, getterFlags.Repository), src)
			gui.Spin.Start()
			gc, prs, err := gitpls.SearchIssues(plsCfg, getterFlags)
			gui.Spin.Stop()
			if err != nil {
				utils.ExitWithError(err)
			}

			if len(prs) == 0 {
				color.HiYellow("no PRs found matching that criteria!")
				gui.Exit()
			}

			pr, prMeta := gitpls.CreateGitIssuesDropdown(prs)
			err = gitpls.ChooseWhatToDoWithIssue(gc, pr, prMeta, plsCfg)
			if err != nil {
				utils.ExitWithError(err)
			}

			break
		default:
			utils.ExitWithError("invalid input, try running `pls get my prs --help`")
		}
	},
}

var myWherePRSubCmd = &cobra.Command{
	Use:     "where",
	Aliases: []string{"w", "search", "s"},
	Example: "pls get my prs where <text to search>",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// flags, search
	},
}

// ------------------------------------------------------
// INIT
// ------------------------------------------------------

func init() {
	getCmd.AddCommand(myGetSubCmd)

	myGetSubCmd.AddCommand(gitMyOrgs)
	myGetSubCmd.AddCommand(gitMyRepos)
	myGetSubCmd.AddCommand(gitMyPRs)
	myGetSubCmd.AddCommand(gitMyIssues)

	// flags
	myGetSubCmd.PersistentFlags().BoolVarP(&work, "work", "w", false, "[ALL] fetch resources via your github enterprise account")
	myGetSubCmd.PersistentFlags().BoolVarP(&mergedOnly, "merged", "m", false, "[PR] fetch only PRs that have been merged")
	myGetSubCmd.PersistentFlags().BoolVarP(&closedOnly, "closed", "c", false, "[PR|ISSUE] fetch only closed prs or issues")
	myGetSubCmd.PersistentFlags().BoolVarP(&pending, "pending", "p", false, "[PR] get only PRs that are pending approval")
	myGetSubCmd.PersistentFlags().BoolVarP(&approved, "approved", "a", false, "[PR] fetch only PRs that have been approved")
	myGetSubCmd.PersistentFlags().BoolVarP(&changesNeeded, "changesneeded", "x", false, "[PR] fetch only PRs where changes have been requested")
	myGetSubCmd.PersistentFlags().BoolVar(&assignedOnly, "assigned", false, "[PR|ISSUE] fetch only PRs or issues assigned to you")
	myGetSubCmd.PersistentFlags().BoolVarP(&forCurrentBranch, "current", "b", false, "[PR] fetch the PR, if one exists, for your current working branch")
	myGetSubCmd.PersistentFlags().BoolVarP(&locked, "locked", "l", false, "[PR|ISSUE] fetch only locked PRs or issues")
	myGetSubCmd.PersistentFlags().BoolVar(&mentionedMe, "mention", false, "[PR|ISSUE] fetch only PRs or issues where i've been mentioned")
	myGetSubCmd.PersistentFlags().BoolVarP(&draftOnly, "draft", "d", false, "[PR] fetch only draft PRs")

	gitMyPRs.AddCommand(myWherePRSubCmd)
	// flags
	myWherePRSubCmd.Flags().StringArrayVarP(&searchTarget, "in", "n", []string{}, fmt.Sprintf("search PRs by text\npls get my prs where add git integration --isin title\n--isin values are title, body, or comments"))
}
