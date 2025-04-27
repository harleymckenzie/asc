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
	case "running", "available", "active", "healthy":
		return text.FgGreen.Sprint(state)
	case "stopped", "failed", "deleting", "deleted", "terminated":
		return text.FgRed.Sprint(state)
	case "pending", "creating", "stopping", "modifying", "rebooting":
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
