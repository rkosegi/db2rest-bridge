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
	"slices"
	"strings"

	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/samber/lo"
)

func mapValue(ct *sql.ColumnType, val interface{}) interface{} {
	switch ct.DatabaseTypeName() {
	case "DATETIME":
		x := &sql.NullTime{}
		if err := x.Scan(val); err == nil {
			return x.Time
		}
	case "DECIMAL":
		x := &sql.NullFloat64{}
		if err := x.Scan(val); err == nil {
			return x.Float64
		}
	case "VARCHAR", "ENUM", "CHAR":
		x := &sql.NullString{}
		if err := x.Scan(val); err == nil {
			return x.String
		}
	}
	// TODO: more types?
	return val
}

func getRowMetadata(rows *sql.Rows) ([]string, []*sql.ColumnType, error) {
	var (
		err      error
		cols     []string
		colTypes []*sql.ColumnType
	)
	if cols, err = rows.Columns(); err != nil {
		return nil, nil, err
	}
	if colTypes, err = rows.ColumnTypes(); err != nil {
		return nil, nil, err
	}
	return cols, colTypes, nil
}

func mapEntity(rows *sql.Rows, columns []string, columnTypes []*sql.ColumnType) (res api.UntypedDto, err error) {
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}
	res = make(api.UntypedDto, len(values))
	if err = rows.Scan(values...); err != nil {
		return nil, err
	}
	for i, column := range columns {
		res[column] = mapValue(columnTypes[i], *(values[i].(*interface{})))
	}
	return res, nil
}

func createInsertQuery(entity string, body api.UntypedDto) (string, []interface{}) {
	sb := strings.Builder{}
	csb := strings.Builder{}
	vsb := strings.Builder{}
	sb.WriteString("INSERT INTO `")
	sb.WriteString(entity)
	sb.WriteString("` ")
	cols := lo.Keys(body)
	slices.Sort(cols)
	colCount := len(cols)
	values := make([]interface{}, 0)
	csb.WriteRune('(')
	vsb.WriteString(" VALUES(")
	for i := 0; i < colCount; i++ {
		col := cols[i]
		csb.WriteRune('`')
		csb.WriteString(col)
		csb.WriteRune('`')
		vsb.WriteRune('?')
		if i < colCount-1 {
			csb.WriteRune(',')
			vsb.WriteRune(',')
		}
		values = append(values, body[col])
	}
	csb.WriteRune(')')
	vsb.WriteRune(')')
	sb.WriteString(csb.String())
	sb.WriteString(vsb.String())
	return sb.String(), values
}

func createUpdateQuery(entity, idColumn string, body api.UntypedDto) (string, []interface{}) {
	sb := strings.Builder{}
	sb.WriteString("UPDATE `")
	sb.WriteString(entity)
	sb.WriteString("` SET ")
	cols := lo.Keys(body)
	slices.Sort(cols)
	colCount := len(cols)
	values := make([]interface{}, 0)
	for i := 0; i < colCount; i++ {
		col := cols[i]
		sb.WriteRune('`')
		sb.WriteString(col)
		sb.WriteString("` = ?")
		values = append(values, body[col])
		if i < colCount-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteRune(' ')
	sb.WriteString(createSingleItemFilter(idColumn))
	return sb.String(), values
}

func createDeleteQueryPrefix(entity string) string {
	sb := strings.Builder{}
	sb.WriteString("DELETE FROM `")
	sb.WriteString(entity)
	sb.WriteString("` ")
	return sb.String()
}

// createSingleDeleteQuery generates `DELETE FROM <entity> WHERE <id> = ? LIMIT 1` query
func createSingleDeleteQuery(entity, idColumn string) string {
	sb := strings.Builder{}
	sb.WriteString(createDeleteQueryPrefix(entity))
	sb.WriteString(createSingleItemFilter(idColumn))
	return sb.String()
}

// createMultiDeleteQuery generates `DELETE FROM <entity> WHERE <id> IN (?,?,?....?)` query.
// idsCount should be > 1
func createMultiDeleteQuery(entity, idColumn string, idsCount int) string {
	sb := strings.Builder{}
	sb.WriteString(createDeleteQueryPrefix(entity))
	sb.WriteString(createMultiItemFilter(idColumn, idsCount))
	return sb.String()
}

func createSingleSelectQuery(entity, idColumn string) string {
	sb := strings.Builder{}
	sb.WriteString("SELECT * FROM `")
	sb.WriteString(entity)
	sb.WriteString("` ")
	sb.WriteString(createSingleItemFilter(idColumn))
	return sb.String()
}

func createSingleItemFilter(idColumn string) string {
	sb := strings.Builder{}
	sb.WriteString("WHERE ")
	sb.WriteRune('`')
	sb.WriteString(idColumn)
	sb.WriteRune('`')
	sb.WriteString(" = ?")
	sb.WriteString(" LIMIT 1")
	return sb.String()
}

// createMultiItemFilter generates `WHERE <id> IN (?,?,...?)` filter.
// caller must ensure that idsCount > 1, either by falling back to createSingleItemFilter for (idsCount==1) or
// by reporting error.
func createMultiItemFilter(idColumn string, idsCount int) string {
	sb := strings.Builder{}
	sb.WriteString("WHERE ")
	sb.WriteRune('`')
	sb.WriteString(idColumn)
	sb.WriteRune('`')
	sb.WriteString(" IN (")
	sb.WriteString(strings.Repeat("?,", idsCount-1))
	sb.WriteString("?)")
	return sb.String()
}
