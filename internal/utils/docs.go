package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// PublishDocsDirectory is the directory where pls documentation is output
const PublishDocsDirectory = "docs/pages"

// FrontMatter creates the frontmatter for pls docs
func FrontMatter(filename string, summary string) string {
	const fmTemplate = `---
title: "%s"
permalink: %s
url: %s
summary: "%s"
layout: default
---
`

	name := filepath.Base(filename)
	base := strings.TrimSuffix(name, path.Ext(name))
	url := "/pls/" + strings.ToLower(base) + "/"
	front := fmt.Sprintf(fmTemplate, strings.Replace(base, "_", " ", -1), base, url, summary)
	return front
}

// command paths not warranting their own documentation page
var excludedCommandPaths = []string{
	"pls make a",
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

	if !cmd.Runnable() && !cmd.HasSubCommands() {
		return nil
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

	if err := GenMarkdownExtraCustom(cmd, f, linkHandler); err != nil {
		return err
	}

	return nil
}

// GenMarkdownExtraCustom is an implementation of the doc.GenMarkdownCustom method with some more custom inclusions
func GenMarkdownExtraCustom(cmd *cobra.Command, w io.Writer, linkHandler func(string) string) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()
	cmd.DisableAutoGenTag = true

	b := new(bytes.Buffer)

	// command info
	name := cmd.CommandPath()
	shortDesc := cmd.Short
	longDesc := cmd.Long
	aliases := cmd.Aliases
	example := cmd.Example

	b.WriteString(fmt.Sprintf("# %s \n\n---\n", name))

	if len(aliases) > 0 {
		b.WriteString(fmt.Sprintf("**Aliases**: %s\n\n", strings.Join(aliases, ",")))
	}

	if shortDesc != "" {
		b.WriteString(fmt.Sprintf("**TL;DR:** %s\n\n", shortDesc))
	}

	if longDesc != "" {
		b.WriteString("## Description\n\n")
		b.WriteString(fmt.Sprintf("%s\n\n", longDesc))
	}

	if cmd.Runnable() {
		b.WriteString("## Usage:\n\n")

		if example != "" {
			b.WriteString("### Examples\n\n")
			b.WriteString(fmt.Sprintf("```\n%s\n```\n\n", example))
		}
	}

	err := printOptions(b, cmd, name)
	if err != nil {
		return err
	}

	if hasSubCommands(cmd) {
		b.WriteString("### Sub Commands\n\n")
		subs := cmd.Commands()
		sort.Sort(byName(subs))
		for _, sub := range subs {
			if (!sub.IsAvailableCommand() || sub.IsAdditionalHelpTopicCommand()) && sub.Runnable() {
				continue
			}

			subName := fmt.Sprintf("%s %s", name, sub.Name())
			subLink := strings.Replace(fmt.Sprintf("%s.md", subName), " ", "_", -1)
			b.WriteString(fmt.Sprintf("* [%s](%s)\t - %s\n", subName, linkHandler(subLink), sub.Short))
		}

		b.WriteString("\n")
	}

	if hasParentCommand(cmd) {
		parent := cmd.Parent()
		if parent.Runnable() {
			b.WriteString("### See Also\n\n")
			parentName := parent.CommandPath()
			parentLink := strings.Replace(fmt.Sprintf("%s.md", parentName), " ", "_", -1)
			b.WriteString(fmt.Sprintf("* [%s](%s)\t - %s\n", parentName, linkHandler(parentLink), parent.Short))
			cmd.VisitParents(func(c *cobra.Command) {
				if c.DisableAutoGenTag {
					cmd.DisableAutoGenTag = c.DisableAutoGenTag
				}
			})
		}
	}

	_, err = b.WriteTo(w)
	return err
}

func hasSubCommands(cmd *cobra.Command) bool {
	if cmd.HasSubCommands() {
		return true
	}

	return false
}

func hasParentCommand(cmd *cobra.Command) bool {
	if cmd.HasParent() {
		return true
	}

	return false
}

func printOptions(b *bytes.Buffer, cmd *cobra.Command, name string) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(b)
	if flags.HasAvailableFlags() {
		b.WriteString("### Local Flags\n\n```\n")
		flags.PrintDefaults()
		b.WriteString("```\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(b)
	if parentFlags.HasAvailableFlags() {
		b.WriteString("### Inherited Flags\n\n```\n")
		parentFlags.PrintDefaults()
		b.WriteString("```\n")
	}

	return nil
}

type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }
