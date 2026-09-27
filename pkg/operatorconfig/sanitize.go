package operatorconfig

import "strings"

// SanitizeSecret normalizes a secret value for storage and use.
// Secrets are single-line; leading and trailing whitespace and line endings are removed.
// Multiline secret values are not supported.
func SanitizeSecret(s string) string {
	return strings.TrimSpace(s)
}
