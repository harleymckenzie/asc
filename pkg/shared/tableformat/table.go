package tableformat

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type TableRenderer interface {
	RenderTable(data interface{}, config TableConfig, options TableOptions) error
}

type TableOptions struct {
	List            bool
	SortOrder       []string
	SelectedColumns []string
}

type BaseTableRenderer struct {
	Writer table.Writer
}

func NewBaseTableRenderer() *BaseTableRenderer {
	return &BaseTableRenderer{
		Writer: table.NewWriter(),
	}
}

func (r *BaseTableRenderer) RenderTable(data interface{}, config TableConfig,
	options TableOptions) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Set headers
	headers := GetSelectedColumns(config, options.SelectedColumns)
	headerRow := make(table.Row, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}
	t.AppendHeader(headerRow)

	// Handle data rows
	rows, ok := data.([]map[string]string)
	if !ok {
		return fmt.Errorf("data must be []map[string]string, got %T", data)
	}

	for _, rowData := range rows {
		row := make(table.Row, len(options.SelectedColumns))
		for i, colID := range options.SelectedColumns {
			row[i] = rowData[colID]
		}
		t.AppendRow(row)
	}

	// Apply column configs
	columnConfigs := make([]table.ColumnConfig, 0)
	for _, col := range config.Columns {
		if col.WidthMin > 0 || col.WidthMax > 0 || col.AutoMerge {
			columnConfigs = append(columnConfigs, table.ColumnConfig{
				Name:      col.Title,
				WidthMin:  col.WidthMin,
				WidthMax:  col.WidthMax,
				AutoMerge: col.AutoMerge,
			})
		}
	}
	if len(columnConfigs) > 0 {
		t.SetColumnConfigs(columnConfigs)
	}

	// Set style and sorting
	SetStyle(t, options.List, config.SeparateRows, config.MergeColumn)
	if len(options.SortOrder) > 0 {
		t.SortBy(SortBy(options.SortOrder))
	}

	t.Render()
	return nil
}
