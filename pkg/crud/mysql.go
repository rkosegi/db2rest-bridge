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
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rkosegi/db2rest-bridge/pkg/query"
	"github.com/rkosegi/db2rest-bridge/pkg/types"
)

var (
	errCreateNotAllowed = errors.New("create is not allowed")
	errReadNotAllowed   = errors.New("read is not allowed")
	errUpdateNotAllowed = errors.New("update is not allowed")
	errDeleteNotAllowed = errors.New("delete is not allowed")
)

type bedb struct {
	io.Closer
	config *types.BackendConfig
	l      *slog.Logger
}

func (c *bedb) logSQL(sql string) string {
	c.l.Debug("SQL", "query", sql)
	return sql
}

func (c *bedb) fetchOneItem(entity, id string, retrieve bool) (res Untyped, err error) {
	qry := c.logSQL(createSingleSelectQuery(entity, c.config.IdColumn(entity)))
	rows, err := c.config.DB().Query(qry, id)
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
			res = make(Untyped, 1)
		}
	}
	return res, nil
}

func (c *bedb) ListItems(entity string, qe query.Interface) (*PagedResult, error) {
	if !*c.config.Read {
		return nil, errReadNotAllowed
	}
	var (
		cols     []string
		cnt      int
		colTypes []*sql.ColumnType
		err      error
		rows     *sql.Rows
		res      []Untyped
		item     Untyped
		qry      string
	)

	res = make([]Untyped, 0)
	if qe == nil {
		qe = query.DefaultQuery
	}
	qry = c.logSQL(fmt.Sprintf("SELECT COUNT(1) FROM `%s` WHERE %s", entity, qe.Filter().String()))
	row := c.config.DB().QueryRow(qry)
	if err = row.Scan(&cnt); err != nil {
		return nil, err
	}

	qry = c.logSQL(fmt.Sprintf("SELECT * FROM `%s` %s", entity, qe.String()))
	rows, err = c.config.DB().Query(qry)
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

func (c *bedb) ListEntities() ([]string, error) {
	if !*c.config.Read {
		return nil, errReadNotAllowed
	}
	var (
		err error
		res []string
	)
	rows, err := c.config.DB().Query(c.logSQL("SHOW TABLES"))
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

func (c *bedb) Exists(entity, id string) (bool, error) {
	if !*c.config.Read {
		return false, errReadNotAllowed
	}
	r, err := c.fetchOneItem(entity, id, false)
	return r != nil, err
}

func (c *bedb) Get(entity, id string) (res Untyped, err error) {
	if !*c.config.Read {
		return nil, errReadNotAllowed
	}
	return c.fetchOneItem(entity, id, true)
}

func (c *bedb) Delete(entity, id string) (err error) {
	if !*c.config.Delete {
		return errDeleteNotAllowed
	}
	qry := c.logSQL(createDeleteQuery(entity, c.config.IdColumn(entity)))
	_, err = c.config.DB().Exec(qry, id)
	return err
}

func (c *bedb) Update(entity, id string, body Untyped) (Untyped, error) {
	if !*c.config.Update {
		return nil, errUpdateNotAllowed
	}
	qry, values := createUpdateQuery(entity, c.config.IdColumn(entity), body)
	values = append(values, id)
	if _, err := c.config.DB().Exec(c.logSQL(qry), values...); err != nil {
		return nil, err
	}
	return c.fetchOneItem(entity, id, true)
}

func (c *bedb) Create(entity string, body Untyped) (Untyped, error) {
	if !*c.config.Create {
		return nil, errCreateNotAllowed
	}
	var (
		err error
		id  int64
		res sql.Result
	)
	qry, values := createInsertQuery(entity, body)
	if res, err = c.config.DB().Exec(c.logSQL(qry), values...); err != nil {
		return nil, err
	}
	if id, err = res.LastInsertId(); err != nil {
		return nil, err
	}
	return c.fetchOneItem(entity, strconv.FormatInt(id, 10), true)
}
