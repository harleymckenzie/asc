// Package format provides utilities for parsing and formatting time strings in various formats.
package format

import (
	"fmt"
	"strings"
	"time"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

// durationMap maps human-readable duration units to their Go time.Duration equivalents.
// For example, "2 hours" will be converted to "2h" for parsing.
var durationMap = map[string]string{
	"hour":    "h",
	"hours":   "h",
	"minute":  "m",
	"minutes": "m",
	"day":     "h24",
	"days":    "h24",
	"week":    "h168",
	"weeks":   "h168",
}

// ParseTime attempts to parse a time string in multiple formats:
// 1. Human-readable duration (e.g., "2 hours", "30 minutes")
// 2. RFC3339 format (e.g., "2024-03-20T15:04:05Z")
// 3. Natural language parsing (e.g., "tomorrow at 3pm", "next week")
//
// Returns the parsed time.Time and any error encountered during parsing.
func ParseTime(timeStr string) (time.Time, error) {
	// First attempt: Parse as human-readable duration
	fields := strings.Fields(timeStr)
	if len(fields) == 2 {
		if unit, ok := durationMap[strings.ToLower(fields[1])]; ok {
			durStr := fields[0] + unit
			if duration, err := time.ParseDuration(durStr); err == nil {
				return time.Now().Add(duration), nil
			}
		}
	}

	// Second attempt: Parse as RFC3339 format
	if parsed, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return parsed, nil
	}

	// Third attempt: Parse using natural language processing
	w := when.New(nil)
	w.Add(en.All...)     // Add English language rules
	w.Add(common.All...) // Add common rules

	parsed, err := w.Parse(timeStr, time.Now())
	if err != nil {
		return time.Time{}, err
	}

	if parsed == nil {
		return time.Time{}, fmt.Errorf("invalid: %s", timeStr)
	}

	return parsed.Time, nil
}
