package tableformat

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

type ListTable struct {
	Instances    []any
	Fields       []Field
	Tags         []string
	GetAttribute AttributeGetter
	GetTagValue  TagValueGetter
}

// WriteHeaders writes the header row for the list table.
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
	if len(lt.Tags) > 0 {
		for _, tag := range lt.Tags {
			headers = append(headers, fmt.Sprintf("Tag: %s", tag))
		}
	}
	if len(headers) == 0 {
		panic("cannot render table: no headers defined")
	}
	t.AppendHeader(headers, table.RowConfig{AutoMerge: true})
}

// WriteRows writes all data rows for the list table.
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
		if len(lt.Tags) > 0 {
			for _, tag := range lt.Tags {
				tagValue, err := lt.GetTagValue(tag, instance)
				if err != nil {
					tagValue = fmt.Sprintf("[error: %v]", err)
				}
				row = append(row, tagValue)
			}
		}
		if len(row) == 0 {
			panic("cannot render table: no data in row")
		}
		rows = append(rows, row)
	}
	t.AppendRows(rows, table.RowConfig{AutoMerge: false})
}

// ColumnConfigs returns the column configuration for the list table.
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
