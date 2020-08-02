package gui

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// SelectPromptWithResponse creates a dropdown selection prompt and records the user's choice
func SelectPromptWithResponse(label string, options []string) string {
	var selection string

	prompt := &survey.Select{
		Message:  label,
		Options:  options,
		PageSize: len(options),
	}

	// bug in survey pkg
	fmt.Printf("\x1b[?7l")
	survey.AskOne(prompt, &selection)
	defer fmt.Printf("\x1b[?7h")

	return selection
}

// ConfirmPrompt prompts the user for a yes/no response to a question, records then returns their response
func ConfirmPrompt(label string, helpMessage string, defaultVal bool) bool {
	var response bool
	prompt := &survey.Confirm{
		Message: label,
		Default: defaultVal,
		Help:    helpMessage,
	}

	survey.AskOne(prompt, &response)
	return response
}
