/*
Copyright 2026 Richard Kosegi

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

package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	ccfg "github.com/rkosegi/go-http-commons/config"
)

var (
	TRUE            = true
	FALSE           = false
	DefaultDbDriver = Mysql
	emptyIdMap      = make(map[string]string)
	beNameRE        = regexp.MustCompile(`^[\w-]{1,63}$`)
	ErrNoBackend    = errors.New("no backend configured")
)

// IdColumn gets ID column for given entity, see IdMap
func (be *BackendConfig) IdColumn(ent string) string {
	if e, hasEnt := (*be.Entities)[ent]; hasEnt {
		if e.IdColumn != nil {
			return *e.IdColumn
		}
	}
	return "id"
}

func (be *BackendConfig) pageSizeLimitForEntity(ent string) *int {
	if e, hasEnt := (*be.Entities)[ent]; hasEnt {
		if e.PageSizeLimit != nil {
			return e.PageSizeLimit
		}
	}
	return be.PageSizeLimit
}

func (be *BackendConfig) ApplyPageSizeLimit(ent string, current int) int {
	var limit *int
	if limit = be.pageSizeLimitForEntity(ent); limit != nil {
		if *limit < current {
			return *limit
		}
	}
	return current
}

func (be *BackendConfig) Open(ctx context.Context) error {
	db, err := sql.Open(string(*be.Driver), be.Dsn)
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
	if be.InitQueries != nil {
		for _, ddl := range *be.InitQueries {
			if _, err = db.ExecContext(ctx, ddl); err != nil {
				return err
			}
		}
	}
	be.Db = db
	return nil
}

func (be *BackendConfig) Close() error {
	if be.Db != nil {
		return be.Db.Close()
	}
	return nil
}

func (be *BackendConfig) DB() *sql.DB {
	return be.Db
}

type Backends map[string]*BackendConfig

var defCorsConfig = ccfg.CorsConfig{
	MaxAge: 600,
	// if you run this in default config, you most likely come from
	// different origin then http://localhost:22001 (or whatever address is this running on).
	// Be sure to set something sane to fit your deployment.
	AllowedOrigins: []string{"*"},
}

// CheckAndNormalize sets any missing optional values and ensures all values are semantically correct.
func (c *Config) CheckAndNormalize() error {
	if c.PageSizeLimit == nil {
		c.PageSizeLimit = new(20)
	}
	if c.Server.ListenAddress == "" {
		c.Server.ListenAddress = ":22001"
	}
	if err := c.Server.Check(); err != nil {
		return err
	}
	if c.Server.Cors == nil {
		c.Server.Cors = &defCorsConfig
	}
	if len(c.Backends) == 0 {
		return ErrNoBackend
	}
	if c.Server.Telemetry == nil {
		c.Server.Telemetry = &ccfg.TelemetryConfig{
			Enabled: FALSE,
		}
	}
	if c.Server.Telemetry.Path == "" {
		c.Server.Telemetry.Path = ccfg.DefaultMetricPath
	}
	if c.Server.APIPrefix == "" {
		c.Server.APIPrefix = "/api/v1"
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
		if v.PageSizeLimit == nil {
			v.PageSizeLimit = c.PageSizeLimit
		}
		if v.Entities == nil {
			v.Entities = &map[string]*EntityConfig{}
		}
		// migrate idMap over to Entities
		for x, y := range *v.IdMap {
			if _, ok := (*v.Entities)[x]; !ok {
				(*v.Entities)[x] = &EntityConfig{}
			}
			(*v.Entities)[x].IdColumn = &y
		}
	}
	return nil
}
