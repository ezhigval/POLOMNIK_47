package backup

import (
	"strings"

	"polomnik/internal/adapters/backup/noop"
	"polomnik/internal/adapters/backup/s3"
	"polomnik/internal/config"
	"polomnik/internal/ports"
)

func New(cfg config.Config) ports.BackupStoragePort {
	adapter := strings.ToLower(strings.TrimSpace(cfg.BackupStorageAdapter))
	switch adapter {
	case "s3", "yandex":
		store := s3.New(
			cfg.S3Endpoint,
			cfg.S3Region,
			cfg.S3Bucket,
			cfg.S3AccessKey,
			cfg.S3SecretKey,
			cfg.S3Prefix,
			cfg.IntegrationHTTPTimeout,
		)
		if !store.Configured() {
			return noop.New()
		}
		return store
	default:
		return noop.New()
	}
}
