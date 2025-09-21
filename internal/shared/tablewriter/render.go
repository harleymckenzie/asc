package tablewriter

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// AscWriter is the interface for the AscTable.
type AscWriter interface {
	AppendRow(row Row)
	AppendRows(rows []Row)
	AppendHeader(headers []string)
	AppendTitleRow(title string)
	AppendHorizontalRow(hr HorizontalRow)
	AppendHorizontalRows([]HorizontalRow)
	AppendGridRow(ar GridRow)
	AppendGridRows([]GridRow)
	Render()
	SetColumnWidth(minWidth int, maxWidth int)
	GetColumns() int
	SetStyle(style string)
	SetRenderStyle(style string)
	SetFieldConfigs(fields []Field, reverse bool)
}

// AscTable is the implementation of the AscWriter interface.
type AscTable struct {
	table         table.Writer
	renderOptions AscTableRenderOptions
	sortByFields  []Field
	style         *string // Lazy style initialization
}

// AscTableRenderOptions is the options for the AscTable.
type AscTableRenderOptions struct {
	Title          string
	Style          string
	Columns        int
	ColumnConfigs  []table.ColumnConfig
	MinColumnWidth int
	MaxColumnWidth int
	MergedColumns  []string
}

// Field is a single field in a row. It contains a name and a value.
type Field struct {
	Category      string
	Name          string
	Value         string
	Visible       bool
	Merge         bool
	DefaultSort   bool
	SortBy        bool
	SortDirection SortDirection
}

// SortDirection is the direction of the sort.
type SortDirection table.SortMode

const (
	Asc  SortDirection = SortDirection(table.AscNumericAlpha)
	Desc SortDirection = SortDirection(table.DscNumericAlpha)
)

// NewAscWriter creates a new AscWriter with the provided number of columns.
func NewAscWriter(renderOptions AscTableRenderOptions) AscWriter {
	return &AscTable{
		table:         table.NewWriter(),
		renderOptions: renderOptions,
		style:         nil,
	}
}

// getStyle returns the style, initializing it with default if needed
func (at *AscTable) getStyle() string {
	if at.style == nil {
		// Lazy initialization - set default if none specified
		if at.renderOptions.Style == "" {
			defaultStyle := "rounded"
			at.style = &defaultStyle
		} else {
			at.style = &at.renderOptions.Style
		}
	}
	return *at.style
}

// SetRenderStyle sets the style for rendering
func (at *AscTable) SetRenderStyle(style string) {
	at.renderOptions.Style = style
	// Reset the cached style so it gets re-initialized
	at.style = nil
}

// Render writes the table to the console.
func (at *AscTable) Render() {
	at.table.SetOutputMirror(os.Stdout)
	at.table.SetTitle(text.Colors{text.Bold}.Sprint(at.renderOptions.Title))
	at.table.SetStyle(TableStyles[at.getStyle()])
	at.SetColumnWidth(at.renderOptions.MinColumnWidth, at.renderOptions.MaxColumnWidth)
	at.table.SetColumnConfigs(at.renderOptions.ColumnConfigs)

	if len(at.sortByFields) > 0 {
		sortBy := parseSortBy(at.sortByFields)
		if len(sortBy) > 0 {
			at.table.SortBy(sortBy)
		}
	}

	at.table.Render()
}

// GetColumns returns the number of columns in the table
func (at *AscTable) GetColumns() int {
	return at.renderOptions.Columns
}

// SetFieldConfigs sets the field configurations for the table.
func (at *AscTable) SetFieldConfigs(fields []Field, reverse bool) {
	for _, field := range fields {
		if field.Merge {
			at.renderOptions.ColumnConfigs = append(at.renderOptions.ColumnConfigs, table.ColumnConfig{Name: field.Name, AutoMerge: true})
		}
		if field.SortBy {
			at.sortByFields = append(at.sortByFields, field)
		}
	}
	if reverse {
		for i := 0; i < len(at.sortByFields)/2; i++ {
			at.sortByFields[i].SortDirection = reverseSortDirection(at.sortByFields[i].SortDirection)
		}
	}
}

// parseSortBy converts fields with SortBy=true to table.SortBy
func parseSortBy(fields []Field) []table.SortBy {
	var sortBy []table.SortBy
	for _, field := range fields {
		if field.SortBy {
			sortBy = append(sortBy, table.SortBy{
				Name:       field.Name,
				Mode:       table.SortMode(field.SortDirection),
				IgnoreCase: true,
			})
		}
	}
	// If no sort fields found, use the first field with DefaultSort=true
	if len(sortBy) == 0 && len(fields) > 0 {
		sortBy = parseDefaultSort(fields)
	}
	return sortBy
}

// parseDefaultSort sets the DefaultSort field to true if it is not already set.
func parseDefaultSort(fields []Field) []table.SortBy {
	var sortBy []table.SortBy
	for _, field := range fields {
		if field.DefaultSort {
			sortBy = append(sortBy, table.SortBy{
				Name:       field.Name,
				Mode:       table.SortMode(field.SortDirection),
				IgnoreCase: true,
			})
		}
	}
	return sortBy
}

// reverseSortDirection reverses the sort direction.
func reverseSortDirection(sortDirection SortDirection) SortDirection {
	if sortDirection == Asc {
		return Desc
	}
	return Asc
}
