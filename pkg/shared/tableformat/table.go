package tableformat

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type AttributeGetter func(fieldID string, instance any) string

type ListTableRenderable interface {
	WriteHeaders(t table.Writer)
	WriteRows(t table.Writer)
	ColumnConfigs() []table.ColumnConfig
}

type DetailTableRenderable interface {
	WriteHeaders(t table.Writer)
	WriteRows(t table.Writer)
	ColumnConfigs() []table.ColumnConfig
}

// ListTable is a struct that defines a list table.
type ListTable struct {

	// Instances is a list of objects to be displayed
	Instances []any

	// Fields is a list of fields to be displayed
	Fields []Field

	// GetAttribute is a function that returns the attribute of the instance
	GetAttribute AttributeGetter
}

// DetailTable is a struct that defines a detailed table.
type DetailTable struct {

	// Instance is the instance of the object to be displayed
	Instance any

	// Fields is a list of fields to be displayed
	Fields []Field

	// GetAttribute is a function that returns the attribute of the instance
	GetAttribute AttributeGetter
}

// Field is a struct that defines a field in a table.
type Field struct {

	// ID is the unique identifier for the field/column
	ID string

	// Visible determines if the field should be shown in the table
	Visible bool

	// Sort determines if the field can be used for sorting
	Sort bool

	// DefaultSort marks this field as the default sort column
	DefaultSort bool

	// Header marks this field as a header row (for grouping or sectioning)
	Header bool

	// Merge determines if adjacent cells with the same value should be merged
	Merge bool

	// SortDirection specifies the default sort direction ("asc" or "dsc")
	SortDirection string
}

// List Table Methods
func (lt *ListTable) WriteHeaders(t table.Writer) {
	headers := table.Row{}
	for _, field := range lt.Fields {
		if field.Visible {
			headers = append(headers, field.ID)
		}
	}
	t.AppendHeader(headers, table.RowConfig{AutoMerge: true})
}

func (lt *ListTable) WriteRows(t table.Writer) {
	rows := []table.Row{}
	for _, instance := range lt.Instances {
		row := table.Row{}
		for _, field := range lt.Fields {
			if field.Visible {
				row = append(row, lt.GetAttribute(field.ID, instance))
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
	}
	return columnConfigs
}

// Detail Table Methods

func (dt *DetailTable) WriteHeaders(t table.Writer) {
	headers := []string{}
	for _, field := range dt.Fields {
		headers = append(headers, field.ID)
	}
	t.AppendHeader(processDetailTableHeaders(dt.Fields))
}

func (dt *DetailTable) WriteRows(t table.Writer) {
	for _, field := range dt.Fields {
		if field.Header {
			row := table.Row{field.ID, field.ID}
			t.AppendSeparator()
			t.AppendRow(row, table.RowConfig{AutoMerge: true})
			t.AppendSeparator()
		} else {
			val := dt.GetAttribute(field.ID, dt.Instance)
			row := table.Row{field.ID, val}
			t.AppendRow(row)
		}
	}
}

func (dt *DetailTable) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{
		{Number: 1, Colors: text.Colors{text.Bold}},
	}
}

// GetSortByField returns a table.SortBy for the given fields
// It will look at the 'Sort', 'DefaultSort' and 'SortDirection' fields to determine the sort order
// If no sort fields are found, the 'DefaultSort' field will be used
// If no 'DefaultSort' field is found, the first field will be used
func GetSortByField(fields []Field, reverseSort bool) []table.SortBy {
	if len(fields) == 0 {
		return nil
	}

	var sortBy []table.SortBy
	// 1. Find fields with Sort == true
	for _, f := range fields {
		if f.Sort {
			sortBy = append(sortBy, table.SortBy{
				Name:       f.ID,
				Mode:       sortDirection(f.SortDirection, reverseSort),
				IgnoreCase: true,
			})
		}
	}

	// 2. If no sort fields found, find fields with DefaultSort == true
	// If no sort direction is specified, use ascending order (and reverse if specified)
	if len(sortBy) == 0 {
		for _, f := range fields {
			if f.DefaultSort {
				if f.SortDirection == "" {
					sortBy = append(sortBy, table.SortBy{
						Name:       f.ID,
						Mode:       sortDirection("asc", reverseSort),
						IgnoreCase: true,
					})
				} else {
					sortBy = append(sortBy, table.SortBy{
						Name:       f.ID,
						Mode:       sortDirection(f.SortDirection, reverseSort),
						IgnoreCase: true,
					})
				}
			}
		}
	}

	// 3. If still none, use the first field
	if len(sortBy) == 0 {
		sortBy = append(sortBy, table.SortBy{
			Name:       fields[0].ID,
			Mode:       sortDirection("asc", reverseSort),
			IgnoreCase: true,
		})
	}

	if len(sortBy) > 1 {
		fmt.Printf(
			"Warning: Multiple sort fields found: %v\nFields will be sorted by the first field\n",
			sortBy,
		)
	}

	return sortBy
}

// processDetailTableHeaders returns a table.Row of column headers
func processDetailTableHeaders(fields []Field) table.Row {
	headers := table.Row{}
	for _, field := range fields {
		headers = append(headers, field.ID)
	}
	return headers
}

func processDetailTableFields(fields []Field) ([]string, string, []string) {
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

// sortDirection returns the sort mode for the given direction
// If reverseSort is true, the direction will be reversed
func sortDirection(direction string, reverseSort bool) table.SortMode {
	if reverseSort {
		switch direction {
		case "asc":
			return table.DscNumericAlpha
		case "dsc":
			return table.AscNumericAlpha
		}
	}

	if direction == "asc" {
		return table.AscNumericAlpha
	}
	return table.DscNumericAlpha
}
