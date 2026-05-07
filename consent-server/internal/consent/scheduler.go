package consent

import (
	"time"

	"github.com/wso2/openfgc/internal/system/log"
)

// ExpirationStatuses groups all status strings needed by the expiration job.
type ExpirationStatuses struct {
	ExpirableConsentStatuses []string // e.g. ["ACTIVE", "CREATED"]
	ExpiredConsentStatus     string   // e.g. "EXPIRED"
	ExpirableAuthStatuses    []string // e.g. ["APPROVED", "CREATED"]
	SystemExpiredAuthStatus  string   // e.g. "SYS_EXPIRED"
}

// StartScheduler starts the consent expiration scheduler at the given interval.
func StartScheduler(interval time.Duration, statuses ExpirationStatuses) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "ConsentScheduler"))

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.Info("Consent expiration scheduler started",
		log.String("interval", interval.String()),
		log.Any("expirable_consent_statuses", statuses.ExpirableConsentStatuses),
		log.Any("expirable_auth_statuses", statuses.ExpirableAuthStatuses),
		log.String("expired_consent_status", statuses.ExpiredConsentStatus),
		log.String("system_expired_auth_status", statuses.SystemExpiredAuthStatus),
	)

	for range ticker.C {
		logger.Debug("Scheduler tick — launching expiration job")
		go RunExpirationJob(statuses)
	}

	// This should never be reached — log if it is
	logger.Error("Consent expiration scheduler exited unexpectedly")
}
