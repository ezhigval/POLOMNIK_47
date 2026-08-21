package dto

type OutboxSummary struct {
	Pending           int     `json:"pending"`
	Failed            int     `json:"failed"`
	Processed         int     `json:"processed"`
	OldestPendingAt   *string `json:"oldest_pending_at,omitempty"`
	LatestFailedAt    *string `json:"latest_failed_at,omitempty"`
	LatestFailedError string  `json:"latest_failed_error,omitempty"`
}

type SystemInfo struct {
	CRMAdapter          string        `json:"crm_adapter"`
	AccountingAdapter   string        `json:"accounting_adapter"`
	NotificationAdapter string        `json:"notification_adapter"`
	TelegramConfigured  bool          `json:"telegram_configured"`
	BitrixConfigured    bool          `json:"bitrix_configured"`
	OneCConfigured      bool          `json:"onec_configured"`
	Outbox              OutboxSummary `json:"outbox"`
}
