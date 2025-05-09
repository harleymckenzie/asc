package tableformat

import (
	"github.com/jedib0t/go-pretty/v6/table"
)

// TableData provides the data needed to render a generic table.
type TableData interface {
	Headers() table.Row
	Rows() []table.Row
	ColumnConfigs() []table.ColumnConfig
	TableStyle() table.Style
}

type TableDataDetail interface {
	Headers() table.Row
	AppendRows(t table.Writer)
	ColumnConfigs() []table.ColumnConfig
	TableStyle() table.Style
}

type Column struct {
	ID          string
	Visible     bool
	Sort        bool
	DefaultSort bool
}

type Field struct {
	ID          string
	Visible     bool
	Sort        bool
	DefaultSort bool
	Header      bool
}

// BuildHeaders returns a table.Row of column headers
func BuildHeaders(columns []string) table.Row {
	headers := table.Row{}
	for _, col := range columns {
		headers = append(headers, col)
	}
	return headers
}

func BuildColumns(columns []Column) ([]string, string) {
	columnIDs := []string{}
	sortBy := ""
	for _, col := range columns {
		if col.Visible {
			columnIDs = append(columnIDs, col.ID)
		}
		if col.Sort {
			sortBy = col.ID
		}
		if col.DefaultSort {
			sortBy = col.ID
		}
	}
	if sortBy == "" {
		sortBy = columnIDs[0]
	}
	return columnIDs, sortBy
}

func RemoveEmptyColumns(header table.Row, rows []table.Row) (table.Row, []table.Row) {
	if len(header) == 0 || len(rows) == 0 {
		return header, rows
	}

	// Track which columns have non-empty values
	hasValues := make([]bool, len(header))
	for _, row := range rows {
		for colIdx, value := range row {
			if str, ok := value.(string); ok && str != "" {
				hasValues[colIdx] = true
			}
		}
	}

	// Create new header and rows with only non-empty columns
	newHeader := table.Row{}
	for colIdx, value := range header {
		if hasValues[colIdx] {
			newHeader = append(newHeader, value)
		}
	}

	newRows := make([]table.Row, len(rows))
	for i, row := range rows {
		newRow := table.Row{}
		for colIdx, value := range row {
			if hasValues[colIdx] {
				newRow = append(newRow, value)
			}
		}
		newRows[i] = newRow
	}

	return newHeader, newRows
}

func BuildDetailFields(fields []Field) ([]string, string, []string) {
	fieldIDs := []string{}
	sortBy := ""
	headerFields := []string{}
	for _, field := range fields {
		if field.Visible {
			fieldIDs = append(fieldIDs, field.ID)
		}
		if field.Sort {
			sortBy = field.ID
		}
		if field.DefaultSort {
			sortBy = field.ID
		}
		if field.Header {
			headerFields = append(headerFields, field.ID)
		}
	}
	if sortBy == "" {
		sortBy = fieldIDs[0]
	}
	return fieldIDs, sortBy, headerFields
}
