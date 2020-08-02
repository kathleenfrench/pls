package git

import (
	"errors"

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
