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

package client

import (
	"context"

	dba "github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/rkosegi/db2rest-bridge/pkg/query"
	"github.com/samber/lo"
)

type imCrud[T any] struct {
	items []*T
	// this function checks if given string is "key" of given item
	isKeyFn func(string, *T) bool
}

func NewInMemory[T any](isKeyFn func(string, *T) bool, initData []*T) GenericInterface[T] {
	return &imCrud[T]{isKeyFn: isKeyFn, items: initData}
}

func (i *imCrud[T]) List(_ context.Context, q query.Interface) ([]*T, int, error) {
	out := lo.Filter(i.items, func(item *T, _ int) bool {
		return FilterExpressionAsPredicate[T](q.Filter(), item)
	})
	// TODO: how about ordering?
	out = ApplyPaging[T](out, q.Paging())
	return out, len(out), nil
}

func (i *imCrud[T]) Create(_ context.Context, t *T) (*T, error) {
	i.items = append(i.items, t)
	return t, nil
}

func (i *imCrud[T]) Get(_ context.Context, key string) (*T, error) {
	if res, ok := lo.Find(i.items, func(item *T) bool {
		return i.isKeyFn(key, item)
	}); ok {
		return res, nil
	}
	return nil, nil
}

func (i *imCrud[T]) Delete(_ context.Context, key string) error {
	i.items = lo.Filter(i.items, func(item *T, _ int) bool {
		return !i.isKeyFn(key, item)
	})
	return nil
}

func (i *imCrud[T]) Update(_ context.Context, key string, t *T) (*T, error) {
	if item, ok := lo.Find(i.items, func(item *T) bool {
		return i.isKeyFn(key, item)
	}); ok {
		*item = *t
	}
	return t, nil
}

func (i *imCrud[T]) BulkUpdate(_ context.Context, _ []*T, _ dba.BulkUpdateMode) error {
	panic("implement me")
}

func (i *imCrud[T]) Query(_ context.Context, _ string, _ query.Interface, _ []string) (*dba.PagedResult, error) {
	panic("implement me")
}
