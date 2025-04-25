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
	SortColumns() []string
	ColumnConfigs() []table.ColumnConfig
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

func SortBy(sortOrder []string) []table.SortBy {
	sortBy := []table.SortBy{}

	if len(sortOrder) == 0 {
		sortOrder = []string{"Identifier"}
	}

	for _, sortField := range sortOrder {
		sortBy = append(sortBy, table.SortBy{Name: sortField, Mode: table.Asc})
	}
	return sortBy
}

func SetStyle(t table.Writer, list bool, separateRows bool, mergeColumn *string) {

	t.SetStyle(table.StyleRounded)
	if list {
		t.Style().Options.DrawBorder = false
		t.Style().Options.SeparateColumns = false
		t.Style().Options.SeparateHeader = false
	} else {
		t.Style().Format.Header = text.FormatTitle
        t.Style().Options.SeparateRows = separateRows
        if mergeColumn != nil {
            t.SetColumnConfigs([]table.ColumnConfig{
                {Name: *mergeColumn, AutoMerge: true},
            })
        }
	}
}