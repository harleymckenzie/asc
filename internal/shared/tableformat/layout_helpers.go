package tableformat

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

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
	for len(fieldIDs) < colsPerRow {
		fieldIDs = append(fieldIDs, "")
		values = append(values, "")
	}
	if len(fieldIDs) != colsPerRow || len(values) != colsPerRow {
		panic(fmt.Sprintf("row has %d columns, expected %d", len(fieldIDs), colsPerRow))
	}
	fieldRow := make([]any, len(fieldIDs))
	for j, h := range fieldIDs {
		fieldRow[j] = text.Colors{text.Bold, text.FgBlue}.Sprint(h)
	}
	t.AppendRow(table.Row(fieldRow))
	t.AppendRow(table.Row(values))
	t.AppendSeparator()
}

func appendVerticalRow(t table.Writer, fieldID string, value any) {
	row := table.Row{
		text.Colors{text.Bold, text.FgHiBlue}.Sprint(fieldID),
		text.Colors{text.FgWhite}.Sprint(value),
	}
	t.AppendRow(row)
}
