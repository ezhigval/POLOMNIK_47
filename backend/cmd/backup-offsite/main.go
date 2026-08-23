package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"palomnik/internal/adapters/backup"
	"palomnik/internal/application"
	"palomnik/internal/config"
	"palomnik/internal/logger"
	"palomnik/internal/ports"
)

func main() {
	os.Exit(run())
}

func run() int {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)
	if len(os.Args) < 2 {
		log.Error("usage: backup-offsite <dump.sql.gz>")
		return 2
	}

	filePath := os.Args[1]
	info, err := os.Stat(filePath)
	if err != nil {
		log.Error("backup file missing", slog.Any("error", err))
		return 1
	}

	status := application.BackupStatus{
		File:  filepath.Base(filePath),
		Bytes: info.Size(),
	}
	now := time.Now().UTC()
	status.At = &now

	store := backup.New(cfg)
	if store.Configured() {
		file, err := os.Open(filePath)
		if err != nil {
			log.Error("open backup", slog.Any("error", err))
			return 1
		}
		uploadErr := store.Upload(context.Background(), ports.BackupObject{
			Name:        filepath.Base(filePath),
			Body:        file,
			Size:        info.Size(),
			ContentType: "application/gzip",
		})
		_ = file.Close()
		if uploadErr != nil {
			status.OffsiteError = uploadErr.Error()
			log.Error("offsite upload failed", slog.Any("error", uploadErr))
		} else {
			status.Offsite = true
			log.Info("offsite backup uploaded", slog.String("file", status.File))
		}
	}

	if err := application.WriteBackupStatus(cfg.BackupLastPath, status); err != nil {
		log.Error("write backup status", slog.Any("error", err))
		return 1
	}
	if status.OffsiteError != "" {
		return 1
	}
	return 0
}
