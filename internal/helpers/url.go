package helpers

import "strings"

// ValidateURL checks if the provided URL starts with "http://" or "https://". Returns a message and validity as a boolean.
func ValidateURL(url string) (string, bool) {
	if !strings.HasPrefix(url, "http://") &&
		!strings.HasPrefix(url, "https://") {
		return "Invalid URL", true
	}
	return "", false
}
