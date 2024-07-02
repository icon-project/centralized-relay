package connutil

import "strings"

func ShouldReconnect(err error) bool {
	if strings.Contains(err.Error(), "operation timed out") ||
		strings.Contains(err.Error(), "HTTP 502") ||
		strings.Contains(err.Error(), "connection refused") {
		return true
	}
	return false
}
