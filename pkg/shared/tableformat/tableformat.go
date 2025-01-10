package tableformat

import "github.com/jedib0t/go-pretty/v6/text"

// ResourceState formats AWS resource states with appropriate colors for table output
func ResourceState(state string) string {
    switch state {
    case "running", "available", "active":
        return text.FgGreen.Sprint(state)
    case "stopped", "failed", "deleted":
        return text.FgRed.Sprint(state)
    case "pending", "creating", "stopping", "modifying", "rebooting":
        return text.FgYellow.Sprint(state)
    default:
        return state
    }
}
