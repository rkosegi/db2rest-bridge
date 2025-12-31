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

type bedb struct {
	io.Closer
	config  *types.BackendConfig
	l       *slog.Logger
	mdCache *ttlcache.Cache[string, map[string]*sql.ColumnType]
}

func (be *bedb) Load(c *ttlcache.Cache[string, map[string]*sql.ColumnType], key string) *ttlcache.Item[string, map[string]*sql.ColumnType] {
	be.l.Debug("loading entity metadata into cache", "entity", key)
	qry := be.logSQL(createSingleSelectQuery(key, be.config.IdColumn(key)))

	rows, err := be.config.DB().Query(qry, "0")
	if err != nil {
		return nil
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	if _, colTypes, err := getRowMetadata(rows); err != nil {
		be.l.Warn("unable to fetch row metadata", "entity", key, "err", err)
		return nil
	} else {
		return c.Set(key, lo.Associate(colTypes, func(item *sql.ColumnType) (string, *sql.ColumnType) {
			return item.Name(), item
		}), ttlcache.DefaultTTL)
	}
}

func (be *bedb) logSQL(sql string) string {
	be.l.Debug("SQL", "query", sql)
	return sql
}

func (be *bedb) fetchOneItem(entity, id string, retrieve bool) (res api.UntypedDto, err error) {
	qry := be.logSQL(createSingleSelectQuery(entity, be.config.IdColumn(entity)))
	rows, err := be.config.DB().Query(qry, id)
	if err != nil {
		return nil, err
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
				return nil, err
			}
			return mapEntity(rows, cols, colTypes)
		} else {
			res = make(api.UntypedDto, 1)
		}
	}
	return res, nil
}

func (be *bedb) ListItems(entity string, qe query.Interface) (*PagedResult, error) {
	if !*be.config.Read {
		return nil, errReadNotAllowed
	}
	var (
		cols     []string
		cnt      int
		colTypes []*sql.ColumnType
		err      error
		rows     *sql.Rows
		res      []api.UntypedDto
		item     api.UntypedDto
		qry      string
	)

	res = make([]api.UntypedDto, 0)
	if qe == nil {
		qe = query.DefaultQuery
	}
	qry = be.logSQL(fmt.Sprintf("SELECT COUNT(1) FROM `%s` WHERE %s", entity, qe.Filter().String()))
	row := be.config.DB().QueryRow(qry)
	if err = row.Scan(&cnt); err != nil {
		return nil, err
	}

	qry = be.logSQL(fmt.Sprintf("SELECT * FROM `%s` %s", entity, qe.String()))
	rows, err = be.config.DB().Query(qry)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	cols, colTypes, err = getRowMetadata(rows)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item, err = mapEntity(rows, cols, colTypes)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}

	return &PagedResult{
		Data:       res,
		TotalCount: cnt,
		Offset:     qe.Paging().Offset(),
	}, nil
}

func (be *bedb) ListEntities() ([]string, error) {
	if !*be.config.Read {
		return nil, errReadNotAllowed
	}
	var (
		err error
		res []string
	)
	rows, err := be.config.DB().Query(be.logSQL("SHOW TABLES"))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var t string
		err = rows.Scan(&t)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}

func (be *bedb) Exists(entity, id string) (bool, error) {
	if !*be.config.Read {
		return false, errReadNotAllowed
	}
	r, err := be.fetchOneItem(entity, id, false)
	return r != nil, err
}

func (be *bedb) Get(entity, id string) (res api.UntypedDto, err error) {
	if !*be.config.Read {
		return nil, errReadNotAllowed
	}
	return be.fetchOneItem(entity, id, true)
}

func (be *bedb) Delete(entity, id string) (err error) {
	if !*be.config.Delete {
		return errDeleteNotAllowed
	}
	qry := be.logSQL(createSingleDeleteQuery(entity, be.config.IdColumn(entity)))
	_, err = be.config.DB().Exec(qry, id)
	return err
}

func (be *bedb) Update(entity, id string, body api.UntypedDto) (api.UntypedDto, error) {
	if !*be.config.Update {
		return nil, errUpdateNotAllowed
	}
	md := be.mdCache.Get(entity)
	if md != nil {
		body = remapBody(md, body)
	}
	qry, values := createUpdateQuery(entity, be.config.IdColumn(entity), body)
	values = append(values, id)
	if _, err := be.config.DB().Exec(be.logSQL(qry), values...); err != nil {
		return nil, err
	}
	return be.fetchOneItem(entity, id, true)
}

func remapValue(v interface{}, ct *sql.ColumnType) interface{} {
	switch v := v.(type) {
	case string:
		if ct.DatabaseTypeName() == "DATETIME" || ct.DatabaseTypeName() == "TIMESTAMP" {
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

func (be *bedb) Create(entity string, body api.UntypedDto) (api.UntypedDto, error) {
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
	if res, err = be.config.DB().Exec(be.logSQL(qry), values...); err != nil {
		return nil, err
	}
	if id, err = res.LastInsertId(); err != nil {
		return nil, err
	}
	return be.fetchOneItem(entity, strconv.FormatInt(id, 10), true)
}

func (be *bedb) MultiDelete(entity string, ids []interface{}) error {
	if !*be.config.Delete {
		return errDeleteNotAllowed
	}
	switch ic := len(ids); {
	case ic == 0:
		return errNoObj
	case ic == 1:
		return be.Delete(entity, fmt.Sprintf("%v", ids[0]))
	default:
		_, err := be.config.DB().Exec(be.logSQL(
			createMultiDeleteQuery(entity, be.config.IdColumn(entity), ic)), ids...)
		return err
	}
}

func (be *bedb) MultiUpdate(entity string, objs []api.UntypedDto) error {
	var (
		err error
		tx  *sql.Tx
	)
	if !*be.config.Update {
		return errUpdateNotAllowed
	}
	ctx := context.Background()
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
		if _, err = tx.ExecContext(ctx, be.logSQL(qry), values...); err != nil {
			be.l.Warn("query execution failed, rollin back", "err", err)
			return tx.Rollback()
		}
	}
	return tx.Commit()
}

func (be *bedb) MultiCreate(entity string, replace bool, objs []api.UntypedDto) error {
	var (
		err error
		tx  *sql.Tx
	)
	if !*be.config.Create {
		return errCreateNotAllowed
	}

	ctx := context.Background()
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

		if _, err = tx.ExecContext(ctx, be.logSQL(qry), values...); err != nil {
			be.l.Warn("query execution failed, rolling back", "err", err)
			return tx.Rollback()
		}
	}

	return tx.Commit()
}
