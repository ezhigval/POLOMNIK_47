package application

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

// ManagementOrigin extracts scheme://host from MANAGEMENT_BASE_URL
// (often ends with /management/bookings) without relying on TrimSuffix alone.
func ManagementOrigin(mgmtBase string) string {
	raw := strings.TrimSpace(mgmtBase)
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		trimmed := strings.TrimRight(strings.TrimSpace(mgmtBase), "/")
		trimmed = strings.TrimSuffix(trimmed, "/bookings")
		trimmed = strings.TrimSuffix(trimmed, "/management")
		return strings.TrimRight(trimmed, "/")
	}
	return strings.TrimRight(u.Scheme+"://"+u.Host, "/")
}

func ManagementSupportThreadURL(mgmtBase string, threadID uuid.UUID) string {
	origin := ManagementOrigin(mgmtBase)
	if origin == "" || threadID == uuid.Nil {
		return ""
	}
	return fmt.Sprintf("%s/management/support/%s", origin, threadID.String())
}
