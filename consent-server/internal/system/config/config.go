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

// Package config provides structures and functions for loading and managing server configurations.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wso2/openfgc/internal/system/log"
	"gopkg.in/yaml.v3"
)

// globalConfig holds the application configuration.
var globalConfig *Config

// Config holds all configuration for the application.
type Config struct {
	Server   ServerConfig    `yaml:"server"`
	Database DatabasesConfig `yaml:"database"`
	Logging  LoggingConfig   `yaml:"logging"`
	Consent  ConsentConfig   `yaml:"consent"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Hostname     string        `yaml:"hostname"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	IdleTimeout  time.Duration `yaml:"idleTimeout"`
}

// DatabasesConfig holds all database configurations.
type DatabasesConfig struct {
	Consent DatabaseConfig `yaml:"consent"`
}

// DatabaseConfig holds individual database configuration.
type DatabaseConfig struct {
	Type            string        `yaml:"type"`
	Hostname        string        `yaml:"hostname"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Database        string        `yaml:"database"`
	Path            string        `yaml:"path"`
	SSLMode         string        `yaml:"sslmode"`
	Options         string        `yaml:"options"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// LoggingConfig holds logging configuration.
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

// Load reads configuration from file and environment variables.
func Load(configPath string) (*Config, error) {
	logger := log.GetLogger()
	logger.Debug("Loading configuration", log.String("config_path", configPath))

	var finalPath string
	if configPath != "" {
		finalPath = configPath
	} else {
		paths := []string{
			"./repository/conf/deployment.yaml",
			"./cmd/server/repository/conf/deployment.yaml",
			"../repository/conf/deployment.yaml",
			"./deployment.yaml",
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				finalPath = path
				break
			}
		}
		if finalPath == "" {
			return nil, fmt.Errorf("no configuration file found in default paths")
		}
	}

	finalPath = filepath.Clean(finalPath)
	data, err := os.ReadFile(finalPath)
	if err != nil {
		logger.Error("Failed to read config file", log.Error(err))
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	data, err = substituteEnvironmentVariables(data)
	if err != nil {
		logger.Error("Failed to substitute environment variables", log.Error(err))
		return nil, fmt.Errorf("failed to substitute environment variables: %w", err)
	}

	logger.Info("Config file loaded", log.String("file", finalPath))

	var config Config
	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	decoder.KnownFields(true)
	if err := decoder.Decode(&config); err != nil {
		logger.Error("Failed to unmarshal config", log.Error(err))
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		logger.Error("Config validation failed", log.Error(err))
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	globalConfig = &config
	logger.Debug("Configuration loaded and validated successfully",
		log.String("server_port", fmt.Sprintf("%d", config.Server.Port)),
		log.String("db_host", config.Database.Consent.Hostname),
	)
	return &config, nil
}

// substituteEnvironmentVariables replaces ${VAR_NAME} patterns with environment variable values.
func substituteEnvironmentVariables(data []byte) ([]byte, error) {
	content := string(data)
	pos := 0
	for pos < len(content) {
		start := strings.Index(content[pos:], "${")
		if start == -1 {
			break
		}
		start += pos
		end := strings.Index(content[start:], "}")
		if end == -1 {
			return nil, fmt.Errorf("unclosed environment variable substitution at position %d", start)
		}
		end += start
		varName := content[start+2 : end]
		varValue := os.Getenv(varName)
		content = content[:start] + varValue + content[end+1:]
		pos = start + len(varValue)
	}
	return []byte(content), nil
}

// validateConfig validates the configuration.
func validateConfig(config *Config) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	switch config.Database.Consent.Type {
	case "sqlite":
		if config.Database.Consent.Path == "" {
			return fmt.Errorf("database path is required for SQLite")
		}
	case "postgres":
		if config.Database.Consent.Hostname == "" {
			return fmt.Errorf("database hostname is required")
		}
		if config.Database.Consent.Database == "" {
			return fmt.Errorf("database name is required")
		}
	default:
		if config.Database.Consent.Hostname == "" {
			return fmt.Errorf("database hostname is required")
		}
		if config.Database.Consent.Database == "" {
			return fmt.Errorf("database name is required")
		}
	}

	if err := validateConsentConfig(&config.Consent); err != nil {
		return err
	}

	return nil
}

// validateConsentConfig validates consent-specific configuration.
// Extracted to keep validateConfig readable and to co-locate consent validation with consent_config.go.
func validateConsentConfig(c *ConsentConfig) error {
	if c.StatusMappings.ActiveStatus == "" {
		return fmt.Errorf("consent active status mapping is required")
	}
	if c.StatusMappings.ExpiredStatus == "" {
		return fmt.Errorf("consent expired status mapping is required")
	}
	if c.StatusMappings.RevokedStatus == "" {
		return fmt.Errorf("consent revoked status mapping is required")
	}
	if c.StatusMappings.CreatedStatus == "" {
		return fmt.Errorf("consent created status mapping is required")
	}
	if c.StatusMappings.RejectedStatus == "" {
		return fmt.Errorf("consent rejected status mapping is required")
	}
	if c.AuthStatusMappings.ApprovedState == "" {
		return fmt.Errorf("auth approved status mapping is required")
	}
	if c.AuthStatusMappings.RejectedState == "" {
		return fmt.Errorf("auth rejected status mapping is required")
	}
	if c.AuthStatusMappings.CreatedState == "" {
		return fmt.Errorf("auth created status mapping is required")
	}
	if c.AuthStatusMappings.SystemExpiredState == "" {
		return fmt.Errorf("auth system expired status mapping is required")
	}
	if c.AuthStatusMappings.SystemRevokedState == "" {
		return fmt.Errorf("auth system revoked status mapping is required")
	}
	if c.ExpirationFrequency != "" {
		if _, err := time.ParseDuration(c.ExpirationFrequency); err != nil {
			return fmt.Errorf("invalid consent expiration_frequency %q: must be a valid duration (e.g. '5m', '1h'): %w",
				c.ExpirationFrequency, err)
		}
	}
	return nil
}

// Get returns the global configuration.
func Get() *Config {
	return globalConfig
}

// SetGlobal sets the global configuration (for testing purposes).
func SetGlobal(cfg *Config) {
	globalConfig = cfg
}

// GetDSN returns the database connection string for the configured database type.
func (d *DatabaseConfig) GetDSN() string {
	switch d.Type {
	case "sqlite":
		options := d.Options
		if options != "" && options[0] != '?' {
			options = "?" + options
		}
		return d.Path + options
	case "postgres":
		escapePGValue := func(v string) string {
			v = strings.ReplaceAll(v, `\`, `\\`)
			v = strings.ReplaceAll(v, `'`, `\'`)
			return "'" + v + "'"
		}
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
			escapePGValue(d.Hostname),
			d.Port,
			escapePGValue(d.User),
			escapePGValue(d.Password),
			escapePGValue(d.Database),
		)
		if d.SSLMode != "" {
			dsn += " sslmode=" + escapePGValue(d.SSLMode)
		}
		if d.Options != "" {
			dsn += " " + d.Options
		}
		return dsn
	default: // mysql
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true",
			d.User,
			d.Password,
			d.Hostname,
			d.Port,
			d.Database,
		)
	}
}

// GetServerAddress returns the server address in host:port format.
func (s *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", s.Hostname, s.Port)
}
