package style

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func bannerString() string {
	return color.HiRedString(fmt.Sprintf("\n%s\n", `
██████╗ ██╗     ███████╗
██╔══██╗██║     ██╔════╝
██████╔╝██║     ███████╗
██╔═══╝ ██║     ╚════██║
██║     ███████╗███████║
╚═╝     ╚══════╝╚══════╝	
`))
}

// PrintBanner prints the pls banner to stdout
func PrintBanner() {
	color.HiRed(`
██████╗ ██╗     ███████╗
██╔══██╗██║     ██╔════╝
██████╔╝██║     ███████╗
██╔═══╝ ██║     ╚════██║
██║     ███████╗███████║
╚═╝     ╚══════╝╚══════╝											
`)
}

// MainMenu styles the main menu output
func MainMenu(usageTemplate string) string {
	cobra.AddTemplateFunc("StyleHeading", color.New(color.FgGreen).SprintFunc())
	usageTemplate = strings.NewReplacer(
		`Usage:`, `{{StyleHeading "Usage:"}}`,
		`Aliases:`, `{{StyleHeading "Aliases:"}}`,
		`Available Commands:`, `{{StyleHeading "Available Commands:"}}`,
		`Global Flags:`, `{{StyleHeading "Global Flags:"}}`,
	).Replace(usageTemplate)
	re := regexp.MustCompile(`(?m)^Flags:\s*$`)
	usageTemplate = re.ReplaceAllLiteralString(usageTemplate, `{{StyleHeading "Flags:"}}`)
	return fmt.Sprintf("%s%s", bannerString(), usageTemplate)
}
