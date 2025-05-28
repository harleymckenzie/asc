package tableformat

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// appendHeaderRow adds a header row for a section, spanning colsPerRow columns.
func appendHeaderRow(t table.Writer, fieldID string, colsPerRow int) {
	headerIDs := make([]string, colsPerRow)
	for j := 0; j < colsPerRow; j++ {
		headerIDs[j] = fieldID
	}
	headerRow := make([]any, len(headerIDs))
	for j, h := range headerIDs {
		formattedID := text.Colors{text.Bold, text.FgWhite}.Sprint(h)
		headerRow[j] = formattedID
	}
	t.AppendRow(table.Row(headerRow), table.RowConfig{
		AutoMerge:      true,
		AutoMergeAlign: text.AlignLeft,
	})
	t.AppendSeparator()
}

// appendHorizontalRow adds a row of field titles and a row of values, padding to colsPerRow.
func appendHorizontalRow(t table.Writer, fieldIDs []string, values []any, colsPerRow int) {
	if colsPerRow == 0 {
		panic("cannot render table: ColumnsPerRow is 0 (not set in layout?)")
	}
	for len(fieldIDs) < colsPerRow {
		fieldIDs = append(fieldIDs, "")
		values = append(values, "")
	}
	if len(fieldIDs) != colsPerRow || len(values) != colsPerRow {
		panic(fmt.Sprintf("row has %d columns, expected %d", len(fieldIDs), colsPerRow))
	}

	// Process the title headers
	fieldRow := make([]any, len(fieldIDs))
	for j, h := range fieldIDs {
		fieldRow[j] = text.Colors{text.Bold, text.FgBlue}.Sprint(h)
	}

	t.AppendRow(table.Row(fieldRow))
	t.AppendRow(table.Row(values))
	t.AppendSeparator()
}

// appendVerticalRow adds a single field-value row to the table.
func appendVerticalRow(t table.Writer, fieldID string, value any) {
	row := table.Row{
		text.Colors{text.Bold, text.FgHiBlue}.Sprint(fieldID),
		text.Colors{}.Sprint(value),
	}
	t.AppendRow(row)
}
