package pls

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func stylizePls() string {
	cobra.AddTemplateFunc("StyleHeading", color.New(color.FgGreen).SprintFunc())
	usageTemplate := rootCmd.UsageTemplate()
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

func printBanner() {
	color.HiRed(`
██████╗ ██╗     ███████╗
██╔══██╗██║     ██╔════╝
██████╔╝██║     ███████╗
██╔═══╝ ██║     ╚════██║
██║     ███████╗███████║
╚═╝     ╚══════╝╚══════╝											
`)
	printVersion()
}
