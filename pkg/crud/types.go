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
	"io"
	"log/slog"

	"github.com/rkosegi/db2rest-bridge/pkg/query"
	"github.com/rkosegi/db2rest-bridge/pkg/types"
	"go.uber.org/multierr"
)

type Untyped map[string]interface{}

// Interface is API to perform CRUD operation against backend
type Interface interface {
	// ListEntities lists all entity types in backend (such as tables)
	ListEntities() ([]string, error)
	// ListItems lists items based on provided query
	ListItems(entity string, qry query.Interface) ([]Untyped, error)
	// Exists checks for existence of item based on ID
	Exists(entity, id string) (bool, error)
	// Get gets item based on ID
	Get(entity, id string) (Untyped, error)
	// Update item by its ID
	Update(entity, id string, body Untyped) (Untyped, error)
	// Create creates new item
	Create(entity string, body Untyped) (Untyped, error)
	// Delete deletes item by its ID
	Delete(entity string, id string) error
}

type NameToCrudMap map[string]Interface

func (m *NameToCrudMap) Close() error {
	errs := make([]error, 0)
	for _, v := range *m {
		if err := v.(io.Closer).Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return multierr.Combine(errs...)
}

func New(be *types.BackendConfig, n string, logger *slog.Logger) Interface {
	return &bedb{
		config: be,
		l:      logger.With("backend", n),
	}
}
