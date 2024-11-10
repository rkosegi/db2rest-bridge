/*
Copyright 2024 Richard Kosegi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"
)

var (
	TRUE                  = true
	FALSE                 = false
	DefaultDbDriver       = "mysql"
	defaultLogLevel       = "info"
	defaultLogFormat      = "json"
	emptyIdMap            = make(map[string]string)
	beNameRE              = regexp.MustCompile(`^[\w-]{1,63}$`)
	ErrNoBackend          = errors.New("no backend configured")
	ErrTlsKeyPathMissing  = errors.New("server.http_tls_config.tls_key_path is required")
	ErrTlsCertPathMissing = errors.New("server.http_tls_config.tls_cert_path is required")
)

type BackendConfig struct {
	// optional name of driver, if omitted, then "mysql"  is assumed
	Driver *string `yaml:"driver,omitempty"`
	DSN    string  `yaml:"dsn"`
	Create *bool   `yaml:"create,omitempty"`
	Read   *bool   `yaml:"read,omitempty"`
	Update *bool   `yaml:"update,omitempty"`
	Delete *bool   `yaml:"delete,omitempty"`
	// Optional mapping from entity (table) name to ID column.
	// If not specified, then "id" is assumed
	IdMap *map[string]string `yaml:"id_map,omitempty"`

	MaxOpenConnections *int           `yaml:"max_open_connections,omitempty"`
	MaxIdleConnections *int           `yaml:"max_idle_connections,omitempty"`
	ConnMaxLifetime    *time.Duration `yaml:"conn_max_lifetime,omitempty"`
	ConnMaxIdleTime    *time.Duration `yaml:"conn_max_idle_time,omitempty"`

	db *sql.DB
}

// IdColumn gets ID column for given entity, see IdMap
func (be *BackendConfig) IdColumn(ent string) string {
	if col, ok := (*be.IdMap)[ent]; ok {
		return col
	}
	return "id"
}

func (be *BackendConfig) Open() error {
	db, err := sql.Open(*be.Driver, be.DSN)
	if err != nil {
		return err
	}
	if be.MaxOpenConnections != nil {
		db.SetMaxOpenConns(*be.MaxOpenConnections)
	}
	if be.MaxIdleConnections != nil {
		db.SetMaxIdleConns(*be.MaxIdleConnections)
	}
	if be.ConnMaxLifetime != nil {
		db.SetConnMaxLifetime(*be.ConnMaxLifetime)
	}
	if be.ConnMaxIdleTime != nil {
		db.SetConnMaxIdleTime(*be.ConnMaxIdleTime)
	}
	be.db = db
	return nil
}

func (be *BackendConfig) Close() error {
	if be.db != nil {
		return be.db.Close()
	}
	return nil
}

func (be *BackendConfig) DB() *sql.DB {
	return be.db
}

type Backends map[string]*BackendConfig

type CorsConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	MaxAge         int      `yaml:"max_age"`
}

var defCorsConfig = CorsConfig{
	MaxAge: 600,
	// if you run this in default config, you most likely come from
	// different origin then http://localhost:22001 (or whatever address is this running on).
	// Be sure to set something sane to fit your deployment.
	AllowedOrigins: []string{"*"},
}

type TLSConfig struct {
	TLSCertPath string `yaml:"cert_file"`
	TLSKeyPath  string `yaml:"key_file"`
	ClientAuth  string `yaml:"client_auth_type"`
	ClientCAs   string `yaml:"client_ca_file"`
}

type LoggingConfig struct {
	Level  *string `yaml:"level,omitempty"`
	Format *string `yaml:"format,omitempty"`
}

type ServerConfig struct {
	HTTPListenAddress string      `yaml:"http_listen_address"`
	HTTPTLSConfig     *TLSConfig  `yaml:"http_tls_config"`
	CorsConfig        *CorsConfig `yaml:"cors"`
}

type Config struct {
	Server        ServerConfig   `yaml:"server"`
	Backends      Backends       `yaml:"backends"`
	LoggingConfig *LoggingConfig `yaml:"logging,omitempty"`
}

// CheckAndNormalize sets any missing optional values and ensures all values are semantically correct.
func (c *Config) CheckAndNormalize() error {
	if c.Server.CorsConfig == nil {
		c.Server.CorsConfig = &defCorsConfig
	}
	if len(c.Backends) == 0 {
		return ErrNoBackend
	}

	for k, v := range c.Backends {
		if !beNameRE.MatchString(k) {
			return fmt.Errorf("invalid backend name: %s", k)
		}
		if v.Driver == nil {
			v.Driver = &DefaultDbDriver
		}
		if v.IdMap == nil {
			v.IdMap = &emptyIdMap
		}
		if v.Create == nil {
			v.Create = &FALSE
		}
		if v.Read == nil {
			v.Read = &TRUE
		}
		if v.Update == nil {
			v.Update = &FALSE
		}
		if v.Delete == nil {
			v.Delete = &FALSE
		}
	}
	if c.Server.HTTPTLSConfig != nil {
		if len(c.Server.HTTPTLSConfig.TLSKeyPath) == 0 {
			return ErrTlsKeyPathMissing
		}
		if len(c.Server.HTTPTLSConfig.TLSCertPath) == 0 {
			return ErrTlsCertPathMissing
		}
	}
	if c.LoggingConfig == nil {
		c.LoggingConfig = &LoggingConfig{Level: &defaultLogLevel, Format: &defaultLogFormat}
	}
	if c.LoggingConfig.Level == nil {
		c.LoggingConfig.Level = &defaultLogLevel
	}
	if c.LoggingConfig.Format == nil {
		c.LoggingConfig.Format = &defaultLogFormat
	}
	return nil
}
