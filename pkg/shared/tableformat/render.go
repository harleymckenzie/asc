package tableformat

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type RenderOptions struct {
	SortBy string
	List   bool
	Title  string
}

func Render(td TableData, opts RenderOptions) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	headers := td.Headers()
	rows := td.Rows()
	headers, rows = RemoveEmptyColumns(headers, rows)	

	t.AppendHeader(headers)
	t.AppendRows(rows)
	t.SetColumnConfigs(td.ColumnConfigs())
	t.SortBy([]table.SortBy{
		{Name: opts.SortBy, Mode: table.Asc},
	})

	// If the list flag is set, use a list style output
	if opts.List {
		t.SetStyle(table.StyleRounded)
		t.Style().Options.DrawBorder = false
		t.Style().Options.SeparateColumns = false
		t.Style().Options.SeparateHeader = false
	} else {
		if opts.Title != "" {
			t.SetTitle(opts.Title)
		}
		t.SetStyle(td.TableStyle())
	}

	t.Render()
}
