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

type Tag struct {
	Name  string
	Value string
}

type TagValueGetter func(tag string, instance any) (string, error)

type TagsGetter func(instance any) (map[string]string, error)

type DetailTableLayout struct {
	Type           string // "horizontal"/"grid" or "vertical"
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
	WriteRows(t table.Writer, layout DetailTableLayout)
	ColumnConfigs() []table.ColumnConfig
}
