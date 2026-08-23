package application

import (
	"path/filepath"
	"testing"
	"time"
)

func TestWriteAndReadBackupStatus(t *testing.T) {
	path := filepath.Join(t.TempDir(), "last-backup.json")
	at := time.Date(2026, 8, 23, 3, 0, 0, 0, time.UTC)
	if err := WriteBackupStatus(path, BackupStatus{
		At:      &at,
		File:    "palomnik-20260823.sql.gz",
		Bytes:   12,
		Offsite: true,
	}); err != nil {
		t.Fatalf("write: %v", err)
	}

	got := ReadBackupStatus(path)
	if got.File != "palomnik-20260823.sql.gz" || !got.Offsite || got.Bytes != 12 {
		t.Fatalf("unexpected status %+v", got)
	}
	if got.At == nil || !got.At.Equal(at) {
		t.Fatalf("unexpected time %+v", got.At)
	}
}
