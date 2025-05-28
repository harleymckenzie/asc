package tableformat

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

type ListTable struct {
	Instances    []any
	Fields       []Field
	GetAttribute AttributeGetter
}

func (lt *ListTable) WriteHeaders(t table.Writer) {
	if len(lt.Fields) == 0 {
		panic("cannot render table: no fields defined")
	}
	headers := table.Row{}
	for _, field := range lt.Fields {
		if field.Display || field.Hidden {
			headers = append(headers, field.ID)
		}
	}
	t.AppendHeader(headers, table.RowConfig{AutoMerge: true})
}

func (lt *ListTable) WriteRows(t table.Writer) {
	if len(lt.Fields) == 0 {
		panic("cannot render table: no fields defined")
	}
	rows := []table.Row{}
	for _, instance := range lt.Instances {
		row := table.Row{}
		for _, field := range lt.Fields {
			if field.Display || field.Hidden {
				val, err := lt.GetAttribute(field.ID, instance)
				if err != nil {
					val = fmt.Sprintf("[error: %v]", err)
				}
				row = append(row, val)
			}
		}
		rows = append(rows, row)
	}
	t.AppendRows(rows, table.RowConfig{AutoMerge: false})
}

func (lt *ListTable) ColumnConfigs() []table.ColumnConfig {
	columnConfigs := []table.ColumnConfig{}
	for _, field := range lt.Fields {
		if field.Merge {
			columnConfig := table.ColumnConfig{Name: field.ID, AutoMerge: true}
			columnConfigs = append(columnConfigs, columnConfig)
		}
		if field.Hidden {
			columnConfig := table.ColumnConfig{Name: field.ID, Hidden: true}
			columnConfigs = append(columnConfigs, columnConfig)
		}
	}
	return columnConfigs
}
