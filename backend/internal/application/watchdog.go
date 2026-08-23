package application

import (
	"context"
	"log/slog"
	"syscall"
	"time"

	"palomnik/internal/ports"
)

// Nightly dump is 03:00. Missing more than 26 hours means the last run did not land.
const BackupOverdueAfter = 26 * time.Hour

type WatchdogReport struct {
	At               time.Time
	Database         string
	DiskPath         string
	DiskUsedBytes    uint64
	DiskTotalBytes   uint64
	DiskPercent      int
	OutboxPending    int
	OutboxFailed     int
	Status5xx        uint64
	BackupAt         *time.Time
	BackupOverdue    bool
	RestartAttempted bool
}

type Watchdog struct {
	Outbox     ports.OutboxRepository
	BackupPath string
	DiskPath   string
	PingDB     func(context.Context) error
	Status5xx  func() uint64
	Now        func() time.Time
}

func (w *Watchdog) Snapshot(ctx context.Context) WatchdogReport {
	now := time.Now().UTC()
	if w != nil && w.Now != nil {
		now = w.Now().UTC()
	}
	report := WatchdogReport{
		At:               now,
		Database:         "unknown",
		RestartAttempted: false,
	}
	if w == nil {
		return report
	}
	report.DiskPath = w.DiskPath
	if w.PingDB != nil {
		if err := w.PingDB(ctx); err != nil {
			report.Database = "error"
		} else {
			report.Database = "ok"
		}
	}
	if used, total, ok := diskUsage(w.DiskPath); ok {
		report.DiskUsedBytes = used
		report.DiskTotalBytes = total
		if total > 0 {
			report.DiskPercent = int((used * 100) / total)
		}
	}
	if w.Outbox != nil {
		if summary, err := w.Outbox.Summarize(ctx); err == nil {
			report.OutboxPending = summary.Pending
			report.OutboxFailed = summary.Failed
		}
	}
	if w.Status5xx != nil {
		report.Status5xx = w.Status5xx()
	}
	backup := ReadBackupStatus(w.BackupPath)
	report.BackupAt = backup.At
	if backup.At == nil || now.Sub(backup.At.UTC()) > BackupOverdueAfter {
		report.BackupOverdue = true
	}
	return report
}

func (w *Watchdog) Log(ctx context.Context, log *slog.Logger) WatchdogReport {
	report := w.Snapshot(ctx)
	if log != nil {
		log.Info(
			"watchdog report",
			slog.Time("at", report.At),
			slog.String("database", report.Database),
			slog.Int("disk_percent", report.DiskPercent),
			slog.Uint64("disk_used_bytes", report.DiskUsedBytes),
			slog.Uint64("disk_total_bytes", report.DiskTotalBytes),
			slog.Int("outbox_pending", report.OutboxPending),
			slog.Int("outbox_failed", report.OutboxFailed),
			slog.Uint64("status_5xx", report.Status5xx),
			slog.Bool("backup_overdue", report.BackupOverdue),
			slog.Bool("restart_attempted", false),
		)
	}
	return report
}

func diskUsage(path string) (used, total uint64, ok bool) {
	if path == "" {
		path = "/"
	}
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil {
		return 0, 0, false
	}
	bsize := uint64(st.Bsize)
	total = st.Blocks * bsize
	free := st.Bavail * bsize
	if total < free {
		return 0, total, true
	}
	return total - free, total, true
}
