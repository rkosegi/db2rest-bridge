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
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jellydator/ttlcache/v3"
	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/rkosegi/db2rest-bridge/pkg/query"
	"github.com/rkosegi/db2rest-bridge/pkg/types"
	"github.com/samber/lo"
)

var (
	errNoObj        = errors.New("require at least one object")
	dateTimeLayouts = []string{time.RFC3339, time.DateOnly}
)

type impl struct {
	io.Closer
	config  *types.BackendConfig
	l       *slog.Logger
	mdCache *ttlcache.Cache[string, map[string]*sql.ColumnType]
}

type Opt func(*impl)

func WithLogger(l *slog.Logger) Opt {
	return func(i *impl) {
		i.l = l
	}
}

func newImpl(be *types.BackendConfig, opts ...Opt) Interface {
	i := &impl{config: be}
	for _, opt := range append([]Opt{
		WithLogger(slog.Default()),
	}, opts...) {
		opt(i)
	}
	i.mdCache = ttlcache.New[string, map[string]*sql.ColumnType](
		ttlcache.WithTTL[string, map[string]*sql.ColumnType](1*time.Hour),
		ttlcache.WithCapacity[string, map[string]*sql.ColumnType](250),
		ttlcache.WithLoader[string, map[string]*sql.ColumnType](i),
	)
	go i.mdCache.Start()
	return i
}

func (be *impl) Load(c *ttlcache.Cache[string, map[string]*sql.ColumnType], key string) *ttlcache.Item[string, map[string]*sql.ColumnType] {
	be.l.Debug("loading entity metadata into cache", "entity", key)
	qry := createSingleSelectQuery(key, be.config.IdColumn(key))
	be.l.Debug("SQL", "query", qry)
	rows, err := be.config.DB().Query(qry, "0")
	if err != nil {
		return nil
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var colTypes []*sql.ColumnType
	if _, colTypes, err = getRowMetadata(rows); err != nil {
		be.l.Warn("unable to fetch row metadata", "entity", key, "err", err)
		return nil
	}
	return c.Set(key, lo.Associate(colTypes, func(item *sql.ColumnType) (string, *sql.ColumnType) {
		return item.Name(), item
	}), ttlcache.DefaultTTL)
}

func (be *impl) fetchOneItem(ctx context.Context, entity, id string, retrieve bool) (res api.UntypedDto, err error) {
	qry := createSingleSelectQuery(entity, be.config.IdColumn(entity))
	be.l.Debug("SQL", "query", qry)
	rows, err := be.config.DB().QueryContext(ctx, qry, id)
	if err != nil {
		return nil, types.WrapError("failed to fetch single row", err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	if rows.Next() {
		if retrieve {
			var (
				cols     []string
				colTypes []*sql.ColumnType
			)
			if cols, colTypes, err = getRowMetadata(rows); err != nil {
				return nil, types.WrapError("failed to get row metadata", err)
			}
			return mapEntity(rows, cols, colTypes)
		}

		res = make(api.UntypedDto, 1)
	}
	return res, nil
}

func (be *impl) ListItems(ctx context.Context, entity string, qe query.Interface) (*api.PagedResult, error) {
	if !*be.config.Read {
		return nil, errReadNotAllowed
	}
	var (
		cnt int
		err error
		qry string
	)

	if qe == nil {
		qe = query.DefaultQuery
	}
	whereExpr := ""
	flt := qe.Filter()
	if flt != nil {
		whereExpr = fmt.Sprintf(" WHERE %v", flt)
	}
	qry = fmt.Sprintf("SELECT COUNT(1) FROM `%s`%s", entity, whereExpr)
	be.l.Debug("SQL", "query", qry)
	row := be.config.DB().QueryRowContext(ctx, qry)
	if err = row.Scan(&cnt); err != nil {
		return nil, types.WrapError("failed to determine resultset size", err)
	}
	res := []api.UntypedDto{}
	if cnt > 0 {
		qry = fmt.Sprintf("SELECT * FROM `%s`%s", entity, qe.String())
		be.l.Debug("SQL", "query", qry)
		if res, err = be.fetchRows(ctx, qry); err != nil {
			return nil, types.WrapError("failed to fetch rows", err)
		}
	}
	return &api.PagedResult{
		Data:       &res,
		TotalCount: &cnt,
		Offset:     lo.ToPtr(float32(qe.Paging().Offset())),
	}, nil
}

func (be *impl) QueryNamed(ctx context.Context, name string, qry query.Interface, args ...interface{}) (*api.PagedResult, error) {
	if !*be.config.Read {
		return nil, errReadNotAllowed
	}
	if qry == nil {
		qry = query.DefaultQuery
	}
	var err error
	items := make([]api.UntypedDto, 0)
	savedQry, ok := be.config.Queries[name]
	if !ok {
		return nil, types.NewErrorWithStatus("no such query: "+name, http.StatusNotFound)
	}

	countQry := fmt.Sprintf("SELECT COUNT(1) FROM (%s) AS wrapper", savedQry)
	be.l.Debug("SQL", "query", countQry)
	row := be.config.DB().QueryRowContext(ctx, countQry, args...)
	var cnt int
	if err = row.Scan(&cnt); err != nil {
		return nil, types.WrapError("failed to determine resultset size", err)
	}

	var offset uint64
	if qry.Paging() != nil {
		savedQry += " LIMIT " + qry.Paging().String()
		offset = qry.Paging().Offset()
	}

	be.l.Debug("SQL", "query", savedQry)
	if items, err = be.fetchRows(ctx, savedQry, args...); err != nil {
		return nil, types.WrapError("failed to execute query "+name, err)
	}
	return &api.PagedResult{
		TotalCount: &cnt,
		Data:       &items,
		Offset:     lo.ToPtr(float32(offset)),
	}, nil
}

func (be *impl) fetchRows(ctx context.Context, qry string, args ...interface{}) ([]api.UntypedDto, error) {
	rows, err := be.config.DB().QueryContext(ctx, qry, args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var (
		cols     []string
		colTypes []*sql.ColumnType
	)

	if cols, colTypes, err = getRowMetadata(rows); err != nil {
		return nil, types.WrapError("failed to get row metadata", err)
	}

	res := []api.UntypedDto{}
	for rows.Next() {
		var item api.UntypedDto
		if item, err = mapEntity(rows, cols, colTypes); err != nil {
			return nil, types.WrapError("failed to map row to entity", err)
		}
		res = append(res, item)
	}
	return res, nil
}

func (be *impl) ListEntities(ctx context.Context) ([]string, error) {
	if !*be.config.Read {
		return nil, errReadNotAllowed
	}
	var (
		err error
		res []string
	)
	qry := "SHOW TABLES"
	be.l.Debug("SQL", "query", qry)
	rows, err := be.config.DB().QueryContext(ctx, qry)
	if err != nil {
		return nil, types.WrapError("failed to list entity tables", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var t string
		if err = rows.Scan(&t); err != nil {
			return nil, types.WrapError("failed to scan table name", err)
		}
		res = append(res, t)
	}
	return res, nil
}

func (be *impl) Exists(ctx context.Context, entity, id string) (bool, error) {
	if !*be.config.Read {
		return false, errReadNotAllowed
	}
	r, err := be.fetchOneItem(ctx, entity, id, false)
	return r != nil, err
}

func (be *impl) Get(ctx context.Context, entity, id string) (res api.UntypedDto, err error) {
	if !*be.config.Read {
		return nil, errReadNotAllowed
	}
	return be.fetchOneItem(ctx, entity, id, true)
}

func (be *impl) Delete(ctx context.Context, entity, id string) (err error) {
	if !*be.config.Delete {
		return errDeleteNotAllowed
	}
	qry := createSingleDeleteQuery(entity, be.config.IdColumn(entity))
	be.l.Debug("SQL", "query", qry)
	_, err = be.config.DB().ExecContext(ctx, qry, id)
	return err
}

func (be *impl) Update(ctx context.Context, entity, id string, body api.UntypedDto) (api.UntypedDto, error) {
	if !*be.config.Update {
		return nil, errUpdateNotAllowed
	}
	md := be.mdCache.Get(entity)
	if md != nil {
		body = remapBody(md, body)
	}
	qry, values := createUpdateQuery(entity, be.config.IdColumn(entity), body)
	values = append(values, id)
	be.l.Debug("SQL", "query", qry)
	if _, err := be.config.DB().ExecContext(ctx, qry, values...); err != nil {
		return nil, types.WrapError("failed to update entity", err)
	}
	return be.fetchOneItem(ctx, entity, id, true)
}

func remapValue(v interface{}, ct *sql.ColumnType) interface{} {
	switch v := v.(type) {
	case string:
		if ct.DatabaseTypeName() == "DATETIME" || ct.DatabaseTypeName() == "TIMESTAMP" || ct.DatabaseTypeName() == "DATE" {
			for _, layout := range dateTimeLayouts {
				t, err := time.Parse(layout, v)
				if err == nil {
					return t
				}
			}
			return v
		}
		if ct.DatabaseTypeName() == "BLOB" {
			bytes, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				return v
			}
			return bytes
		}
	}

	return v
}

func remapBody(md *ttlcache.Item[string, map[string]*sql.ColumnType], body api.UntypedDto) api.UntypedDto {
	for key, val := range body {
		if ct, ok := md.Value()[key]; ok {
			body[key] = remapValue(val, ct)
		}
	}
	return body
}

func (be *impl) Create(ctx context.Context, entity string, body api.UntypedDto) (api.UntypedDto, error) {
	if !*be.config.Create {
		return nil, errCreateNotAllowed
	}
	var (
		err error
		id  int64
		res sql.Result
	)
	md := be.mdCache.Get(entity)
	if md != nil {
		body = remapBody(md, body)
	}
	qry, values := createInsertQuery(entity, body)
	be.l.Debug("SQL", "query", qry)
	if res, err = be.config.DB().ExecContext(ctx, qry, values...); err != nil {
		return nil, err
	}
	// TODO: this could be configurable. There are scenarios where you don't use auto increment
	if id, err = res.LastInsertId(); err != nil {
		return nil, types.WrapError("failed to retrieve last insert ID", err)
	}
	return be.fetchOneItem(ctx, entity, strconv.FormatInt(id, 10), true)
}

func (be *impl) MultiDelete(ctx context.Context, entity string, ids []interface{}) error {
	if !*be.config.Delete {
		return errDeleteNotAllowed
	}
	switch ic := len(ids); {
	case ic == 0:
		return errNoObj
	case ic == 1:
		return be.Delete(ctx, entity, fmt.Sprintf("%v", ids[0]))
	default:
		qry := createMultiDeleteQuery(entity, be.config.IdColumn(entity), ic)
		be.l.Debug("SQL", "query", qry)
		_, err := be.config.DB().ExecContext(ctx, qry, ids...)
		return err
	}
}

func (be *impl) MultiUpdate(ctx context.Context, entity string, objs []api.UntypedDto) error {
	var (
		err error
		tx  *sql.Tx
	)
	if !*be.config.Update {
		return errUpdateNotAllowed
	}
	tx, err = be.config.DB().BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	for _, obj := range objs {
		idCol := be.config.IdColumn(entity)
		id := obj[idCol].(string)
		md := be.mdCache.Get(entity)
		if md != nil {
			obj = remapBody(md, obj)
		}
		qry, values := createUpdateQuery(entity, idCol, obj)
		values = append(values, id)
		be.l.Debug("SQL", "query", qry)
		if _, err = tx.ExecContext(ctx, qry, values...); err != nil {
			be.l.ErrorContext(ctx, "query execution failed, rolling back", "err", err)
			return errors.Join(err, tx.Rollback())
		}
	}
	return tx.Commit()
}

func (be *impl) MultiCreate(ctx context.Context, entity string, replace bool, objs []api.UntypedDto) error {
	var (
		err error
		tx  *sql.Tx
	)
	if !*be.config.Create {
		return errCreateNotAllowed
	}

	tx, err = be.config.DB().BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	for _, obj := range objs {
		md := be.mdCache.Get(entity)
		if md != nil {
			obj = remapBody(md, obj)
		}
		var (
			qry    string
			values []interface{}
		)
		if replace {
			qry, values = createReplaceQuery(entity, obj)
		} else {
			qry, values = createInsertQuery(entity, obj)
		}

		be.l.Debug("SQL", "query", qry)
		if _, err = tx.ExecContext(ctx, qry, values...); err != nil {
			be.l.ErrorContext(ctx, "query execution failed, rolling back", "err", err)
			return errors.Join(types.WrapErrorWithStatus(
				"query failed: "+qry, err, http.StatusInternalServerError),
				tx.Rollback())
		}
	}

	return tx.Commit()
}
