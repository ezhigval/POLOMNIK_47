package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv            string
	HTTPAddr          string
	DatabaseURL       string
	RedisURL          string
	AdminToken        string
	LogLevel          string
	CRMAdapter        string
	AccountingAdapter string
	BitrixWebhookURL  string
	BitrixOriginatorID string
	BitrixDealCategoryID int
	BitrixStageNew string
	BitrixStageContacted string
	BitrixStageConfirmed string
	BitrixStageCompleted string
	BitrixStageCancelled string
	BitrixUFTourID string
	BitrixUFPeopleCount string
	BitrixInboundToken string
	OneCProtocol string
	OneCBaseURL string
	OneCODataPath string
	OneCODataCounterparty string
	OneCODataSalesOrder string
	OneCExportBookingPath string
	OneCCounterpartyPath string
	OneCUsername string
	OneCPassword string
	OneCAPIKey string
	IntegrationHTTPTimeout time.Duration
	CORSAllowOrigins  string
	JWTSecret         string
	JWTTokenTTL       time.Duration
	InternalAPISecret string
	ShutdownTimeout   time.Duration
	OutboxWorkerPollInterval time.Duration
	OutboxWorkerBatchSize    int
	OutboxWorkerMaxAttempts  int
	NotificationAdapter      string
	TelegramBotToken         string
	TelegramChatID           string
	TelegramAPIBase          string
	ManagementBaseURL        string
	UploadDir                string
	UploadPublicBaseURL      string
	UploadMaxBytes           int
	WorkerHeartbeatPath      string
	WorkerHeartbeatMaxAge    time.Duration
}

func Load() Config {
	return Config{
		AppEnv:            envString("APP_ENV", "local"),
		HTTPAddr:          envString("HTTP_ADDR", ":8080"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		RedisURL:          os.Getenv("REDIS_URL"),
		AdminToken:        os.Getenv("ADMIN_TOKEN"),
		LogLevel:          envString("LOG_LEVEL", "info"),
		CRMAdapter:        envString("CRM_ADAPTER", "noop"),
		AccountingAdapter: envString("ACCOUNTING_ADAPTER", "noop"),
		BitrixWebhookURL:  os.Getenv("BITRIX_WEBHOOK_URL"),
		BitrixOriginatorID: envString("BITRIX_ORIGINATOR_ID", "polomnik"),
		BitrixDealCategoryID: envInt("BITRIX_DEAL_CATEGORY_ID", 0),
		BitrixStageNew: envString("BITRIX_STAGE_NEW", ""),
		BitrixStageContacted: envString("BITRIX_STAGE_CONTACTED", ""),
		BitrixStageConfirmed: envString("BITRIX_STAGE_CONFIRMED", ""),
		BitrixStageCompleted: envString("BITRIX_STAGE_COMPLETED", ""),
		BitrixStageCancelled: envString("BITRIX_STAGE_CANCELLED", ""),
		BitrixUFTourID: envString("BITRIX_UF_TOUR_ID", ""),
		BitrixUFPeopleCount: envString("BITRIX_UF_PEOPLE_COUNT", ""),
		BitrixInboundToken: os.Getenv("BITRIX_INBOUND_TOKEN"),
		OneCProtocol: envString("ONEC_PROTOCOL", "http"),
		OneCBaseURL:       os.Getenv("ONEC_BASE_URL"),
		OneCODataPath: envString("ONEC_ODATA_PATH", "/odata/standard.odata"),
		OneCODataCounterparty: envString("ONEC_ODATA_COUNTERPARTY", "Catalog_Контрагенты"),
		OneCODataSalesOrder: envString("ONEC_ODATA_SALES_ORDER", "Document_ЗаказКлиента"),
		OneCExportBookingPath: envString("ONEC_EXPORT_BOOKING_PATH", "/hs/polomnik/booking"),
		OneCCounterpartyPath: envString("ONEC_COUNTERPARTY_PATH", "/hs/polomnik/counterparty"),
		OneCUsername: os.Getenv("ONEC_USERNAME"),
		OneCPassword: os.Getenv("ONEC_PASSWORD"),
		OneCAPIKey: os.Getenv("ONEC_API_KEY"),
		IntegrationHTTPTimeout: envDuration("INTEGRATION_HTTP_TIMEOUT", 15*time.Second),
		CORSAllowOrigins:  envString("CORS_ALLOW_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000"),
		JWTSecret:         envString("JWT_SECRET", DefaultJWTSecret),
		JWTTokenTTL:       envDuration("JWT_TOKEN_TTL", 7*24*time.Hour),
		InternalAPISecret: envString("INTERNAL_API_SECRET", DefaultInternalAPISecret),
		ShutdownTimeout:   envDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
		OutboxWorkerPollInterval: envDuration("OUTBOX_WORKER_POLL_INTERVAL", 30*time.Second),
		OutboxWorkerBatchSize:    envInt("OUTBOX_WORKER_BATCH_SIZE", 10),
		OutboxWorkerMaxAttempts:  envInt("OUTBOX_WORKER_MAX_ATTEMPTS", 5),
		NotificationAdapter:      envString("NOTIFICATION_ADAPTER", "noop"),
		TelegramBotToken:         os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:           os.Getenv("TELEGRAM_CHAT_ID"),
		TelegramAPIBase:          envString("TELEGRAM_API_BASE", "https://api.telegram.org"),
		ManagementBaseURL:        envString("MANAGEMENT_BASE_URL", "http://localhost:3000/management/bookings"),
		UploadDir:                envString("UPLOAD_DIR", "./data/uploads"),
		UploadPublicBaseURL:      envString("UPLOAD_PUBLIC_BASE_URL", "http://localhost:8080"),
		UploadMaxBytes:           envInt("UPLOAD_MAX_BYTES", 5*1024*1024),
		WorkerHeartbeatPath:      envString("WORKER_HEARTBEAT_PATH", "/tmp/polomnik-worker-heartbeat"),
		WorkerHeartbeatMaxAge:    envDuration("WORKER_HEARTBEAT_MAX_AGE", 2*time.Minute),
	}
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
