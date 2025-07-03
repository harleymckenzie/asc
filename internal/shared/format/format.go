package format

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// DecodeAndFormatJSON decodes a JSON string and formats it with indentation.
func DecodeAndFormatJSON(encoded *string) string {
    if encoded == nil {
        return ""
    }
    decoded, err := url.QueryUnescape(*encoded)
    if err != nil {
        return *encoded // fallback to raw if decode fails
    }
    // Optional: pretty-print JSON
    var prettyJSON map[string]interface{}
    if err := json.Unmarshal([]byte(decoded), &prettyJSON); err == nil {
        pretty, _ := json.MarshalIndent(prettyJSON, "", "  ")
        return string(pretty)
    }
    return decoded
}

// StringOrEmpt safely dereferences a *string, returning "" if nil.
func StringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// StringOrDefault safely dereferences a *string, returning "" if nil.
func StringOrDefault(s *string, defaultValue string) string {
	if s == nil {
		return defaultValue
	}
	return *s
}

// Int64ToStringOrDefault safely converts *int64 to string, "" if nil.
func Int64ToStringOrDefault(i *int64, defaultValue string) string {
	if i == nil {
		return defaultValue
	}
	return strconv.FormatInt(*i, 10)
}

// IntToStringOrDefault safely converts *int to string, "" if nil.
func IntToStringOrDefault(i *int, defaultValue string) string {
	if i == nil {
		return defaultValue
	}
	return strconv.Itoa(*i)
}

// Int32ToStringOrDefault safely converts *int32 to string, "" if nil.
func Int32ToStringOrDefault(i *int32, defaultValue string) string {
	if i == nil {
		return defaultValue
	}
	return strconv.Itoa(int(*i))
}

// BoolToStringOrDefault safely converts *bool to string, "" if nil.
func BoolToStringOrDefault(b *bool, defaultValue string) string {
	if b == nil {
		return defaultValue
	}
	return strconv.FormatBool(*b)
}

// TimeToStringOrDefault formats *time.Time or returns "" if nil.
func TimeToStringOrDefault(t *time.Time, defaultValue string) string {
	if t == nil {
		return defaultValue
	}
	return t.Local().Format("2006-01-02 15:04:05 MST")
}

// StatusOrDefault returns format.Status(s) or "" if s is empty.
func StatusOrDefault(s string, defaultValue string) string {
	if s == "" {
		return defaultValue
	}
	return Status(s)
}

// BoolToLabel returns trueLabel if b is true, falseLabel if b is false, or "" if b is nil.
func BoolToLabel(b *bool, trueLabel, falseLabel string) string {
	if b == nil {
		return ""
	}
	if *b {
		return trueLabel
	}
	return falseLabel
}

// Int64ToStringOrEmpty safely converts a *int64 to string, returns "" if nil.
func Int64ToStringOrEmpty(i *int64) string {
	if i == nil {
		return ""
	}
	return strconv.FormatInt(int64(*i), 10)
}

// Int32ToStringOrEmpty safely converts a *int32 to string, returns "" if nil.
func Int32ToStringOrEmpty(i *int32) string {
	if i == nil {
		return ""
	}
	return strconv.FormatInt(int64(*i), 10)
}

// TimeToStringOrEmpty formats a *time.Time or returns "" if nil.
func TimeToStringOrEmpty(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Local().Format("2006-01-02 15:04:05 MST")
}

// TimeToStringRelative formats a *time.Time or returns "" if nil. It shows the relative time since the time.
func TimeToStringRelative(t *time.Time) string {
	if t == nil {
		return ""
	}
	now := time.Now()
	diff := now.Sub(*t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	case diff < 30*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	default:
		return t.Local().Format("2006-01-02 15:04:05 MST")
	}
}
