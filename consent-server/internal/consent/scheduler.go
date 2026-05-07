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

// ExpirationStatuses groups all status strings needed by the expiration job.
type ExpirationStatuses struct {
	ExpirableConsentStatuses []string // e.g. ["ACTIVE", "CREATED"]
	ExpiredConsentStatus     string   // e.g. "EXPIRED"
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