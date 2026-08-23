package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv                   string
	HTTPAddr                 string
	DatabaseURL              string
	RedisURL                 string
	AdminToken               string
	LogLevel                 string
	CRMAdapter               string
	AccountingAdapter        string
	BitrixWebhookURL         string
	BitrixOriginatorID       string
	BitrixDealCategoryID     int
	BitrixStageNew           string
	BitrixStageContacted     string
	BitrixStageConfirmed     string
	BitrixStageCompleted     string
	BitrixStageCancelled     string
	BitrixUFTourID           string
	BitrixUFPeopleCount      string
	BitrixInboundToken       string
	OneCProtocol             string
	OneCBaseURL              string
	OneCODataPath            string
	OneCODataCounterparty    string
	OneCODataSalesOrder      string
	OneCExportBookingPath    string
	OneCCounterpartyPath     string
	OneCUsername             string
	OneCPassword             string
	OneCAPIKey               string
	IntegrationHTTPTimeout   time.Duration
	CORSAllowOrigins         string
	JWTSecret                string
	JWTTokenTTL              time.Duration
	InternalAPISecret        string
	ShutdownTimeout          time.Duration
	OutboxWorkerPollInterval time.Duration
	OutboxWorkerBatchSize    int
	OutboxWorkerMaxAttempts  int
	NotificationAdapter      string
	TelegramBotToken         string
	TelegramBotUsername      string
	TelegramChatID           string
	TelegramAPIBase          string
	TelegramWebhookURL       string
	TelegramLoginBotToken    string
	TelegramLoginBotUsername string
	MaxBotToken              string
	PhoneAdapter             string
	SMSRUAPIID               string
	MailAdapter              string
	SMTPHost                 string
	SMTPPort                 string
	SMTPUsername             string
	SMTPPassword             string
	SMTPFrom                 string
	SMTPTimeout              time.Duration
	MailForwardToFallback    string
	YandexOAuthClientID      string
	YandexOAuthClientSecret  string
	VKOAuthClientID          string
	VKOAuthClientSecret      string
	MaxOAuthClientID         string
	MaxOAuthClientSecret     string
	MaxOAuthAuthorizeURL     string
	MaxOAuthTokenURL         string
	MaxOAuthUserInfoURL      string
	PublicSiteURL            string
	ManagementBaseURL        string
	UploadDir                string
	UploadPublicBaseURL      string
	UploadMaxBytes           int
	WorkerHeartbeatPath      string
	WorkerHeartbeatMaxAge    time.Duration
	CaptchaAdapter           string
	SmartCaptchaServerKey    string
	SmartCaptchaClientKey    string
	BackupStorageAdapter     string
	S3Endpoint               string
	S3Region                 string
	S3Bucket                 string
	S3AccessKey              string
	S3SecretKey              string
	S3Prefix                 string
	BackupLastPath           string
}

func Load() Config {
	return Config{
		AppEnv:                   envString("APP_ENV", "local"),
		HTTPAddr:                 envString("HTTP_ADDR", ":8080"),
		DatabaseURL:              os.Getenv("DATABASE_URL"),
		RedisURL:                 os.Getenv("REDIS_URL"),
		AdminToken:               os.Getenv("ADMIN_TOKEN"),
		LogLevel:                 envString("LOG_LEVEL", "info"),
		CRMAdapter:               envString("CRM_ADAPTER", "noop"),
		AccountingAdapter:        envString("ACCOUNTING_ADAPTER", "noop"),
		BitrixWebhookURL:         os.Getenv("BITRIX_WEBHOOK_URL"),
		BitrixOriginatorID:       envString("BITRIX_ORIGINATOR_ID", "palomnik"),
		BitrixDealCategoryID:     envInt("BITRIX_DEAL_CATEGORY_ID", 0),
		BitrixStageNew:           envString("BITRIX_STAGE_NEW", ""),
		BitrixStageContacted:     envString("BITRIX_STAGE_CONTACTED", ""),
		BitrixStageConfirmed:     envString("BITRIX_STAGE_CONFIRMED", ""),
		BitrixStageCompleted:     envString("BITRIX_STAGE_COMPLETED", ""),
		BitrixStageCancelled:     envString("BITRIX_STAGE_CANCELLED", ""),
		BitrixUFTourID:           envString("BITRIX_UF_TOUR_ID", ""),
		BitrixUFPeopleCount:      envString("BITRIX_UF_PEOPLE_COUNT", ""),
		BitrixInboundToken:       os.Getenv("BITRIX_INBOUND_TOKEN"),
		OneCProtocol:             envString("ONEC_PROTOCOL", "http"),
		OneCBaseURL:              os.Getenv("ONEC_BASE_URL"),
		OneCODataPath:            envString("ONEC_ODATA_PATH", "/odata/standard.odata"),
		OneCODataCounterparty:    envString("ONEC_ODATA_COUNTERPARTY", "Catalog_Контрагенты"),
		OneCODataSalesOrder:      envString("ONEC_ODATA_SALES_ORDER", "Document_ЗаказКлиента"),
		OneCExportBookingPath:    envString("ONEC_EXPORT_BOOKING_PATH", "/hs/palomnik/booking"),
		OneCCounterpartyPath:     envString("ONEC_COUNTERPARTY_PATH", "/hs/palomnik/counterparty"),
		OneCUsername:             os.Getenv("ONEC_USERNAME"),
		OneCPassword:             os.Getenv("ONEC_PASSWORD"),
		OneCAPIKey:               os.Getenv("ONEC_API_KEY"),
		IntegrationHTTPTimeout:   envDuration("INTEGRATION_HTTP_TIMEOUT", 15*time.Second),
		CORSAllowOrigins:         envString("CORS_ALLOW_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000"),
		JWTSecret:                envString("JWT_SECRET", DefaultJWTSecret),
		JWTTokenTTL:              envDuration("JWT_TOKEN_TTL", 7*24*time.Hour),
		InternalAPISecret:        envString("INTERNAL_API_SECRET", DefaultInternalAPISecret),
		ShutdownTimeout:          envDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
		OutboxWorkerPollInterval: envDuration("OUTBOX_WORKER_POLL_INTERVAL", 30*time.Second),
		OutboxWorkerBatchSize:    envInt("OUTBOX_WORKER_BATCH_SIZE", 10),
		OutboxWorkerMaxAttempts:  envInt("OUTBOX_WORKER_MAX_ATTEMPTS", 5),
		NotificationAdapter:      envString("NOTIFICATION_ADAPTER", "noop"),
		TelegramBotToken:         os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramBotUsername:      os.Getenv("TELEGRAM_BOT_USERNAME"),
		TelegramChatID:           os.Getenv("TELEGRAM_CHAT_ID"),
		TelegramAPIBase:          envString("TELEGRAM_API_BASE", "https://api.telegram.org"),
		TelegramWebhookURL:       os.Getenv("TELEGRAM_WEBHOOK_URL"),
		TelegramLoginBotToken:    os.Getenv("TELEGRAM_LOGIN_BOT_TOKEN"),
		TelegramLoginBotUsername: os.Getenv("TELEGRAM_LOGIN_BOT_USERNAME"),
		MaxBotToken:              os.Getenv("MAX_BOT_TOKEN"),
		PhoneAdapter:             phoneAdapterFromEnv(),
		SMSRUAPIID:               os.Getenv("SMSRU_API_ID"),
		MailAdapter:              envString("MAIL_ADAPTER", "noop"),
		SMTPHost:                 os.Getenv("SMTP_HOST"),
		SMTPPort:                 envString("SMTP_PORT", "587"),
		SMTPUsername:             os.Getenv("SMTP_USERNAME"),
		SMTPPassword:             os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:                 os.Getenv("SMTP_FROM"),
		SMTPTimeout:              envDuration("SMTP_TIMEOUT", 15*time.Second),
		MailForwardToFallback:    os.Getenv("MAIL_FORWARD_TO"),
		YandexOAuthClientID:      os.Getenv("YANDEX_OAUTH_CLIENT_ID"),
		YandexOAuthClientSecret:  os.Getenv("YANDEX_OAUTH_CLIENT_SECRET"),
		VKOAuthClientID:          os.Getenv("VK_OAUTH_CLIENT_ID"),
		VKOAuthClientSecret:      os.Getenv("VK_OAUTH_CLIENT_SECRET"),
		MaxOAuthClientID:         os.Getenv("MAX_OAUTH_CLIENT_ID"),
		MaxOAuthClientSecret:     os.Getenv("MAX_OAUTH_CLIENT_SECRET"),
		MaxOAuthAuthorizeURL:     os.Getenv("MAX_OAUTH_AUTHORIZE_URL"),
		MaxOAuthTokenURL:         os.Getenv("MAX_OAUTH_TOKEN_URL"),
		MaxOAuthUserInfoURL:      os.Getenv("MAX_OAUTH_USERINFO_URL"),
		PublicSiteURL:            envString("PUBLIC_SITE_URL", "https://tikhvin-palomnik.ru"),
		ManagementBaseURL:        envString("MANAGEMENT_BASE_URL", "http://localhost:3000/management/bookings"),
		UploadDir:                envString("UPLOAD_DIR", "./data/uploads"),
		UploadPublicBaseURL:      envString("UPLOAD_PUBLIC_BASE_URL", "http://localhost:8080"),
		UploadMaxBytes:           envInt("UPLOAD_MAX_BYTES", 5*1024*1024),
		WorkerHeartbeatPath:      envString("WORKER_HEARTBEAT_PATH", "/tmp/palomnik-worker-heartbeat"),
		WorkerHeartbeatMaxAge:    envDuration("WORKER_HEARTBEAT_MAX_AGE", 2*time.Minute),
		CaptchaAdapter:           envString("CAPTCHA_ADAPTER", "noop"),
		SmartCaptchaServerKey:    os.Getenv("SMARTCAPTCHA_SERVER_KEY"),
		SmartCaptchaClientKey:    os.Getenv("SMARTCAPTCHA_CLIENT_KEY"),
		BackupStorageAdapter:     envString("BACKUP_STORAGE_ADAPTER", "noop"),
		S3Endpoint:               envString("S3_ENDPOINT", "https://storage.yandexcloud.net"),
		S3Region:                 envString("S3_REGION", "ru-central1"),
		S3Bucket:                 os.Getenv("S3_BUCKET"),
		S3AccessKey:              os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey:              os.Getenv("S3_SECRET_KEY"),
		S3Prefix:                 envString("S3_PREFIX", "polomnik-backups"),
		BackupLastPath:           envString("BACKUP_LAST_PATH", "./backups/last-backup.json"),
	}
}

func (c Config) EffectiveTelegramLoginBotToken() string {
	if value := strings.TrimSpace(c.TelegramLoginBotToken); value != "" {
		return value
	}
	return strings.TrimSpace(c.TelegramBotToken)
}

func (c Config) EffectiveTelegramLoginBotUsername() string {
	if value := strings.TrimSpace(c.TelegramLoginBotUsername); value != "" {
		return strings.TrimPrefix(value, "@")
	}
	return strings.TrimPrefix(strings.TrimSpace(c.TelegramBotUsername), "@")
}

func (c Config) EffectiveTelegramWebhookURL() string {
	if value := strings.TrimSpace(c.TelegramWebhookURL); value != "" {
		return value
	}
	base := strings.TrimRight(strings.TrimSpace(c.UploadPublicBaseURL), "/")
	if base == "" {
		return ""
	}
	return base + "/api/v1/webhooks/telegram"
}

func (c Config) TelegramAPIHostIsOfficial() bool {
	base := strings.ToLower(strings.TrimSpace(c.TelegramAPIBase))
	return base == "" || strings.Contains(base, "api.telegram.org")
}

// phoneAdapterFromEnv prefers PHONE_ADAPTER; SMS_ADAPTER is a legacy alias.
func phoneAdapterFromEnv() string {
	if value := strings.TrimSpace(os.Getenv("PHONE_ADAPTER")); value != "" {
		return value
	}
	if value := strings.TrimSpace(os.Getenv("SMS_ADAPTER")); value != "" {
		return value
	}
	return "noop"
}

func envString(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return duration
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
