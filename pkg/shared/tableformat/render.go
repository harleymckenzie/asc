package tableformat

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func Render(td TableData) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(td.Headers())
	t.AppendRows(td.Rows())
	t.SetColumnConfigs(td.ColumnConfigs())
	t.SortBy(SortBy(td.SortColumns()))
	t.SetStyle(td.TableStyle())

	t.Render()
}