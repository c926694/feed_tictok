package util

import (
	"simple_tiktok/internal/pkg/constants"
	"strings"
)

// EnsureHTTPPath prepends HttpPath for relative resource paths.
func EnsureHTTPPath(path string) string {
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	base := strings.TrimRight(constants.HttpPath, "/")
	if strings.HasPrefix(path, "/") {
		return base + path
	}
	return base + "/" + path
}

