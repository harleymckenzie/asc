package smoke

import "strings"

// containsSortWarning returns true if the output contains the sort warning message.
func containsSortWarning(output string) bool {
	return strings.Contains(output, "Warning: Multiple sort fields found")
}

// containsAttributeError returns true if the output contains the attribute error message.
func containsAttributeError(output string) bool {
	return strings.Contains(output, "error getting attribute")
}
