package tableformat

import (
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
)

// StateColours maps AWS resource states to appropriate colours for table output
var StateColours = map[string]text.Color{
	"active":                      text.FgGreen,
	"available":                   text.FgGreen,
	"create_complete":             text.FgGreen,
	"create_in_progress":          text.FgYellow,
	"creating":                    text.FgYellow,
	"deleted":                     text.FgRed,
	"deleting":                    text.FgRed,
	"failed":                      text.FgRed,
	"healthy":                     text.FgGreen,
	"import_complete":             text.FgGreen,
	"import_in_progress":          text.FgYellow,
	"import_rollback_complete":    text.FgGreen,
	"import_rollback_failed":      text.FgRed,
	"import_rollback_in_progress": text.FgYellow,
	"modifying":                   text.FgYellow,
	"pending":                     text.FgYellow,
	"rebooting":                   text.FgYellow,
	"rollback_complete":           text.FgRed,
	"rollback_failed":             text.FgRed,
	"rollback_in_progress":        text.FgYellow,
	"running":                     text.FgGreen,
	"shutting-down":               text.FgRed,
	"stopped":                     text.FgRed,
	"stopping":                    text.FgYellow,
	"terminated":                  text.FgRed,
	"update_complete":             text.FgGreen,
}

// FormatState formats AWS resource states with appropriate colors for table output
func FormatState(state string) string {
	if colour, exists := StateColours[strings.ToLower(state)]; exists {
		return colour.Sprint(state)
	}
	return state
}
