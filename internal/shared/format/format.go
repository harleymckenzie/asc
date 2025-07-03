package format

import (
	"strconv"
	"time"
)

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

// Float64ToStringOrEmpty safely converts a *float64 to string, returns "" if nil.
func Float64ToStringOrEmpty(f *float64) string {
	if f == nil {
		return ""
	}
	return strconv.FormatFloat(*f, 'f', -1, 64)
}

// TimeToStringOrEmpty formats a *time.Time or returns "" if nil.
func TimeToStringOrEmpty(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Local().Format("2006-01-02 15:04:05 MST")
}
