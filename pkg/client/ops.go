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
	"github.com/rkosegi/db2rest-bridge/pkg/types"
	"github.com/samber/lo"
)

func (g *generic[T]) List(ctx context.Context, qry query.Interface) ([]*T, int, error) {
	resp, err := g.r.List(ctx, qry)
	if err != nil {
		return nil, 0, err
	}
	res, err := lo.MapErr(*resp.Data, func(item api.UntypedDto, _ int) (*T, error) {
		return g.decFn(item)
	})
	if err != nil {
		return nil, 0, err
	}

	return res, *resp.TotalCount, nil
}

func (g *generic[T]) Create(ctx context.Context, t *T) (*T, error) {
	var (
		err error
		m   map[string]any
	)
	if m, err = g.encFn(t); err != nil {
		return nil, err
	}
	res, err := g.r.Create(ctx, m)
	if err != nil {
		return nil, err
	}
	return g.decFn(*res)
}

func (g *generic[T]) Get(ctx context.Context, id string) (*T, error) {
	resp, err := g.r.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return g.decFn(*resp)
}

func (g *generic[T]) Update(ctx context.Context, id string, obj *T) (*T, error) {
	var (
		err error
		m   map[string]any
		res *api.UntypedDto
	)
	if m, err = g.encFn(obj); err != nil {
		return nil, err
	}

	if res, err = g.r.Update(ctx, id, m); err != nil {
		return nil, err
	}
	return g.decFn(*res)
}

func (g *generic[T]) Delete(ctx context.Context, id string) error {
	return g.r.Delete(ctx, id)
}

func (g *generic[T]) BulkUpdate(ctx context.Context, objs []*T, mode api.BulkUpdateMode) error {
	dtos, err := lo.MapErr(objs, func(item *T, _ int) (api.UntypedDto, error) {
		return g.encFn(item)
	})
	if err != nil {
		return err
	}
	return g.r.BulkUpdate(ctx, dtos, mode)
}

func (g *generic[T]) Query(ctx context.Context, s string, q query.Interface, args []string) ([]*T, int, error) {
	pr, err := g.r.Query(ctx, s, q, args)
	if err != nil {
		return nil, 0, err
	}
	res, err := lo.MapErr(*pr.Data, func(item api.UntypedDto, _ int) (*T, error) {
		return g.decFn(item)
	})
	if err != nil {
		return nil, 0, err
	}
	return res, *pr.TotalCount, nil
}

func (g *generic[T]) Raw() RawInterface {
	return g.r
}

func (r *rawImpl) List(ctx context.Context, qry query.Interface) (*api.PagedResult, error) {
	var (
		err    error
		params *api.ListItemsParams
	)
	if qry == nil {
		qry = query.DefaultQuery
	}
	r.l.Debug("Listing instances", "query", qry.String())
	params, err = query.ToParams(qry)
	if err != nil {
		return nil, err
	}
	resp, err := r.c.ListItemsWithResponse(ctx, r.be, r.ent, params)
	if err != nil {
		return nil, err
	}
	if err = ensureResponseCode(resp.HTTPResponse, http.StatusOK, resp.Body); err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (r *rawImpl) Create(ctx context.Context, dto api.UntypedDto) (*api.UntypedDto, error) {
	var (
		err error
		cir *api.CreateItemResponse
	)
	r.l.DebugContext(ctx, "Creating instance")
	if cir, err = r.c.CreateItemWithResponse(ctx, r.be, r.ent, excludeProps(dto, r.roProps)); err != nil {
		return nil, err
	}
	switch cir.StatusCode() {
	case http.StatusCreated:
		return cir.JSON201, nil
	case http.StatusInternalServerError:
		return nil, &types.ErrorWithStatus{
			Status: http.StatusInternalServerError,
			Msg:    cir.JSON500.Message,
		}
	default:
		return nil, errorFromResponse(cir.HTTPResponse)
	}
}

func (r *rawImpl) Get(ctx context.Context, id string) (*api.UntypedDto, error) {
	r.l.DebugContext(ctx, "Getting instance", "id", id)
	if resp, err := r.c.GetItemByIdWithResponse(ctx, r.be, r.ent, id); err != nil {
		return nil, err
	} else {
		switch resp.StatusCode() {
		case http.StatusOK:
			return resp.JSON200, nil
		case http.StatusNotFound:
			return nil, errorFromResponseWithMsg(resp.HTTPResponse, resp.JSON404.Message)
		default:
			return nil, errorFromResponse(resp.HTTPResponse)
		}
	}
}

func (r *rawImpl) Query(ctx context.Context, name string, qry query.Interface, args []string) (*api.PagedResult, error) {
	var (
		resp *api.QueryNamedResponse
		err  error
	)
	if qry.Paging() == nil {
		qry = query.DefaultQuery
	}
	r.l.DebugContext(ctx, "Running named query", "name", name)
	if resp, err = r.c.QueryNamedWithResponse(ctx, r.be, name, &api.QueryNamedParams{
		PageSize:   new(qry.Paging().Size()),
		PageOffset: new(api.PageOffset(qry.Paging().Offset())),
		Arg:        &args,
	}); err != nil {
		return nil, err
	}
	switch resp.HTTPResponse.StatusCode {
	case http.StatusOK:
		return resp.JSON200, nil
	case http.StatusNotFound:
		return nil, errorFromResponseWithMsg(resp.HTTPResponse, resp.JSON404.Message)
	default:
		return nil, errorFromResponse(resp.HTTPResponse)
	}
}

func (r *rawImpl) Update(ctx context.Context, id string, dto api.UntypedDto) (*api.UntypedDto, error) {
	var (
		err  error
		resp *api.UpdateItemByIdResponse
	)
	r.l.DebugContext(ctx, "Updating instance", "id", id)
	if resp, err = r.c.UpdateItemByIdWithResponse(ctx, r.be, r.ent, id, excludeProps(dto, r.roProps)); err != nil {
		return nil, err
	}
	switch resp.StatusCode() {
	case http.StatusAccepted:
		return resp.JSON202, nil
	case http.StatusNotFound:
		return nil, errorFromResponseWithMsg(resp.HTTPResponse, resp.JSON404.Message)
	default:
		return nil, errorFromResponse(resp.HTTPResponse)
	}
}

func (r *rawImpl) BulkUpdate(ctx context.Context, objs []api.UntypedDto, mode api.BulkUpdateMode) error {
	var (
		resp *api.BulkUpdateResponse
		err  error
	)
	encObjs := make([]api.UntypedDto, 0)
	for _, obj := range objs {
		switch mode {
		case api.DELETE:
			encObjs = append(encObjs, onlyProps(obj, []string{r.idProp}))
		default:
			encObjs = append(encObjs, excludeProps(obj, r.roProps))
		}
	}
	if resp, err = r.c.BulkUpdateWithResponse(ctx, r.be, r.ent, api.BulkUpdateRequest{
		Mode:    mode,
		Objects: encObjs,
	}); err != nil {
		return err
	}
	return ensureResponseCode(resp.HTTPResponse, http.StatusOK, resp.Body)
}

func (r *rawImpl) Delete(ctx context.Context, id string) error {
	r.l.DebugContext(ctx, "Deleting instance", "id", id)
	if resp, err := r.c.DeleteItemByIdWithResponse(ctx, r.be, r.ent, id); err != nil {
		return err
	} else {
		switch resp.StatusCode() {
		case http.StatusNoContent:
			return nil
		default:
			return errorFromResponse(resp.HTTPResponse)
		}
	}
}
