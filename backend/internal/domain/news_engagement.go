package domain

import (
	"strings"
	"unicode/utf8"
)

const (
	MaxNewsCommentRunes = 2000
	MaxVisitorIDLen     = 64
)

func NormalizeVisitorID(raw string) (string, error) {
	id := strings.TrimSpace(raw)
	if id == "" || utf8.RuneCountInString(id) > MaxVisitorIDLen {
		return "", ErrInvalidVisitorID
	}
	return id, nil
}

func ValidateNewsCommentBody(body string) (string, error) {
	text := strings.TrimSpace(body)
	if text == "" {
		return "", ErrInvalidCommentBody
	}
	if utf8.RuneCountInString(text) > MaxNewsCommentRunes {
		return "", ErrInvalidCommentBody
	}
	return text, nil
}
