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
