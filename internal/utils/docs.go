package utils

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// FrontMatter creates the frontmatter for pls docs
func FrontMatter(filename string, summary string) string {
	const fmTemplate = `---
title: "%s"
slug: %s
url: %s
summary: "%s"
---
`

	name := filepath.Base(filename)
	base := strings.TrimSuffix(name, path.Ext(name))
	url := "/commands/" + strings.ToLower(base) + "/"
	front := fmt.Sprintf(fmTemplate, strings.Replace(base, "_", " ", -1), base, url, summary)
	return front
}

// GenMarkdownDocumentation generates markdown of pls commands into a chosen directory
func GenMarkdownDocumentation(cmd *cobra.Command, dir string, filePrefix func(string, string) string, linkHandler func(string) string) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}

		if err := GenMarkdownDocumentation(c, dir, filePrefix, linkHandler); err != nil {
			return err
		}
	}

	base := fmt.Sprintf("%s.md", strings.Replace(cmd.CommandPath(), " ", "_", -1))
	file := filepath.Join(dir, base)
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := io.WriteString(f, filePrefix(file, cmd.Short)); err != nil {
		return err
	}

	if err := doc.GenMarkdownCustom(cmd, f, linkHandler); err != nil {
		return err
	}

	return nil
}
