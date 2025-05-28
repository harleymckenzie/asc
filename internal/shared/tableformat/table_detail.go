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
		values = append(values, text.Colors{text.FgWhite}.Sprint(val))

		if len(fieldIDs) == colsPerRow || i == len(dt.Fields)-1 {
			appendHorizontalRow(t, fieldIDs, values, colsPerRow)
			fieldIDs = nil
			values = nil
		}
	}
}

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

func (dt *DetailTable) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{
		{Number: 1, Colors: text.Colors{text.Bold}},
	}
}

func processDetailTableHeaders(fields []Field) table.Row {
	headers := table.Row{}
	for _, field := range fields {
		headers = append(headers, field.ID)
	}
	return headers
}
