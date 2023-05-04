package gui

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// SelectPromptWithResponse creates a dropdown selection prompt and records the user's choice
func SelectPromptWithResponse(label string, options []string, defaultValue interface{}, disableClear bool) string {
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

	if defaultValue != nil {
		prompt.Default = defaultValue
	}

	// see: https://github.com/AlecAivazis/survey/issues/101
	fmt.Printf("\x1b[?7l")
	_ = survey.AskOne(prompt, &selection)
	defer fmt.Printf("\x1b[?7h")

	return selection
}

// ConfirmPrompt prompts the user for a yes/no response to a question, records then returns their response
func ConfirmPrompt(label string, helpMessage string, defaultVal bool, disableClear bool) bool {
	if !disableClear {
		clearScreen()
	}

	var response bool
	prompt := &survey.Confirm{
		Message: label,
		Default: defaultVal,
		Help:    helpMessage,
	}

	_ = survey.AskOne(prompt, &response)
	return response
}

// TextEditorInputAndSave launches a temporary file with a text editor, captures the input on save, and removes the tmp file while closing the editor
func TextEditorInputAndSave(label string, defaultText string, editor string) string {
	var content string

	prompt := &survey.Editor{
		Message:       label,
		Default:       defaultText,
		AppendDefault: true,
		HideDefault:   true,
		Editor:        editor,
		Help:          "enter text, save it, and pls will handle the rest!",
	}

	_ = survey.AskOne(prompt, &content)
	return content
}

// InputPromptWithResponse accepts a user's typed input to a question as a response
func InputPromptWithResponse(label string, defaultVal string, disableClear bool) string {
	if !disableClear {
		clearScreen()
	}

	var response string
	prompt := &survey.Input{
		Message: label,
		Default: defaultVal,
		Help:    ":q or :quit to exit",
	}

	// see: https://github.com/AlecAivazis/survey/issues/101
	fmt.Printf("\x1b[?7l")
	_ = survey.AskOne(prompt, &response)
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
