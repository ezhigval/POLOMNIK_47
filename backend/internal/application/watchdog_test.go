package application

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestWatchdogBackupOverdueWhenMissing(t *testing.T) {
	fixed := time.Date(2026, 8, 23, 20, 0, 0, 0, time.UTC)
	w := &Watchdog{Now: func() time.Time { return fixed }}
	report := w.Snapshot(context.Background())
	if !report.BackupOverdue {
		t.Fatal("missing backup must be overdue")
	}
	if report.RestartAttempted {
		t.Fatal("must not restart")
	}
}

func TestWatchdogBackupFreshIsNotOverdue(t *testing.T) {
	fixed := time.Date(2026, 8, 23, 20, 0, 0, 0, time.UTC)
	at := fixed.Add(-2 * time.Hour)
	path := filepath.Join(t.TempDir(), "last-backup.json")
	if err := WriteBackupStatus(path, BackupStatus{At: &at, File: "dump.sql.gz"}); err != nil {
		t.Fatalf("write: %v", err)
	}
	w := &Watchdog{
		Now:        func() time.Time { return fixed },
		BackupPath: path,
	}
	report := w.Snapshot(context.Background())
	if report.BackupOverdue {
		t.Fatalf("fresh backup marked overdue: %+v", report)
	}
}

func TestWatchdogBackupOlderThan26hIsOverdue(t *testing.T) {
	fixed := time.Date(2026, 8, 23, 20, 0, 0, 0, time.UTC)
	at := fixed.Add(-27 * time.Hour)
	path := filepath.Join(t.TempDir(), "last-backup.json")
	if err := WriteBackupStatus(path, BackupStatus{At: &at, File: "dump.sql.gz"}); err != nil {
		t.Fatalf("write: %v", err)
	}
	w := &Watchdog{
		Now:        func() time.Time { return fixed },
		BackupPath: path,
	}
	report := w.Snapshot(context.Background())
	if !report.BackupOverdue {
		t.Fatal("expected overdue after 26h")
	}
}
