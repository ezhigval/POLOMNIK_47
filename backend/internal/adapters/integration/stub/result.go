package stub

import "palomnik/internal/ports"

func NotConfigured(system string) ports.IntegrationResult {
	return ports.IntegrationResult{
		Status:  ports.IntegrationStatusNotConfigured,
		Message: system + " adapter is not configured",
	}
}

func Pending(system, message string) ports.IntegrationResult {
	return ports.IntegrationResult{
		Status:  ports.IntegrationStatusPending,
		Message: system + " adapter stub: " + message,
	}
}
