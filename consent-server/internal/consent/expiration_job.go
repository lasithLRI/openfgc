/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package consent

import (
	"time"

	"github.com/wso2/openfgc/internal/system/log"
)

// RunExpirationJob finds all consents whose VALIDITY_TIME has passed and marks them
// as expired, along with all their auth resources.
// It is safe to run concurrently — each tick launches a new goroutine, and any
// panic is recovered so the scheduler goroutine is never killed.
func RunExpirationJob(statuses ExpirationStatuses) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "ConsentExpirationJob"))

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
		err := s.ExpireConsent(
			nowMs,
			consentID,
			statuses.ExpiredConsentStatus,
			statuses.SystemExpiredAuthStatus,
		)
		if err != nil {
			logger.Error("Failed to expire consent",
				log.Error(err),
				log.String("consent_id", consentID),
			)
			continue
		}
		logger.Info("Consent expired successfully", log.String("consent_id", consentID))
	}
}
