package consent

import (
	"time"

	"github.com/wso2/openfgc/internal/system/log"
)

// RunExpirationJob finds consents whose VALIDITY_TIME has passed and marks them as expired.
func RunExpirationJob(statuses ExpirationStatuses) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "ConsentExpirationJob"))

	// Catch any silent panics so the scheduler goroutine is never killed
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic recovered in expiration job", log.Any("panic", r))
		}
	}()

	logger.Debug("Running consent expiration job")

	s := &store{}
	nowMs := time.Now().UnixMilli()

	consentIDs, err := s.GetExpiredConsentIDs(nowMs, statuses.ExpirableConsentStatuses)
	if err != nil {
		logger.Error("Failed to query expired consents", log.Error(err))
		return
	}

	if len(consentIDs) == 0 {
		logger.Debug("No consents to expire")
		return
	}

	logger.Info("Found consents to expire", log.Int("count", len(consentIDs)))

	for _, consentID := range consentIDs {
		err := s.ExpireConsent(nowMs, consentID, statuses.ExpiredConsentStatus, statuses.SystemExpiredAuthStatus, statuses.ExpirableAuthStatuses)
		if err != nil {
			logger.Error("Failed to expire consent", log.Error(err), log.String("consent_id", consentID))
			continue
		}
		logger.Info("Consent expired successfully", log.String("consent_id", consentID))
	}
}
