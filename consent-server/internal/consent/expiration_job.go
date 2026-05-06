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

	dbmodel "github.com/wso2/openfgc/internal/system/database/model"
	"github.com/wso2/openfgc/internal/system/database/provider"
	"github.com/wso2/openfgc/internal/system/log"
)

var (
	QuerySelectExpiredConsents = dbmodel.DBQuery{
		ID:            "SELECT_EXPIRED_CONSENTS",
		Query:         "SELECT CONSENT_ID FROM CONSENT WHERE VALIDITY_TIME < ? AND CURRENT_STATUS IN (?)",
		PostgresQuery: "SELECT CONSENT_ID FROM CONSENT WHERE VALIDITY_TIME < $1 AND CURRENT_STATUS IN ($2)",
	}

	QueryExpireConsent = dbmodel.DBQuery{
		ID:            "EXPIRE_CONSENT",
		Query:         "UPDATE CONSENT SET CURRENT_STATUS = ?, UPDATED_TIME = ? WHERE CONSENT_ID = ?",
		PostgresQuery: "UPDATE CONSENT SET CURRENT_STATUS = $1, UPDATED_TIME = $2 WHERE CONSENT_ID = $3",
	}

	QueryExpireConsentAuthResources = dbmodel.DBQuery{
		ID:            "EXPIRE_CONSENT_AUTH_RESOURCES",
		Query:         "UPDATE CONSENT_AUTH_RESOURCE SET AUTH_STATUS = ?, UPDATED_TIME = ? WHERE CONSENT_ID = ? AND AUTH_STATUS IN (?)",
		PostgresQuery: "UPDATE CONSENT_AUTH_RESOURCE SET AUTH_STATUS = $1, UPDATED_TIME = $2 WHERE CONSENT_ID = $3 AND AUTH_STATUS IN ($4)",
	}
)

// RunExpirationJob finds consents whose VALIDITY_TIME has passed and marks them as expired.
func RunExpirationJob(activeStatus, expiredStatus, approvedAuthStatus, systemExpiredAuthStatus, systemRevokedAuthStatus string) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "ConsentExpirationJob"))

	dbClient, err := provider.GetDBProvider().GetConsentDBClient()
	if err != nil {
		logger.Error("Failed to get database client", log.Error(err))
		return
	}

	nowMs := time.Now().UnixMilli()

	rows, err := dbClient.Query(QuerySelectExpiredConsents, nowMs, activeStatus)
	if err != nil {
		logger.Error("Failed to query expired consents", log.Error(err))
		return
	}

	if len(rows) == 0 {
		return
	}

	logger.Info("Found consents to expire", log.Int("count", len(rows)))

	for _, row := range rows {
		rawID, ok := row["consent_id"].([]byte)
		if !ok {
			logger.Warn("Failed to parse consent_id from row", log.Any("row", row))
			continue
		}
		consentID := string(rawID)

		_, err := dbClient.Execute(QueryExpireConsent, expiredStatus, nowMs, consentID)
		if err != nil {
			logger.Error("Failed to expire consent", log.Error(err), log.String("consent_id", consentID))
			continue
		}

		_, err = dbClient.Execute(QueryExpireConsentAuthResources,
			systemExpiredAuthStatus, nowMs, consentID,
			approvedAuthStatus, // ← only expire auth resources that are currently approved
		)
		if err != nil {
			logger.Error("Failed to expire consent auth resources",
				log.Error(err),
				log.String("consent_id", consentID),
			)
			continue
		}

		logger.Info("Consent expired successfully", log.String("consent_id", consentID))
	}
}
