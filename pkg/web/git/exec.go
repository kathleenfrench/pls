package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
