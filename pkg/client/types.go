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
	List(context.Context, query.Interface) ([]*T, error)
	Create(context.Context, *T) (*T, error)
	Get(context.Context, string) (*T, error)
	Delete(context.Context, string) error
	Update(context.Context, string, *T) (*T, error)
}

type Opt[T any] func(*generic[T])

type EncoderFn[T any] func(*T) (api.UntypedDto, error)
type DecoderFn[T any] func(dto api.UntypedDto) (*T, error)
