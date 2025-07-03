package utils

import (
	"strconv"
	"strings"
)

// ApplyRelativeOrAbsolute accepts a string (input) that may contain a relative or absolute value,
// and a current value. It will return the new value after applying the relative or absolute value.
// If the string is a relative value, it will be applied to the current value
// If the string is an absolute value, it will become the new current value
func ApplyRelativeOrAbsolute(input string, currentValue int32) (int32, error) {
	if strings.HasPrefix(input, "+") || strings.HasPrefix(input, "-") {
		delta, err := strconv.Atoi(input)
		if err != nil {
			return 0, err
		}
		return currentValue + int32(delta), nil
	}
	newValue, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	return int32(newValue), nil
}

// SlicesToAny converts a slice of any type to a slice of any.
func SlicesToAny[T any](slices []T) []any {
	anySlices := make([]any, len(slices))
	for i, v := range slices {
		anySlices[i] = v
	}
	return anySlices
}

// StringPtr returns a pointer to the string, or nil if the string is empty.
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
