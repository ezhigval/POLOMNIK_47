package postformat

import (
	"strings"

	"palomnik/internal/ports"
)

// Post joins title, body and URL. No extra marketing copy.
func Post(content ports.PublishContent) string {
	var parts []string
	if title := strings.TrimSpace(content.Title); title != "" {
		parts = append(parts, title)
	}
	if body := strings.TrimSpace(content.Body); body != "" {
		parts = append(parts, body)
	}
	if link := strings.TrimSpace(content.URL); link != "" {
		parts = append(parts, link)
	}
	return strings.Join(parts, "\n\n")
}

// EscapeHTML keeps Telegram parse_mode=HTML from treating post text as markup.
func EscapeHTML(value string) string {
	value = strings.ReplaceAll(value, "&", "&amp;")
	value = strings.ReplaceAll(value, "<", "&lt;")
	value = strings.ReplaceAll(value, ">", "&gt;")
	return value
}
