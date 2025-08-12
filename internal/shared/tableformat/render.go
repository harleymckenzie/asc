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

// RenderTableList renders a list table using the provided options.
func RenderTableList(tl ListTableRenderable, opts RenderOptions) {
	t := table.NewWriter()

	// Set default style if not specified
	if opts.Style == "" {
		opts.Style = DefaultTableStyle
	}

	style := TableStyles[opts.Style]
	sortBy := opts.SortBy

	t.SetOutputMirror(os.Stdout)
	t.SetTitle(opts.Title)
	t.SetStyle(style)
	tl.WriteHeaders(t)
	tl.WriteRows(t)
	if listTable, ok := tl.(*ListTable); ok && len(listTable.Instances) > 0 {
		t.SuppressEmptyColumns()
	}
	t.SetColumnConfigs(tl.ColumnConfigs())
	t.SortBy(sortBy)
	t.Render()
}

// RenderTableDetail renders a detailed table using the provided options.
func RenderTableDetail(td DetailTableRenderable, opts RenderOptions) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(opts.Title)
	t.SetStyle(TableStyles[opts.Style])

	// Ensure ColumnsPerRow is set to a sensible default
	if opts.Layout.ColumnsPerRow <= 0 {
		opts.Layout.ColumnsPerRow = 3
	}

	td.WriteRows(t, opts.Layout)
	t.SetColumnConfigs(td.ColumnConfigs())

	// Optionally set column widths if specified in layout
	if (opts.Layout.ColumnMinWidth > 0 || opts.Layout.ColumnMaxWidth > 0) && opts.Layout.ColumnsPerRow > 0 {
		SetColumnWidths(t, opts.Layout.ColumnsPerRow, opts.Layout.ColumnMinWidth, opts.Layout.ColumnMaxWidth)
	}

	if t.Length() == 0 {
		panic("cannot render table: no columns defined (header missing?)")
	}
	t.Render()
	return nil
}

// SetColumnWidths sets the minimum and maximum width for each column in the table.
// This can be used in RenderTableDetail if you want to enforce column widths.
func SetColumnWidths(t table.Writer, cols int, minWidth int, maxWidth int) {
	columnConfigs := make([]table.ColumnConfig, cols)
	for i := 0; i < cols; i++ {
		columnConfigs[i] = table.ColumnConfig{
			Number: i + 1,
		}
		if minWidth > 0 {
			columnConfigs[i].WidthMin = minWidth
		}
		if maxWidth > 0 {
			columnConfigs[i].WidthMax = maxWidth
		}
	}
	t.SetColumnConfigs(columnConfigs)
}
