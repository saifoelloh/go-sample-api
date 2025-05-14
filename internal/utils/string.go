package utils

import (
	"strings"
)

// GetUniqueLabel returns the first non-empty value between email and ssoID
func GetUniqueLabel(email, ssoID *string) string {
	if email != nil && strings.TrimSpace(*email) != "" {
		return *email
	}
	if ssoID != nil && strings.TrimSpace(*ssoID) != "" {
		return *ssoID
	}
	return "unknown"
}
