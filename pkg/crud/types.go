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

package crud

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/rkosegi/db2rest-bridge/pkg/query"
	"github.com/rkosegi/db2rest-bridge/pkg/types"
)

type OpNotAllowedError struct {
	Op string
}

func (o OpNotAllowedError) Error() string {
	return o.Op + " is not allowed by configuration"
}

var (
	errCreateNotAllowed = OpNotAllowedError{Op: "create"}
	errReadNotAllowed   = OpNotAllowedError{Op: "read"}
	errUpdateNotAllowed = OpNotAllowedError{Op: "update"}
	errDeleteNotAllowed = OpNotAllowedError{Op: "delete"}
)

type PagedResult struct {
	Data       []api.UntypedDto `json:"data"`
	TotalCount int              `json:"total_count"`
	Offset     uint64           `json:"offset"`
}

// Interface is API to perform CRUD operation against backend
type Interface interface {
	// ListEntities lists all entity types in backend (such as tables)
	ListEntities(ctx context.Context) ([]string, error)
	// ListItems lists items based on provided query
	ListItems(ctx context.Context, entity string, qry query.Interface) (*PagedResult, error)
	// Exists checks for existence of item based on ID
	Exists(ctx context.Context, entity, id string) (bool, error)
	// Get gets item based on ID
	Get(ctx context.Context, entity, id string) (api.UntypedDto, error)
	// Update item by its ID
	Update(ctx context.Context, entity, id string, body api.UntypedDto) (api.UntypedDto, error)
	// Create creates new item
	Create(ctx context.Context, entity string, body api.UntypedDto) (api.UntypedDto, error)
	// Delete deletes item by its ID
	Delete(ctx context.Context, entity string, id string) error
	// MultiDelete deletes items that has provided ids
	MultiDelete(ctx context.Context, entity string, ids []interface{}) error
	// MultiUpdate updates multiple items in one shot
	MultiUpdate(ctx context.Context, entity string, objs []api.UntypedDto) error
	// MultiCreate creates multiple items in one shot
	// if replace is set to true, then items are removed in backend prior to creating, if they exist.
	MultiCreate(ctx context.Context, entity string, replace bool, objs []api.UntypedDto) error
}

type NameToCrudMap map[string]Interface

func (m *NameToCrudMap) Close() error {
	errs := make([]error, 0)
	for _, v := range *m {
		if err := v.(io.Closer).Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func New(be *types.BackendConfig, n string, logger *slog.Logger) Interface {
	impl := &bedb{
		config: be,
		l:      logger.With("backend", n),
	}
	c := ttlcache.New[string, map[string]*sql.ColumnType](
		ttlcache.WithTTL[string, map[string]*sql.ColumnType](1*time.Hour),
		ttlcache.WithCapacity[string, map[string]*sql.ColumnType](250),
		ttlcache.WithLoader[string, map[string]*sql.ColumnType](impl),
	)
	impl.mdCache = c
	go impl.mdCache.Start()
	return impl
}
