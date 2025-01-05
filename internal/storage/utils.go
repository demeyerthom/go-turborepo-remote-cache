package storage

import (
	"regexp"
)

func sanitizeSlug(slug string) string {
	// Replace any non-alphanumeric characters with an underscore
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return re.ReplaceAllString(slug, "_")
}
