package smoke

import (
	"regexp"
	"strings"
)

// containsSortWarning returns true if the output contains the sort warning message.
func containsSortWarning(output string) bool {
	return strings.Contains(output, "Warning: Multiple sort fields found")
}

// containsAttributeError returns true if the output contains the attribute error message.
func containsAttributeError(output string) bool {
	return strings.Contains(output, "error getting attribute")
}

// containsMissingAttributeError returns true if the output contains a missing attribute error for any attribute.
func containsMissingAttributeError(output string) bool {
	missingAttrRe := regexp.MustCompile(`\[error: attribute ".*" does not exist\]`)
	return missingAttrRe.MatchString(output)
}
