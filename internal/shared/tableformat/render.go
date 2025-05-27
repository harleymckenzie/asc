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

	t.SetTitle(opts.Title)
	t.SetStyle(TableStyles[opts.Style])

	td.WriteRows(t, opts.Layout)
	t.SetColumnConfigs(td.ColumnConfigs())

	//REVIEW - Set column widths?
	// 	columnConfigs := make([]table.ColumnConfig, colsPerRow)
	// 	for i := range colsPerRow {
	// 		columnConfigs[i] = table.ColumnConfig{
	// 			Number:   i + 1,
	// 			WidthMin: 20,
	// 			// WidthMax: 20,
	// 		}
	// 	}
	// 	t.SetColumnConfigs(columnConfigs)

	t.Render()

	return nil
}
