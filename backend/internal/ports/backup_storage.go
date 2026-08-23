package ports

import (
	"context"
	"errors"
	"io"
)

var ErrBackupStorageNotConfigured = errors.New("backup storage adapter not configured")

type BackupObject struct {
	Name        string
	Body        io.Reader
	Size        int64
	ContentType string
}

// BackupStoragePort uploads nightly dumps offsite when a bucket key exists.
type BackupStoragePort interface {
	Configured() bool
	Upload(ctx context.Context, object BackupObject) error
}
