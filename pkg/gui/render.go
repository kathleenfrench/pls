package gui

import (
	"fmt"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/fatih/color"
)

// RenderMarkdown accepts a markdown formatted string and renders it in the terminal
func RenderMarkdown(body string) string {
	markdown.BlueBgItalic = color.New(color.FgBlue).SprintFunc()
	out := markdown.Render(body, 80, 6)
	return string(fmt.Sprintf("\n%s", out))
}
