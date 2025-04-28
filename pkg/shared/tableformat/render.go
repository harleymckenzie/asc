package tableformat

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func Render(td TableData, sortBy string, list bool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	headers := td.Headers()
	rows := td.Rows()
	headers, rows = RemoveEmptyColumns(headers, rows)	

	t.AppendHeader(headers)
	t.AppendRows(rows)
	t.SetColumnConfigs(td.ColumnConfigs())
	t.SortBy([]table.SortBy{
		{Name: sortBy, Mode: table.Asc},
	})

	// If the list flag is set, use a list style output
	if list {
		t.SetStyle(table.StyleRounded)
		t.Style().Options.DrawBorder = false
		t.Style().Options.SeparateColumns = false
		t.Style().Options.SeparateHeader = false
	} else {
		t.SetStyle(td.TableStyle())
	}

	t.Render()
}
