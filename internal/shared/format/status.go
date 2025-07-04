package format

import (
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
)

// Status formats the provided state with the appropriate colour.
// If the state is not found in the StateColours map, the state is returned unchanged.
func Status(state string) string {
	if colour, exists := StateColours[strings.ToLower(state)]; exists {
		return colour.Sprint(state)
	}
	return state
}

// StateColours is a map of state to colour.
var StateColours = map[string]text.Color{
	"100%":                        text.FgGreen,
	"active":                      text.FgGreen,
	"allow":                       text.FgGreen,
	"available":                   text.FgGreen,
	"attached":                    text.FgGreen,
	"blackhole":                   text.FgRed,
	"completed":                   text.FgGreen,
	"create_complete":             text.FgGreen,
	"create-complete":             text.FgGreen,
	"create_in_progress":          text.FgYellow,
	"create-in-progress":          text.FgYellow,
	"create_failed":               text.FgRed,
	"create-failed":               text.FgRed,
	"creating":                    text.FgYellow,
	"delete_complete":             text.FgGreen,
	"delete-complete":             text.FgGreen,
	"delete_failed":               text.FgRed,
	"delete-failed":               text.FgRed,
	"deleted":                     text.FgRed,
	"deleting":                    text.FgRed,
	"deny":                        text.FgRed,
	"error":                       text.FgRed,
	"failed":                      text.FgRed,
	"healthy":                     text.FgGreen,
	"import_complete":             text.FgGreen,
	"import_in_progress":          text.FgYellow,
	"import_rollback_complete":    text.FgGreen,
	"import_rollback_failed":      text.FgRed,
	"import_rollback_in_progress": text.FgYellow,
	"in-use":                      text.FgGreen,
	"modify_complete":             text.FgGreen,
	"modify-complete":             text.FgGreen,
	"modify_failed":               text.FgRed,
	"modify-failed":               text.FgRed,
	"modify_in_progress":          text.FgYellow,
	"modify-in-progress":          text.FgYellow,
	"modifying":                   text.FgYellow,
	"pending":                     text.FgYellow,
	"rebooting":                   text.FgYellow,
	"restore_complete":            text.FgGreen,
	"restore-complete":            text.FgGreen,
	"restore_failed":              text.FgRed,
	"restore-failed":              text.FgRed,
	"restore_in_progress":         text.FgYellow,
	"restore-in-progress":         text.FgYellow,
	"rollback_complete":           text.FgRed,
	"rollback-complete":           text.FgRed,
	"rollback_failed":             text.FgRed,
	"rollback-failed":             text.FgRed,
	"rollback_in_progress":        text.FgYellow,
	"rollback-in-progress":        text.FgYellow,
	"running":                     text.FgGreen,
	"shutting-down":               text.FgRed,
	"stopped":                     text.FgRed,
	"stopping":                    text.FgYellow,
	"terminated":                  text.FgRed,
	"update_complete":             text.FgGreen,
	"update_in_progress":          text.FgBlue,
	"update_rollback_complete":    text.FgRed,
	"update_rollback_failed":      text.FgRed,
	"update_rollback_in_progress": text.FgRed,
	"upgrading":                   text.FgYellow,
}
