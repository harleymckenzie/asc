package tablewriter

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// AscWriter is the interface for the AscTable.
type AscWriter interface {
	AppendAttributeRow(ar AttributeRow)
	AppendAttributeRows([]AttributeRow)
	AppendHeaderRow(hr HeaderRow)
	AppendHorizontalRow(hr HorizontalRow)
	AppendHorizontalRows([]HorizontalRow)
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

// Field is a single field in a row. It contains a name and a value.
type Field struct {
	Name  string
	Value string
}

// AttributeRow is made up of two go-pretty table.Row objects.
// The first row is one or more field names, and the second row is the value(s) for each field.
type AttributeRow struct {
	Fields []Field
}

// HeaderRow is made up of a single field. The field name is the first column, and the value is the rest of the columns.
type HeaderRow struct {
	Title string
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

// AppendAttributeRow creates a new attribute row with the provided fields and values.
func (at *AscTable) AppendAttributeRow(ar AttributeRow) {
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

// AppendAttributeRows accepts a list of AttributeRows and creates new attribute rows for each.
func (at *AscTable) AppendAttributeRows(ar []AttributeRow) {
	for _, row := range ar {
		at.AppendAttributeRow(row)
	}
}

// AppendHeaderRow creates a new row that is made up of the provided title * at.columns times.
func (at *AscTable) AppendHeaderRow(hr HeaderRow) {
	row := make(table.Row, at.renderOptions.Columns)
	for i := 0; i < at.renderOptions.Columns; i++ {
		row[i] = hr.Title
	}

	at.table.AppendRow(row, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignLeft})
	at.table.AppendSeparator()
}

// AppendHorizontalRow creates a new row that is made up of the provided name (column 1)
// and the value (all remaining columns, merged)
func (at *AscTable) AppendHorizontalRow(hr HorizontalRow) {
	row := make(table.Row, at.renderOptions.Columns)
	row[0] = text.Colors{text.Bold, text.FgBlue}.Sprint(hr.Field.Name)
	for i := 1; i < at.renderOptions.Columns; i++ {
		row[i] = text.Colors{}.Sprint(hr.Field.Value)
	}
	at.table.AppendRow(row, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignLeft})
	at.table.AppendSeparator()
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
	at.table.SetTitle(at.renderOptions.Title)
	at.table.SetStyle(table.StyleRounded)
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
