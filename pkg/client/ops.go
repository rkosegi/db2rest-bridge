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
	"net/http"

	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/rkosegi/db2rest-bridge/pkg/query"
)

func (g *generic[T]) List(ctx context.Context, qry query.Interface) ([]*T, int, error) {
	var (
		err    error
		params *api.ListItemsParams
	)
	res := make([]*T, 0)
	if qry == nil {
		qry = query.DefaultQuery
	}
	g.l.Debug("Listing entities", "query", qry)
	params, err = query.ToParams(qry)
	if err != nil {
		return nil, 0, err
	}
	resp, err := g.c.ListItemsWithResponse(ctx, g.be, g.ent, params)
	if err != nil {
		return nil, 0, err
	}
	if err = ensureResponseCode(resp.HTTPResponse, http.StatusOK); err != nil {
		return nil, 0, err
	}
	for _, e := range *resp.JSON200.Data {
		var dto *T
		dto, err = g.decFn(e)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, dto)
	}
	return res, int(*resp.JSON200.TotalCount), nil
}

func (g *generic[T]) Create(ctx context.Context, t *T) (*T, error) {
	var (
		err error
		m   map[string]interface{}
		cir *api.CreateItemResponse
	)
	if m, err = g.encFn(t); err != nil {
		return nil, err
	}
	if cir, err = g.c.CreateItemWithResponse(ctx, g.be, g.ent, filterProps(m, g.skipProps)); err != nil {
		return nil, err
	}
	if err = ensureResponseCode(cir.HTTPResponse, http.StatusCreated); err != nil {
		return nil, err
	}
	return g.decFn(*cir.JSON201)
}

func (g *generic[T]) Get(ctx context.Context, id string) (*T, error) {
	if resp, err := g.c.GetItemByIdWithResponse(ctx, g.be, g.ent, id); err != nil {
		return nil, err
	} else {
		if err = ensureResponseCode(resp.HTTPResponse, http.StatusOK); err != nil {
			return nil, err
		}
		return g.decFn(*resp.JSON200)
	}
}

func (g *generic[T]) Delete(ctx context.Context, id string) error {
	if resp, err := g.c.DeleteItemByIdWithResponse(ctx, g.be, g.ent, id); err != nil {
		return err
	} else {
		if err = ensureResponseCode(resp.HTTPResponse, http.StatusNoContent); err != nil {
			return err
		}
	}
	return nil
}

func (g *generic[T]) Update(ctx context.Context, id string, obj *T) (*T, error) {
	var (
		err error
		m   map[string]interface{}
		cir *api.UpdateItemByIdResponse
	)
	if m, err = g.encFn(obj); err != nil {
		return nil, err
	}
	if cir, err = g.c.UpdateItemByIdWithResponse(ctx, g.be, g.ent, id, filterProps(m, g.skipProps)); err != nil {
		return nil, err
	}
	if err = ensureResponseCode(cir.HTTPResponse, http.StatusAccepted); err != nil {
		return nil, err
	}
	return g.decFn(*cir.JSON202)
}
