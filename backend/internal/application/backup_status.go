package application

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type BackupStatus struct {
	At           *time.Time
	File         string
	Bytes        int64
	Offsite      bool
	OffsiteError string
}

type backupStatusFile struct {
	At           string `json:"at"`
	File         string `json:"file"`
	Bytes        int64  `json:"bytes"`
	Offsite      bool   `json:"offsite"`
	OffsiteError string `json:"offsite_error"`
}

func ReadBackupStatus(path string) BackupStatus {
	path = strings.TrimSpace(path)
	if path == "" {
		return BackupStatus{}
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return BackupStatus{}
	}
	var payload backupStatusFile
	if json.Unmarshal(raw, &payload) != nil {
		return BackupStatus{}
	}
	status := BackupStatus{
		File:         payload.File,
		Bytes:        payload.Bytes,
		Offsite:      payload.Offsite,
		OffsiteError: payload.OffsiteError,
	}
	if parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(payload.At)); err == nil {
		utc := parsed.UTC()
		status.At = &utc
	}
	return status
}

func WriteBackupStatus(path string, status BackupStatus) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	payload := backupStatusFile{
		File:         status.File,
		Bytes:        status.Bytes,
		Offsite:      status.Offsite,
		OffsiteError: status.OffsiteError,
	}
	if status.At != nil {
		payload.At = status.At.UTC().Format(time.RFC3339)
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
