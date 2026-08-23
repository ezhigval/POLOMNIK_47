package dto

import (
	"time"

	"palomnik/internal/application"
	"palomnik/internal/domain"
)

type SupportDraftResponse struct {
	Configured bool   `json:"configured"`
	Escalate   bool   `json:"escalate"`
	Draft      string `json:"draft"`
	Note       string `json:"note"`
}

func ToSupportDraftResponse(draft application.SupportDraft) SupportDraftResponse {
	return SupportDraftResponse{
		Configured: draft.Configured,
		Escalate:   draft.Escalate,
		Draft:      draft.Draft,
		Note:       draft.Note,
	}
}

type MetricsDigestResponse struct {
	Configured         bool           `json:"configured"`
	BookingsByStatus   map[string]int `json:"bookings_by_status"`
	ActiveTours        int            `json:"active_tours"`
	OpenSupportThreads int            `json:"open_support_threads"`
	OutboxPending      int            `json:"outbox_pending"`
	OutboxFailed       int            `json:"outbox_failed"`
	Summary            string         `json:"summary,omitempty"`
}

func ToMetricsDigestResponse(digest application.MetricsDigest) MetricsDigestResponse {
	byStatus := digest.BookingsByStatus
	if byStatus == nil {
		byStatus = map[string]int{}
	}
	return MetricsDigestResponse{
		Configured:         digest.Configured,
		BookingsByStatus:   byStatus,
		ActiveTours:        digest.ActiveTours,
		OpenSupportThreads: digest.OpenSupportThreads,
		OutboxPending:      digest.OutboxPending,
		OutboxFailed:       digest.OutboxFailed,
		Summary:            digest.Summary,
	}
}

type WatchdogResponse struct {
	Configured       bool    `json:"configured"`
	At               string  `json:"at"`
	Database         string  `json:"database"`
	DiskPath         string  `json:"disk_path,omitempty"`
	DiskUsedBytes    uint64  `json:"disk_used_bytes"`
	DiskTotalBytes   uint64  `json:"disk_total_bytes"`
	DiskPercent      int     `json:"disk_percent"`
	OutboxPending    int     `json:"outbox_pending"`
	OutboxFailed     int     `json:"outbox_failed"`
	Status5xx        uint64  `json:"status_5xx"`
	BackupAt         *string `json:"backup_at,omitempty"`
	BackupOverdue    bool    `json:"backup_overdue"`
	RestartAttempted bool    `json:"restart_attempted"`
	Summary          string  `json:"summary,omitempty"`
}

func ToWatchdogResponse(report application.AIWatchdogReport) WatchdogResponse {
	out := WatchdogResponse{
		Configured:       report.Configured,
		At:               report.At.UTC().Format(time.RFC3339),
		Database:         report.Database,
		DiskPath:         report.DiskPath,
		DiskUsedBytes:    report.DiskUsedBytes,
		DiskTotalBytes:   report.DiskTotalBytes,
		DiskPercent:      report.DiskPercent,
		OutboxPending:    report.OutboxPending,
		OutboxFailed:     report.OutboxFailed,
		Status5xx:        report.Status5xx,
		BackupOverdue:    report.BackupOverdue,
		RestartAttempted: false,
		Summary:          report.Summary,
	}
	if report.BackupAt != nil {
		value := report.BackupAt.UTC().Format(time.RFC3339)
		out.BackupAt = &value
	}
	return out
}

func ToTourList(tours []domain.Tour) []TourResponse {
	items := make([]TourResponse, 0, len(tours))
	for _, tour := range tours {
		items = append(items, ToTourResponse(tour))
	}
	return items
}
