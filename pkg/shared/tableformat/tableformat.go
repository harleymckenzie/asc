package tableformat

import (
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// TableData provides the data needed to render a generic table.
type TableData interface {
	Headers() table.Row
	Rows() []table.Row
	ColumnConfigs() []table.ColumnConfig
	TableStyle() table.Style
}

type Column struct {
	ID      string
	Visible bool
	Sort    bool
}

// ResourceState formats AWS resource states with appropriate colors for table output
func ResourceState(state string) string {
	stateLower := strings.ToLower(state)
	switch stateLower {
	case "running", "available", "active", "healthy", "create_complete", "update_complete", "import_complete":
		return text.FgGreen.Sprint(state)
	case "stopped", "failed", "deleting", "deleted", "terminated", "rollback_failed", "rollback_complete", "update_rollback_complete", "create_failed", "delete_in_progress", "delete_failed", "delete_complete", "import_rollback_in_progress", "import_rollback_complete", "import_rollback_failed":
		return text.FgRed.Sprint(state)
	case "pending", "creating", "stopping", "modifying", "rebooting", "create_in_progress", "rollback_in_progress", "update_in_progress", "update_complete_cleanup_in_progress", "update_rollback_in_progress", "update_rollback_complete_cleanup_in_progress", "import_in_progress":
		return text.FgYellow.Sprint(state)
	default:
		return state
	}
}

// BuildHeaders returns a table.Row of column headers
func BuildHeaders(columns []string) table.Row {
	headers := table.Row{}
	for _, col := range columns {
		headers = append(headers, col)
	}
	return headers
}

func BuildColumns(columns []Column) ([]string, string) {
	columnIDs := []string{}
	sortBy := ""
	for _, col := range columns {
		if col.Visible {
			columnIDs = append(columnIDs, col.ID)
		}
		if col.Sort {
			sortBy = col.ID
		}
	}
	if sortBy == "" {
		sortBy = columnIDs[0]
	}
	return columnIDs, sortBy
}

func RemoveEmptyColumns(header table.Row, rows []table.Row) (table.Row, []table.Row) {
	if len(header) == 0 || len(rows) == 0 {
		return header, rows
	}

	// Track which columns have non-empty values
	hasValues := make([]bool, len(header))
	for _, row := range rows {
		for colIdx, value := range row {
			if str, ok := value.(string); ok && str != "" {
				hasValues[colIdx] = true
			}
		}
	}

	// Create new header and rows with only non-empty columns
	newHeader := table.Row{}
	for colIdx, value := range header {
		if hasValues[colIdx] {
			newHeader = append(newHeader, value)
		}
	}

	newRows := make([]table.Row, len(rows))
	for i, row := range rows {
		newRow := table.Row{}
		for colIdx, value := range row {
			if hasValues[colIdx] {
				newRow = append(newRow, value)
			}
		}
		newRows[i] = newRow
	}

	return newHeader, newRows
}
