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
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package config

// ConsentStatus represents a typed consent status.
type ConsentStatus string

// AuthStatus represents a typed authorization status.
type AuthStatus string

// ConsentConfig holds consent-related configuration.
type ConsentConfig struct {
	ExpirationFrequency ExpirationFrequencyConfig `yaml:"expiration_frequency"`
	EligibleStatuses    EligibleStatusesConfig    `yaml:"eligible_statuses"`
	StatusMappings      ConsentStatusMappings     `yaml:"status_mappings"`
	AuthStatusMappings  AuthStatusMappings        `yaml:"auth_status_mappings"`
}

// ExpirationFrequencyConfig holds the expiration scheduler timing configuration.
type ExpirationFrequencyConfig struct {
	Frequency string `yaml:"frequency"`
}

// EligibleStatusesConfig holds the eligible status lists for the expiration scheduler.
type EligibleStatusesConfig struct {
	ConsentStatuses []string `yaml:"consent_statuses"`
}

// ConsentStatusMappings holds the mapping of specific consent lifecycle states.
type ConsentStatusMappings struct {
	ActiveStatus   string `yaml:"active_status"`
	ExpiredStatus  string `yaml:"expired_status"`
	RevokedStatus  string `yaml:"revoked_status"`
	CreatedStatus  string `yaml:"created_status"`
	RejectedStatus string `yaml:"rejected_status"`
}

// AuthStatusMappings holds the mapping of authorization resource lifecycle states.
type AuthStatusMappings struct {
	ApprovedState      string `yaml:"approved_state"`
	RejectedState      string `yaml:"rejected_state"`
	CreatedState       string `yaml:"created_state"`
	SystemExpiredState string `yaml:"system_expired_state"`
	SystemRevokedState string `yaml:"system_revoked_state"`
}

// ── Consent status getters ────────────────────────────────────────────────────

// GetActiveConsentStatus returns the typed active status from config.
func (c *ConsentConfig) GetActiveConsentStatus() ConsentStatus {
	return ConsentStatus(c.StatusMappings.ActiveStatus)
}

// GetExpiredConsentStatus returns the typed expired status from config.
func (c *ConsentConfig) GetExpiredConsentStatus() ConsentStatus {
	return ConsentStatus(c.StatusMappings.ExpiredStatus)
}

// GetRevokedConsentStatus returns the typed revoked status from config.
func (c *ConsentConfig) GetRevokedConsentStatus() ConsentStatus {
	return ConsentStatus(c.StatusMappings.RevokedStatus)
}

// GetCreatedConsentStatus returns the typed created status from config.
func (c *ConsentConfig) GetCreatedConsentStatus() ConsentStatus {
	return ConsentStatus(c.StatusMappings.CreatedStatus)
}

// GetRejectedConsentStatus returns the typed rejected status from config.
func (c *ConsentConfig) GetRejectedConsentStatus() ConsentStatus {
	return ConsentStatus(c.StatusMappings.RejectedStatus)
}

// ── Auth status getters ───────────────────────────────────────────────────────

// GetApprovedAuthStatus returns the typed approved auth status from config.
func (c *ConsentConfig) GetApprovedAuthStatus() AuthStatus {
	return AuthStatus(c.AuthStatusMappings.ApprovedState)
}

// GetRejectedAuthStatus returns the typed rejected auth status from config.
func (c *ConsentConfig) GetRejectedAuthStatus() AuthStatus {
	return AuthStatus(c.AuthStatusMappings.RejectedState)
}

// GetCreatedAuthStatus returns the typed created auth status from config.
func (c *ConsentConfig) GetCreatedAuthStatus() AuthStatus {
	return AuthStatus(c.AuthStatusMappings.CreatedState)
}

// GetSystemExpiredAuthStatus returns the typed system expired auth status from config.
func (c *ConsentConfig) GetSystemExpiredAuthStatus() AuthStatus {
	return AuthStatus(c.AuthStatusMappings.SystemExpiredState)
}

// GetSystemRevokedAuthStatus returns the typed system revoked auth status from config.
func (c *ConsentConfig) GetSystemRevokedAuthStatus() AuthStatus {
	return AuthStatus(c.AuthStatusMappings.SystemRevokedState)
}

// ── Consent status checks ─────────────────────────────────────────────────────

// IsStatusAllowed checks if a given status is a valid consent status.
func (c *ConsentConfig) IsStatusAllowed(status ConsentStatus) bool {
	return status == c.GetActiveConsentStatus() ||
		status == c.GetExpiredConsentStatus() ||
		status == c.GetRevokedConsentStatus() ||
		status == c.GetCreatedConsentStatus() ||
		status == c.GetRejectedConsentStatus()
}

// IsActiveStatus checks if the given status represents an active consent.
func (c *ConsentConfig) IsActiveStatus(status ConsentStatus) bool {
	return status == c.GetActiveConsentStatus()
}

// IsExpiredStatus checks if the given status represents an expired consent.
func (c *ConsentConfig) IsExpiredStatus(status ConsentStatus) bool {
	return status == c.GetExpiredConsentStatus()
}

// IsRevokedStatus checks if the given status represents a revoked consent.
func (c *ConsentConfig) IsRevokedStatus(status ConsentStatus) bool {
	return status == c.GetRevokedConsentStatus()
}

// IsCreatedStatus checks if the given status represents a created consent.
func (c *ConsentConfig) IsCreatedStatus(status ConsentStatus) bool {
	return status == c.GetCreatedConsentStatus()
}

// IsRejectedStatus checks if the given status represents a rejected consent.
func (c *ConsentConfig) IsRejectedStatus(status ConsentStatus) bool {
	return status == c.GetRejectedConsentStatus()
}

// IsTerminalStatus checks if the given status is a terminal state (expired or revoked).
func (c *ConsentConfig) IsTerminalStatus(status ConsentStatus) bool {
	return c.IsExpiredStatus(status) || c.IsRevokedStatus(status)
}

// GetAllowedConsentStatuses returns a list of all valid consent statuses.
func (c *ConsentConfig) GetAllowedConsentStatuses() []ConsentStatus {
	return []ConsentStatus{
		c.GetCreatedConsentStatus(),
		c.GetActiveConsentStatus(),
		c.GetRejectedConsentStatus(),
		c.GetRevokedConsentStatus(),
		c.GetExpiredConsentStatus(),
	}
}

// ── Auth status checks ────────────────────────────────────────────────────────

// IsAuthStatusAllowed checks if a given status is a valid authorization status.
func (c *ConsentConfig) IsAuthStatusAllowed(status AuthStatus) bool {
	return status == c.GetCreatedAuthStatus() ||
		status == c.GetApprovedAuthStatus() ||
		status == c.GetRejectedAuthStatus() ||
		status == c.GetSystemExpiredAuthStatus() ||
		status == c.GetSystemRevokedAuthStatus()
}

// GetAllowedAuthStatuses returns a list of all valid authorization statuses.
func (c *ConsentConfig) GetAllowedAuthStatuses() []AuthStatus {
	return []AuthStatus{
		c.GetCreatedAuthStatus(),
		c.GetApprovedAuthStatus(),
		c.GetRejectedAuthStatus(),
		c.GetSystemExpiredAuthStatus(),
		c.GetSystemRevokedAuthStatus(),
	}
}
