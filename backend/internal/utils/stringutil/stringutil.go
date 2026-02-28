package stringutil

import "strings"

// MaskSensitiveValue masks a sensitive value for display purposes, showing only the first 4 and last 4 characters.
// For values 8 characters or shorter, returns all asterisks.
// Empty strings are returned as-is.
func MaskSensitiveValue(value string) string {
	if value == "" {
		return ""
	}
	if len(value) <= 8 {
		// For short values, return all asterisks
		return strings.Repeat("*", len(value))
	}
	// Show first 4 and last 4 characters with asterisks in between
	return value[:4] + "..." + value[len(value)-4:]
}
