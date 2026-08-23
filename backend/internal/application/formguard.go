package application

import "strings"

func HoneypotTriggered(value string) bool {
	return strings.TrimSpace(value) != ""
}
