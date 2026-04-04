// Package modifier
package modifier

import (
	"net/http"
	"path/filepath"
	"strings"
)

func DetectContentType(filename string, data []byte) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".css":
		return "text/css"
	case ".html":
		return "text/html"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".svg":
		return "image/svg+xml"
	case ".txt":
		return "text/plain"
	}

	// fallback
	return http.DetectContentType(data)
}
