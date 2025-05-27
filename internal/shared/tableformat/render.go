// Package tableformat: rendering logic for tables.
// This file contains only rendering functions and options.

package tableformat

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type RenderOptions struct {
	Title  string
	Style  string
	Layout DetailTableLayout
	SortBy []table.SortBy
}

type DetailTableLayout struct {
	Type          DetailTableLayoutType
	ColumnsPerRow int
}

type DetailTableLayoutType string

const (
	DetailTableLayoutClassic DetailTableLayoutType = "classic"
	DetailTableLayoutAlt     DetailTableLayoutType = "alt"
)

// RenderTableList renders a list table.
func RenderTableList(tl ListTableRenderable, opts RenderOptions) {
	// Confirm the GetAttribute function is set
	if listTable, ok := tl.(*ListTable); ok {
		if listTable.GetAttribute == nil {
			panic("GetAttribute function is not set")
		}
	}

	t := table.NewWriter()
	style := TableStyles[opts.Style]
	sortBy := opts.SortBy

	t.SetOutputMirror(os.Stdout)
	t.SetTitle(opts.Title)
	t.SetStyle(style)
	tl.WriteHeaders(t)
	tl.WriteRows(t)
	// Only suppress empty columns if there is at least one row of data
	// This requires a type assertion to access Instances
	if listTable, ok := tl.(*ListTable); ok && len(listTable.Instances) > 0 {
		t.SuppressEmptyColumns()
	}
	t.SetColumnConfigs(tl.ColumnConfigs())
	t.SortBy(sortBy)

	t.Render()
}

// RenderDetail renders a detailed table.
func RenderTableDetail(td DetailTableRenderable, opts RenderOptions) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// td.WriteHeaders(t)

	
	t.SetTitle(opts.Title)
	t.SetStyle(TableStyles[opts.Style])
	if opts.Layout.Type == DetailTableLayoutAlt {
		td.WriteAltRows(t, opts.Layout.ColumnsPerRow)
		columnConfigs := make([]table.ColumnConfig, opts.Layout.ColumnsPerRow)
		for i := 0; i < opts.Layout.ColumnsPerRow; i++ {
			columnConfigs[i] = table.ColumnConfig{
				Number: i,
				WidthMin: 20,
				WidthMax: 20,
			}
		}
		t.SetColumnConfigs(columnConfigs)
	} else {
		td.WriteRows(t)
		t.SetColumnConfigs(td.ColumnConfigs())
	}

	t.Render()

	return nil
}

// func RenderTableDetailAlt(td DetailTableRenderable, opts RenderOptions) error {
// 	t := table.NewWriter()
// 	t.SetOutputMirror(os.Stdout)

// 	// td.WriteHeaders(t)
// 	td.WriteRows(t)

// 	t.SetTitle(opts.Title)
// 	t.SetColumnConfigs(td.ColumnConfigs())
// 	t.SetStyle(TableStyles[opts.Style])

// 	t.Style().Format.Header = text.FormatDefault
// 	t.Style().Size.WidthMin = 70
// 	t.Style().Color.Header = text.Colors{text.Bold}
// 	t.Style().Format.Header = text.FormatDefault

// 	t.Render()

// 	return nil
// }
