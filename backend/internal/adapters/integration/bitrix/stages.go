package bitrix

import (
	"strings"

	"polomnik/internal/config"
	"polomnik/internal/domain"
)

var defaultForwardStages = map[domain.BookingStatus]string{
	domain.BookingStatusNew:       "NEW",
	domain.BookingStatusContacted: "PREPARATION",
	domain.BookingStatusConfirmed: "EXECUTING",
	domain.BookingStatusCompleted: "WON",
	domain.BookingStatusCancelled: "LOSE",
}

func bookingStageID(cfg config.Config, status domain.BookingStatus) string {
	if stage := stageFromConfig(cfg, status); stage != "" {
		return stage
	}
	return defaultForwardStages[status]
}

func stageFromConfig(cfg config.Config, status domain.BookingStatus) string {
	switch status {
	case domain.BookingStatusNew:
		return strings.TrimSpace(cfg.BitrixStageNew)
	case domain.BookingStatusContacted:
		return strings.TrimSpace(cfg.BitrixStageContacted)
	case domain.BookingStatusConfirmed:
		return strings.TrimSpace(cfg.BitrixStageConfirmed)
	case domain.BookingStatusCompleted:
		return strings.TrimSpace(cfg.BitrixStageCompleted)
	case domain.BookingStatusCancelled:
		return strings.TrimSpace(cfg.BitrixStageCancelled)
	default:
		return ""
	}
}

func bookingStatusFromStage(cfg config.Config, stageID string) (domain.BookingStatus, bool) {
	stageID = strings.TrimSpace(stageID)
	if stageID == "" {
		return "", false
	}

	for status, mapped := range map[domain.BookingStatus]string{
		domain.BookingStatusNew:       stageFromConfig(cfg, domain.BookingStatusNew),
		domain.BookingStatusContacted: stageFromConfig(cfg, domain.BookingStatusContacted),
		domain.BookingStatusConfirmed: stageFromConfig(cfg, domain.BookingStatusConfirmed),
		domain.BookingStatusCompleted: stageFromConfig(cfg, domain.BookingStatusCompleted),
		domain.BookingStatusCancelled: stageFromConfig(cfg, domain.BookingStatusCancelled),
	} {
		if mapped != "" && mapped == stageID {
			return status, true
		}
	}

	for status, mapped := range defaultForwardStages {
		if mapped == stageID {
			return status, true
		}
	}

	return "", false
}
