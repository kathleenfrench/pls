package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/kyokomi/emoji"
)

// CheckForGitUsername checks for a git username set in the expected location
func CheckForGitUsername() (string, error) {
	username, err := utils.BashExec("git config user.name")
	if err != nil {
		return "", err
	}

	if username == "" {
		return "", errors.New("no github username set")
	}

	return username, nil
}

// PushBranchToOrigin pushes the current working directory's branch to origin
func PushBranchToOrigin(cb string) (err error) {
	if cb == "" {
		cb, err = CurrentBranch()
		if err != nil {
			return err
		}
	}

	gui.PleaseHold("attempting a push", cb)

	cmd := exec.Command("bash", "-c", fmt.Sprintf("git push -u origin %s", cb))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout)
	return nil

}

// CheckoutMasterAndPull checks out the current branch to master and pulls down the latest
func CheckoutMasterAndPull() error {
	cmd := exec.Command("bash", "-c", "git checkout master && git pull")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout)
	return nil
}

// CloneRepository accepts an ssh url and clones a repository to a specified directory
func CloneRepository(name string, sshURL string, path string) error {
	cloneCmd := fmt.Sprintf("git clone %s", sshURL)
	if path != "" {
		path = strings.TrimSuffix(path, "/")
		cloneCmd = fmt.Sprintf("git clone %s %s/%s", sshURL, path, name)
	}

	cmd := exec.Command("bash", "-c", cloneCmd)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout)
	return nil
}

// CurrentBranch returns the name of the branch in the current working directory
func CurrentBranch() (string, error) {
	inWorkTree, err := utils.BashExec("git rev-parse --is-inside-work-tree")
	if err != nil {
		return "", err
	}

	if inWorkTree != "true" {
		return "", errors.New("no git work tree found")
	}

	currentBranch, err := utils.BashExec("git rev-parse --abbrev-ref HEAD")
	if err != nil {
		return "", err
	}

	return currentBranch, nil
}

// GetCurrentGitBaseURL returns the current base git url
func GetCurrentGitBaseURL() string {
	currentRemoteOriginURL, err := utils.BashExec("git config --local --get remote.origin.url")
	if err != nil {
		return ""
	}

	gitBaseCheck := regexp.MustCompile(`github.*.com`)
	val := gitBaseCheck.FindString(currentRemoteOriginURL)
	return strings.TrimSpace(val)
}

// IsEnterpriseGit is a helper for determining whether or not the active repository is from github.com or an enterprise instance
func IsEnterpriseGit() (bool, error) {
	if GetCurrentGitBaseURL() == "" {
		return false, errors.New("an error occurred parsing your current git configuration")
	}

	if GetCurrentGitBaseURL() != "github.com" {
		return true, nil
	}

	return false, nil
}

// CurrentRepositoryOrganization parses the local git config's remote.origin.url to determine the 'organization' or top-level 'user' of a repository
func CurrentRepositoryOrganization() (string, error) {
	var (
		org      string
		gitSplit string
	)

	currentRemoteOriginURL, err := utils.BashExec("git config --local --get remote.origin.url")
	if err != nil {
		return "", err
	}

	if len(currentRemoteOriginURL) == 0 {
		return "", errors.New("could not fetch the remote origin URL of your current working directory's repository")
	}

	gitBaseCheck := regexp.MustCompile(`github.*.com`)
	val := gitBaseCheck.FindString(currentRemoteOriginURL)

	switch strings.Contains(currentRemoteOriginURL, "https") {
	case true:
		// https, like: https://github.com/kathleenfrench/pls.git
		base := fmt.Sprintf("https://%s/", val)
		gitSplit = strings.Split(currentRemoteOriginURL, base)[1]
	case false:
		// ssh, like: git@github.com:kathleenfrench/pls.git
		base := fmt.Sprintf("git@%s:", val)
		gitSplit = strings.Split(currentRemoteOriginURL, base)[1]
	}

	org = strings.Split(gitSplit, "/")[0]
	return org, nil
}

// CurrentRepositoryName returns the name of the repository of the current working directory from any of its subdirectories
func CurrentRepositoryName() (string, error) {
	currentRepo, err := utils.BashExec("basename -s .git `git config --local --get remote.origin.url`")
	if err != nil {
		return "", fmt.Errorf("%s - you have to be in a git repository to run this", err)
	}

	if len(currentRepo) == 0 {
		return "", errors.New("could not determine the name of this git repository")
	}

	return currentRepo, nil
}

// RemoteRefOfCurrentBranchExists checks whether a remote ref of the current working directory's branch exists
func RemoteRefOfCurrentBranchExists() (bool, error) {
	cb, err := CurrentBranch()
	if err != nil {
		return false, err
	}

	return RemoteRefExists(cb), nil
}

func getStatus() (string, error) {
	status, err := utils.BashExec("git status")
	if err != nil {
		return status, err
	}

	return status, nil
}

func gitAddAll() error {
	_, err := utils.BashExec("git add .")
	if err != nil {
		return err
	}

	gui.Log(":white_check_mark:", "added changes", nil)

	return nil
}

func printCodeChanges() (err error) {
	changes, err := utils.BashExec("git status -b")
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("%s\n%s", color.HiRedString("%s UNTRACKED CHANGES FOUND %s", emoji.Sprint(":police_car_light:"), emoji.Sprint(":police_car_light:")), changes))
	return nil
}

func gitCommit() error {
	commitMessage := gui.InputPromptWithResponse("add a commit message", "", true)
	if commitMessage == "" {
		gui.OhNo("you must enter a commit message to proceed")
		return gitCommit()
	}

	_, err := utils.BashExec(fmt.Sprintf("git commit -am %q", commitMessage))
	if err != nil {
		return err
	}

	return nil
}

const (
	notStaged      = "changes not staged for commit"
	needCommitting = "changes to be committed"
	branchAhead    = "your branch is ahead"
)

// HasUnpushedChangesOrCommits checks whether there are unpushed local commits
func HasUnpushedChangesOrCommits() (bool, error) {
	status, err := getStatus()
	if err != nil {
		return false, err
	}

	normalizeStatus := strings.ToLower(status)

	if strings.Contains(normalizeStatus, notStaged) {
		err := printCodeChanges()
		if err != nil {
			return false, err
		}

		confirmAddCommitPush := gui.ConfirmPrompt("you have un-added changes - want me to add, commit, and push them?", "", true, true)
		if !confirmAddCommitPush {
			return true, nil
		}

		err = gitAddAll()
		if err != nil {
			return false, err
		}

		err = gitCommit()
		if err != nil {
			return false, err
		}

		err = PushBranchToOrigin("")
		if err != nil {
			return false, err
		}

		return HasUnpushedChangesOrCommits()
	}

	if strings.Contains(normalizeStatus, needCommitting) {
		err := printCodeChanges()
		if err != nil {
			return false, err
		}

		confirmCommitPush := gui.ConfirmPrompt("you have un-committed changes - want me to commit and push them?", "", true, true)
		if !confirmCommitPush {
			return true, nil
		}

		err = gitCommit()
		if err != nil {
			return false, err
		}

		err = PushBranchToOrigin("")
		if err != nil {
			return false, err
		}

		return HasUnpushedChangesOrCommits()
	}

	if strings.Contains(normalizeStatus, branchAhead) {
		err := printCodeChanges()
		if err != nil {
			return false, err
		}

		confirmPush := gui.ConfirmPrompt("you have un-pushed commits - want me to push them?", "", true, true)
		if !confirmPush {
			return true, nil
		}

		err = PushBranchToOrigin("")
		if err != nil {
			return false, err
		}

		return HasUnpushedChangesOrCommits()
	}

	return false, nil
}

// RemoteRefExists returns a bool for whether a remote reference to a pull request exists
func RemoteRefExists(ref string) bool {
	check := exec.Command("git", "show-ref", "--verify", "--quiet", fmt.Sprintf("refs/remotes/origin/%s", ref))
	err := check.Run()
	return err == nil
}
