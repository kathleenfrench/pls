package utils

import (
	"fmt"
	"io/ioutil"
	"os/exec"

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

	return out
}

// BashExec executes a command appended to a 'bash -c' command
func BashExec(cmd string) (string, error) {
	r, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}

	return string(r), nil
}
