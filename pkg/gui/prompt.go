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
