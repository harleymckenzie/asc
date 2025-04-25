package tableformat

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func Render(td TableData, list bool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(td.Headers())
	t.AppendRows(td.Rows())
	t.SetColumnConfigs(td.ColumnConfigs())
	t.SortBy(SortBy(td.SortColumns()))
	SetStyle(t, list, false, nil)

	t.Render()
}