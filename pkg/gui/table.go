package gui

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

// RenderTable outputs data into a table
func RenderTable(header table.Row, rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(header)
	t.AppendRows(rows)
	t.SetStyle(table.StyleColoredBright)
	t.Render()
}

// SideBySideTable is used for displaying information
func SideBySideTable(rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", ""})
	t.AppendRows(rows)
	t.SetStyle(table.StyleColoredBright)
	t.Style().Color.Header = text.Colors{text.BgHiRed, text.FgHiRed}
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:     "#",
			Align:    text.AlignLeft,
			Colors:   text.Colors{text.FgHiRed},
			WidthMax: 60,
		},
	})

	t.Render()
}
