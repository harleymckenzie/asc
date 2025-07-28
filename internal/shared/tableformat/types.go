package tableformat

import "github.com/jedib0t/go-pretty/v6/table"

type Field struct {
	DefaultSort   bool
	Display       bool
	Header        bool
	ID            string
	Hidden        bool
	Merge         bool
	Sort          bool
	SortDirection string
}

type AttributeGetter func(fieldID string, instance any) (string, error)

type TagGetter func(tag string, instance any) (string, error)

type DetailTableLayout struct {
	Type           string // "horizontal" or "vertical"
	ColumnsPerRow  int
	ColumnMinWidth int
	ColumnMaxWidth int
}

type ListTableRenderable interface {
	WriteHeaders(t table.Writer)
	WriteRows(t table.Writer)
	ColumnConfigs() []table.ColumnConfig
}

type DetailTableRenderable interface {
	WriteHeaders(t table.Writer)
	WriteRows(t table.Writer, layout DetailTableLayout)
	ColumnConfigs() []table.ColumnConfig
}
