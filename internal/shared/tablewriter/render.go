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
	SortBy(fields []Field, reverse bool)
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
	MinColumnWidth int
	MaxColumnWidth int
}

// Field is a single field in a row. It contains a name and a value.
type Field struct {
	Category      string
	Name          string
	Value         string
	Visible       bool
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
	// at.table.SetStyle(StyleRounded)
	at.SetColumnWidth(at.renderOptions.MinColumnWidth, at.renderOptions.MaxColumnWidth)

	if len(at.sortByFields) > 0 {
		sortBy := parseSortBy(at.sortByFields)
		if len(sortBy) > 0 {
			at.table.SortBy(sortBy)
		}
	}

	at.table.Render()
}

// SetColumnWidth sets the minimum and maximum width for all columns.
func (at *AscTable) SetColumnWidth(minWidth int, maxWidth int) {
	configs := make([]table.ColumnConfig, at.renderOptions.Columns)
	for i := 0; i < at.renderOptions.Columns; i++ {
		configs[i] = table.ColumnConfig{Number: i + 1, WidthMin: minWidth, WidthMax: maxWidth}
	}
	at.table.SetColumnConfigs(configs)
}

// GetColumns returns the number of columns in the table
func (at *AscTable) GetColumns() int {
	return at.renderOptions.Columns
}

// SortBy applies sorting based on fields configuration
func (at *AscTable) SortBy(fields []Field, reverse bool) {
	at.sortByFields = fields
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
				Name: field.Name,
				Mode: table.SortMode(field.SortDirection),
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
