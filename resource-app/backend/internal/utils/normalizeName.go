package utils

import (
	"strings"
)

// NormalizeName converts input like:
// " admin   group " -> "admin group"
// "pHySiCs   lab"   -> "physics lab"
func NormalizeName(input string) string {
	// Trim + collapse multiple spaces
	cleaned := strings.Join(strings.Fields(strings.TrimSpace(input)), " ")
	if cleaned == "" {
		return ""
	}
	return strings.ToLower(cleaned)
}