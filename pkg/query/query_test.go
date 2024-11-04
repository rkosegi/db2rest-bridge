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

package query

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestExpressionString(t *testing.T) {
	assert.Equal(t, "((name = 'John') AND (age = 53))",
		Junction(OpAnd,
			SimpleExpr("name", "=", "John"),
			SimpleExpr("age", "=", 53)).
			String())
	assert.Equal(t, "((name = 'John') OR (name = 'Tom'))",
		Junction(OpOr,
			SimpleExpr("name", "=", "John"),
			SimpleExpr("name", "=", "Tom")).
			String())
}

func TestParseFilter(t *testing.T) {
	var (
		fe  FilterExpression
		err error
	)
	fe, err = DecodeFilter(`{"not": { "simple": {"op":">", "name": "salary", "val": 10000}}}`)
	assert.NoError(t, err)
	assert.NotNil(t, fe)
	assert.Equal(t, "salary", fe.(*notExpr).expr.(*simpleExpr).name)
	fe, err = DecodeFilter(`{"junction": {"op": "AND", "sub":[{"simple": {"name":1,"op":"=","val":1 }}]}}`)
	assert.NoError(t, err)
	assert.NotNil(t, fe)
	assert.Equal(t, OpAnd, fe.(*junctionExpr).op)
	assert.Equal(t, "1", fe.(*junctionExpr).sub[0].(*simpleExpr).name)
}

func TestOrdersString(t *testing.T) {
	assert.Equal(t, "`name` ASC, `age` DESC",
		Orders{
			OrderBy("name", true),
			OrderBy("age", false),
		}.String(),
	)
}

func TestQueryString(t *testing.T) {
	var (
		qry *qryData
	)
	qry = &qryData{}
	assert.Equal(t, "", qry.String())
	qry = &qryData{
		orders: Orders{OrderBy("name", true)},
		paging: Page(5, 10),
		filter: SimpleExpr("name", "=", "John"),
	}
	assert.Equal(t, " WHERE name = 'John' ORDER BY `name` ASC LIMIT 5, 10", qry.String())
	qry = &qryData{
		filter: Not(SimpleExpr("salary", ">", 5000)),
	}
	assert.Equal(t, " WHERE NOT (salary > 5000)", qry.String())
}

func TestPage(t *testing.T) {
	var p Paging
	p = Page(0, 100)
	assert.Equal(t, uint64(0), p.Offset())
	assert.Equal(t, 100, p.Size())
	assert.Equal(t, "0, 100", p.String())
	p = Page(100, 10)
	assert.Equal(t, "100, 10", p.String())
}

func TestDecodeEncode(t *testing.T) {
	var (
		data []byte
		err  error
		fe   FilterExpression
	)
	data, err = json.Marshal(Junction(OpAnd,
		SimpleExpr("name", "=", "John"),
		SimpleExpr("age", "=", 53)))
	assert.NoError(t, err)
	assert.NotNil(t, data)
	fe, err = DecodeFilter(string(data))
	assert.NoError(t, err)
	assert.NotNil(t, fe)
	assert.Equal(t, OpAnd, fe.(*junctionExpr).op)
	assert.Equal(t, 2, len(fe.(*junctionExpr).sub))
	assert.Equal(t, "name", fe.(*junctionExpr).sub[0].(*simpleExpr).name)
	assert.Equal(t, Op("="), fe.(*junctionExpr).sub[0].(*simpleExpr).op)

	data, err = EncodeFilter(
		Not(
			Junction(OpOr,
				SimpleExpr("name", "=", "John"),
				SimpleExpr("age", "=", 53),
			),
		),
	)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, `{"not":{"junction":{"op":"OR","sub":[{"simple":`+
		`{"name":"name","op":"=","val":"John"}},{"simple":{"name":"age","op":"=","val":53}}]}}}`,
		strings.TrimSpace(string(data)))
}

func TestToParams(t *testing.T) {
	var (
		params *api.ListItemsParams
		err    error
	)
	params, err = ToParams(&qryData{
		orders: Orders{OrderBy("name", true)},
		paging: DefaultPaging,
		filter: DefaultFilter,
	})
	assert.NoError(t, err)
	assert.NotNil(t, params)
	assert.Equal(t, 0, *params.PageOffset)
	assert.Equal(t, 20, *params.PageSize)
}

func TestDecodeRequest(t *testing.T) {
	var (
		req *http.Request
		qry Interface
		err error
	)
	qry, err = FromParams(api.ListItemsParams{})
	assert.NoError(t, err)
	assert.NotNil(t, qry)
	req, _ = http.NewRequest(http.MethodGet, "", nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	s := `{"not":{"simple":{"op":"=","name":"name","val":"John"}}}`
	qry, err = FromParams(api.ListItemsParams{
		PageOffset: lo.ToPtr(31),
		PageSize:   lo.ToPtr(20),
		Order:      lo.ToPtr([]string{"name=asc", "age=desc"}),
		Filter:     &s,
	})
	assert.NoError(t, err)
	assert.NotNil(t, qry)
	assert.Equal(t, Op("="), qry.Filter().(*notExpr).expr.(*simpleExpr).op)
	assert.Equal(t, "name", qry.Orders()[0].Name())
	assert.Equal(t, true, qry.Orders()[0].Asc())
	assert.Equal(t, "age", qry.Orders()[1].Name())
	assert.Equal(t, false, qry.Orders()[1].Asc())
}

func TestDecodeExprFromMap(t *testing.T) {
	assert.Nil(t, decodeExprFromMap(map[string]interface{}{
		"invalid": map[string]interface{}{},
	}))
}

func TestEncodeRequest(t *testing.T) {
	var (
		err error
		req *http.Request
	)
	req, err = http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	err = EncodeRequest(req, &qryData{
		orders: Orders{
			OrderBy("name", true),
			OrderBy("age", false),
		},
		paging: Page(5, 10),
		filter: Junction(OpAnd),
	})
	assert.Contains(t, req.URL.RawQuery, "order%5B%5D=name")
	assert.Contains(t, req.URL.RawQuery, "order%5B%5D=age")
	assert.NoError(t, err)
}
