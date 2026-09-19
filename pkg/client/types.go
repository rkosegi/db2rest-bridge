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

package client

import (
	"context"

	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/rkosegi/db2rest-bridge/pkg/query"
)

type GenericInterface[T any] interface {
	// Raw gets access to untyped API
	Raw() RawInterface
	Create(context.Context, *T) (*T, error)
	Get(context.Context, string) (*T, error)
	List(context.Context, query.Interface) ([]*T, int, error)
	Update(context.Context, string, *T) (*T, error)
	// Delete removes single entity instance by its ID
	Delete(context.Context, string) error
	BulkUpdate(context.Context, []*T, api.BulkUpdateMode) error
	Query(context.Context, string, query.Interface, []string) ([]*T, int, error)
}

type RawInterface interface {
	// Create creates new entity instance
	Create(context.Context, api.UntypedDto) (*api.UntypedDto, error)
	// Get gets single entity instance by its ID
	Get(context.Context, string) (*api.UntypedDto, error)
	// List lists entity instances based on provided query
	List(context.Context, query.Interface) (*api.PagedResult, error)
	// Update updates single entity identified by its key with values provided as map.
	// This method exists to address limitation of generic API to perform updates by i.e. skipping empty strings.
	Update(context.Context, string, api.UntypedDto) (*api.UntypedDto, error)
	// BulkUpdate perform multiple mutation operations in one transaction
	BulkUpdate(context.Context, []api.UntypedDto, api.BulkUpdateMode) error
	// Query executes named query already defined in server-side configuration and returns result.
	Query(context.Context, string, query.Interface, []string) (*api.PagedResult, error)
}

type Opt[T any] func(*generic[T])

type EncoderFn[T any] func(*T) (api.UntypedDto, error)
type DecoderFn[T any] func(dto api.UntypedDto) (*T, error)
