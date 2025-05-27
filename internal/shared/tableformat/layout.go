package tableformat

import (
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type DetailTableLayout struct {
	Type          string
	ColumnsPerRow int
}

func appendHeaderRow(t table.Writer, fieldID string, colsPerRow int) {
	headerIDs := make([]string, colsPerRow)
	for j := 0; j < colsPerRow; j++ {
		headerIDs[j] = strings.ToUpper(fieldID)
	}
	headerRow := make([]any, len(headerIDs))
	for j, h := range headerIDs {
		formattedID := text.Colors{text.Bold, text.FgBlue}.Sprint(h)
		headerRow[j] = formattedID
	}
	t.AppendRow(table.Row(headerRow), table.RowConfig{
		AutoMerge:      true,
		AutoMergeAlign: text.AlignLeft,
	})
	t.AppendSeparator()
}

func appendHorizontalRow(t table.Writer, fieldIDs []string, values []any, colsPerRow int) {
	// Pad fieldIDs and values to colsPerRow
	for len(fieldIDs) < colsPerRow {
		fieldIDs = append(fieldIDs, "")
		values = append(values, "")
	}
	fieldRow := make([]any, len(fieldIDs))
	for j, h := range fieldIDs {
		fieldRow[j] = text.Colors{text.Bold, text.FgBlue}.Sprint(strings.ToUpper(h))
	}
	t.AppendRow(table.Row(fieldRow))

	// Distribute the values across the columns
	t.AppendRow(table.Row(values))
	t.AppendSeparator()
}

func appendVerticalRow(t table.Writer, fieldID string, value any) {
	row := table.Row{text.Bold.Sprint(fieldID), value}
	t.AppendRow(row)
}