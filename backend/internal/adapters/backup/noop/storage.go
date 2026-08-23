package noop

import (
	"context"

	"palomnik/internal/ports"
)

type Storage struct{}

func New() Storage { return Storage{} }

var _ ports.BackupStoragePort = Storage{}

func (Storage) Configured() bool { return false }

func (Storage) Upload(context.Context, ports.BackupObject) error {
	return ports.ErrBackupStorageNotConfigured
}
