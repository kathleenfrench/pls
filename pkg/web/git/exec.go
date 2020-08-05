package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/utils"
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
	color.HiRed("val: %v", val)

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

// RemoteRefExists returns a bool for whether a remote reference to a pull request exists
func RemoteRefExists(ref string) bool {
	check := exec.Command("git", "show-ref", "--verify", "--quiet", fmt.Sprintf("refs/remotes/origin/%s", ref))
	err := check.Run()
	return err == nil
}
