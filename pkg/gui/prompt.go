package gui

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// SelectPromptWithResponse creates a dropdown selection prompt and records the user's choice
func SelectPromptWithResponse(label string, options []string, disableClear bool) string {
	if !disableClear {
		clearScreen()
	}

	var selection string
	var pageSize = len(options)

	if pageSize > 20 {
		pageSize = 20
	}

	prompt := &survey.Select{
		Message:  label,
		Options:  options,
		PageSize: pageSize,
	}

	// see: https://github.com/AlecAivazis/survey/issues/101
	fmt.Printf("\x1b[?7l")
	survey.AskOne(prompt, &selection)
	defer fmt.Printf("\x1b[?7h")

	return selection
}

// ConfirmPrompt prompts the user for a yes/no response to a question, records then returns their response
func ConfirmPrompt(label string, helpMessage string, defaultVal bool) bool {
	clearScreen()

	var response bool
	prompt := &survey.Confirm{
		Message: label,
		Default: defaultVal,
		Help:    helpMessage,
	}

	survey.AskOne(prompt, &response)
	return response
}

// InputPromptWithResponse accepts a user's typed input to a question as a response
func InputPromptWithResponse(label string, defaultVal string) string {
	clearScreen()

	var response string
	prompt := &survey.Input{
		Message: label,
		Default: defaultVal,
		Help:    ":q or :quit to exit",
	}

	// see: https://github.com/AlecAivazis/survey/issues/101
	fmt.Printf("\x1b[?7l")
	survey.AskOne(prompt, &response)
	defer fmt.Printf("\x1b[?7h")

	response = strings.TrimSpace(response)
	if wantsToExit(response) {
		os.Exit(0)
	}

	return response
}

func wantsToExit(v string) bool {
	if v == ":q" || v == ":quit" {
		return true
	}

	return false
}
