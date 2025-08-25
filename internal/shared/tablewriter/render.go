package tablewriter

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// AscWriter is the interface for the AscTable.
type AscWriter interface {
	AppendRow(row Row)
	AppendHeaderRow(headers []string)
	AppendTitleRow(title string)
	AppendHorizontalRow(hr HorizontalRow)
	AppendHorizontalRows([]HorizontalRow)
	AppendGridRow(ar GridRow)
	AppendGridRows([]GridRow)
	Render()
	SetColumnWidth(minWidth int, maxWidth int)
	GetColumns() int
}

// AscTable is the implementation of the AscWriter interface.
type AscTable struct {
	table         table.Writer
	renderOptions AscTableRenderOptions
}

// AscTableRenderOptions is the options for the AscTable.
type AscTableRenderOptions struct {
	Title          string
	Style          string
	Columns        int
	MinColumnWidth int
	MaxColumnWidth int
}

type Row struct {
	Values []string
}

// Field is a single field in a row. It contains a name and a value.
type Field struct {
	Name  string
	Value string
}

// GridRow is made up of two go-pretty table.Row objects.
// The first row is one or more field names, and the second row is the value(s) for each field.
type GridRow struct {
	Fields []Field
}

// HorizontalRow is made up of a single field. The field name is the first column, and the value is the rest of the columns.
type HorizontalRow struct {
	Field Field
}

// NewAscWriter creates a new AscWriter with the provided number of columns.
func NewAscWriter(renderOptions AscTableRenderOptions) AscWriter {
	return &AscTable{
		table:         table.NewWriter(),
		renderOptions: renderOptions,
	}
}

// AppendRow creates a standard row with the provided values.
func (at *AscTable) AppendRow(row Row) {
	rowValues := make(table.Row, len(row.Values))
	for i := 0; i < len(row.Values); i++ {
		rowValues[i] = text.Colors{}.Sprint(row.Values[i])
	}
	at.table.AppendRow(rowValues)
}

// AppendHeader creates a header row that will be formatted according to the table style.
//
//	┌───────────────────────┬─────────────────────────┬───────────────────────┐
//	│ FirstName             │ LastName                │ Age                   │
//	├───────────────────────┼─────────────────────────┼───────────────────────┤
func (at *AscTable) AppendHeaderRow(headers []string) {
	headerRow := make(table.Row, len(headers))
	for i, header := range headers {
		headerRow[i] = header
	}
	at.table.AppendHeader(headerRow)
}

// AppendGridRow creates a new grid row with the provided fields and values.
//
//	┌───────────────────────┬─────────────────────────┬───────────────────────┐
//	│ FirstName             │ LastName                │ Age                   │
//	├───────────────────────┼─────────────────────────┼───────────────────────┤
//	│ John                  │ Doe                     │ 30                    │
//	├───────────────────────┼─────────────────────────┼───────────────────────┤
func (at *AscTable) AppendGridRow(ar GridRow) {
	nr := make(table.Row, at.renderOptions.Columns)
	vr := make(table.Row, at.renderOptions.Columns)

	for i := 0; i < at.renderOptions.Columns; i++ {
		if i < len(ar.Fields) {
			nr[i] = text.Colors{text.Bold, text.FgBlue}.Sprint(ar.Fields[i].Name)
			vr[i] = ar.Fields[i].Value
		} else {
			nr[i] = ""
			vr[i] = ""
		}
	}

	at.table.AppendRow(nr)
	at.table.AppendRow(vr)
	at.table.AppendSeparator()
}

// AppendGridRows accepts a list of GridRows and creates new grid rows for each.
func (at *AscTable) AppendGridRows(ar []GridRow) {
	for _, row := range ar {
		at.AppendGridRow(row)
	}
}

// AppendHeaderRow creates a new row that is made up of the provided title * at.columns times.
//
//	┌─────────────────────────────────────────────────────────────────────────┐
//	│ Title                                                                   │
//	╰─────────────────────────────────────────────────────────────────────────╯
func (at *AscTable) AppendTitleRow(title string) {
	row := make(table.Row, at.renderOptions.Columns)
	for i := 0; i < at.renderOptions.Columns; i++ {
		row[i] = text.Colors{text.Bold}.Sprint(title)
	}

	at.table.AppendSeparator()
	at.table.AppendRow(row, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignLeft})
	at.table.AppendSeparator()
}

// AppendHorizontalRow creates a new row that is made up of the provided name (column 1)
// and the value (all remaining columns, merged).
//
//	┌───────────────────────┬─────────────────────────────────────────────────┐
//	│ Name                  │ John Doe                                        │
//	╰───────────────────────┴─────────────────────────────────────────────────╯
func (at *AscTable) AppendHorizontalRow(hr HorizontalRow) {
	row := make(table.Row, at.renderOptions.Columns)
	row[0] = text.Colors{text.Bold, text.FgBlue}.Sprint(hr.Field.Name)
	for i := 1; i < at.renderOptions.Columns; i++ {
		row[i] = text.Colors{}.Sprint(hr.Field.Value)
	}
	at.table.AppendRow(row, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignLeft})
}

// AppendHorizontalRow accepts a list of Rows and creates new horizontal rows for each.
func (at *AscTable) AppendHorizontalRows(hr []HorizontalRow) {
	for _, row := range hr {
		at.AppendHorizontalRow(row)
	}
}

// Render writes the table to the console.
func (at *AscTable) Render() {
	at.table.SetOutputMirror(os.Stdout)
	at.table.SetTitle(text.Colors{text.Bold}.Sprint(at.renderOptions.Title))
	at.table.SetStyle(table.StyleRounded)
	at.table.Style().Format.Header = text.FormatUpper
	at.SetColumnWidth(at.renderOptions.MinColumnWidth, at.renderOptions.MaxColumnWidth)
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
