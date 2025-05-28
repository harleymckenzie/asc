package tableformat

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type DetailTable struct {
	Instance     any
	Fields       []Field
	GetAttribute AttributeGetter
}

// WriteHeaders writes the header row for the detail table.
func (dt *DetailTable) WriteHeaders(t table.Writer) {
	if len(dt.Fields) == 0 {
		panic("cannot render table: no fields defined")
	}
	headers := []string{}
	for _, field := range dt.Fields {
		headers = append(headers, field.ID)
	}
	t.AppendHeader(processDetailTableHeaders(dt.Fields))
}

// WriteRows writes the data rows for the detail table, using the specified layout.
func (dt *DetailTable) WriteRows(t table.Writer, layout DetailTableLayout) {
	if len(dt.Fields) == 0 {
		panic("cannot render table: no fields defined")
	}
	switch layout.Type {
	case "vertical":
		dt.writeRowsVertical(t)
	default:
		dt.writeRowsHorizontal(t, layout.ColumnsPerRow)
	}
}

// writeRowsHorizontal writes rows in horizontal layout (fields and values in rows).
func (dt *DetailTable) writeRowsHorizontal(t table.Writer, colsPerRow int) {
	var fieldIDs []string
	var values []any

	for i, field := range dt.Fields {
		if field.Header {
			if len(fieldIDs) > 0 {
				appendHorizontalRow(t, fieldIDs, values, colsPerRow)
				fieldIDs = nil
				values = nil
			}
			appendHeaderRow(t, field.ID, colsPerRow)
			continue
		}

		fieldIDs = append(fieldIDs, field.ID)
		val, err := dt.GetAttribute(field.ID, dt.Instance)
		if err != nil {
			val = fmt.Sprintf("[error: %v]", err)
		}
		if val == "" {
			val = "-"
		}
		values = append(values, text.Colors{}.Sprint(val))

		if len(fieldIDs) == colsPerRow || i == len(dt.Fields)-1 {
			appendHorizontalRow(t, fieldIDs, values, colsPerRow)
			fieldIDs = nil
			values = nil
		}
	}
}

// writeRowsVertical writes rows in vertical layout (field, value per row).
func (dt *DetailTable) writeRowsVertical(t table.Writer) {
	for _, field := range dt.Fields {
		if field.Header {
			t.AppendSeparator()
			appendHeaderRow(t, field.ID, 2)
			continue
		}
		val, err := dt.GetAttribute(field.ID, dt.Instance)
		if err != nil {
			val = fmt.Sprintf("[error: %v]", err)
		}
		if val == "" {
			val = "-"
		}
		appendVerticalRow(t, field.ID, val)
	}
}

// REVIEW: Confirm if this is needed
func (dt *DetailTable) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{}
}

// processDetailTableHeaders returns a table.Row of column headers for the detail table.
func processDetailTableHeaders(fields []Field) table.Row {
	headers := table.Row{}
	for _, field := range fields {
		headers = append(headers, field.ID)
	}
	return headers
}
