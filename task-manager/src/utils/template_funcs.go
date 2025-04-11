package utils

import (
	"html/template"
	"time"
)

// GetTemplateFunctions returns a map of custom functions for use in templates
func GetTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		// Truncate a string to the specified length and add "..." if truncated
		"truncate": func(s string, length int) string {
			if len(s) <= length {
				return s
			}
			return s[0:length] + "..."
		},
		// Format a date to a human-readable string
		"formatDate": func(t time.Time) string {
			return t.Format("Jan 02, 2006")
		},
	}
}
