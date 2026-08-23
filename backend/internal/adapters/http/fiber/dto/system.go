package dto

type OutboxSummary struct {
	Pending           int     `json:"pending"`
	Failed            int     `json:"failed"`
	Processed         int     `json:"processed"`
	OldestPendingAt   *string `json:"oldest_pending_at,omitempty"`
	LatestFailedAt    *string `json:"latest_failed_at,omitempty"`
	LatestFailedError string  `json:"latest_failed_error,omitempty"`
}

type LatencyInfo struct {
	LastMS   int64  `json:"last_ms"`
	AvgMS    int64  `json:"avg_ms"`
	Requests uint64 `json:"requests"`
}

type BackupInfo struct {
	At           *string `json:"at,omitempty"`
	File         string  `json:"file,omitempty"`
	Bytes        int64   `json:"bytes,omitempty"`
	Offsite      bool    `json:"offsite"`
	OffsiteError string  `json:"offsite_error,omitempty"`
}

type SystemInfo struct {
	CRMAdapter          string        `json:"crm_adapter"`
	AccountingAdapter   string        `json:"accounting_adapter"`
	NotificationAdapter string        `json:"notification_adapter"`
	MessengerAdapter    string        `json:"messenger_adapter"`
	PublisherAdapter    string        `json:"publisher_adapter"`
	AIAdapter           string        `json:"ai_adapter"`
	TelegramConfigured  bool          `json:"telegram_configured"`
	MessengerConfigured bool          `json:"messenger_configured"`
	PublisherConfigured bool          `json:"publisher_configured"`
	AIConfigured        bool          `json:"ai_configured"`
	BitrixConfigured    bool          `json:"bitrix_configured"`
	OneCConfigured      bool          `json:"onec_configured"`
	Outbox              OutboxSummary `json:"outbox"`
	Latency             LatencyInfo   `json:"latency"`
	LastBackup          BackupInfo    `json:"last_backup"`
}
