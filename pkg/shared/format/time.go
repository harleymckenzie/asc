package format

import (
	"fmt"
	"strings"
	"time"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/en"
	"github.com/olebedev/when/rules/common"
)

var durationMap = map[string]string{
    "hour":   "h",
    "hours":  "h",
    "minute": "m",
    "minutes": "m",
    "day":    "h24",
    "days":   "h24",
    "week":   "h168",
    "weeks":  "h168",
}

func ParseTime(timeStr string) (time.Time, error) {
    // Try parsing as human-readable duration (e.g., "2 hours")
    fields := strings.Fields(timeStr)
    if len(fields) == 2 {
        if unit, ok := durationMap[strings.ToLower(fields[1])]; ok {
            durStr := fields[0] + unit
            if duration, err := time.ParseDuration(durStr); err == nil {
                return time.Now().Add(duration), nil
            }
        }
    }

	// Fall back to natural language parsing
	if parsed, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return parsed, nil
	}

	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	parsed, err := w.Parse(timeStr, time.Now())
	if err != nil {
		return time.Time{}, err
	}

	if parsed == nil {
		return time.Time{}, fmt.Errorf("invalid: %s", timeStr)
	}

	return parsed.Time, nil
}