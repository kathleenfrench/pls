package utils

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

// ExecuteCommand runs an external shell command
func ExecuteCommand(cmd *exec.Cmd) string {
	var out string

	if cmd == nil {
		return out
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		out = color.HiRedString(fmt.Sprintf("[ERROR]: %s", err))
		return out
	}

	err = cmd.Start()
	if err != nil {
		out = color.HiRedString(fmt.Sprintf("[ERROR]: %s", err))
		return out
	}

	res, err := ioutil.ReadAll(stdout)
	if err != nil {
		out = color.HiRedString(fmt.Sprintf("[ERROR]: %s", err))
		return out
	}

	out += string(res)

	err = cmd.Wait()
	if err != nil {
		out = color.HiRedString(fmt.Sprintf("[ERROR]: %s", err))
		return out
	}

	return strings.TrimSpace(out)
}

// BashExec executes a command appended to a 'bash -c' command
func BashExec(cmd string) (string, error) {
	r, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}

	trimmed := strings.TrimSpace(string(r))
	return trimmed, nil
}

// OpenURLInDefaultBrowser launches the user system's default browser with a given URL as the target
func OpenURLInDefaultBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		err := exec.Command("open", url).Start()
		if err != nil {
			ExitWithError(err)
		}
	case "linux":
		err := exec.Command("xdg-open", url).Start()
		if err != nil {
			ExitWithError(err)
		}
	default:
		ExitWithError(fmt.Sprintf("%s is not a supported platform", runtime.GOOS))
	}

	return nil
}
