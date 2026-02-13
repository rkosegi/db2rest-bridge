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

	t.Run("IN expression", func(t *testing.T) {
		assert.Equal(t, "user_id IN (1,2,3)",
			In("user_id", []interface{}{1, 2, 3}).String())

		assert.Equal(t, "((salary > 1200) AND (department IN ('HR','management')))",
			Junction(OpAnd,
				SimpleExpr("salary", ">", 1200),
				In("department", []interface{}{"HR", "management"}),
			).String())
	})
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

	t.Run("IN expression (num)", func(t *testing.T) {
		fe, err = DecodeFilter(`{"in": { "name": "user_id", "val": [9,8,7]}}`)
		assert.NoError(t, err)
		assert.NotNil(t, fe)
		assert.Len(t, fe.(*inExpr).val, 3)
		assert.Equal(t, float64(9), fe.(*inExpr).val[0])
	})

	t.Run("IN expression (str)", func(t *testing.T) {
		fe, err = DecodeFilter(`{"in": { "name": "code", "val": ["x", "Y"]}}`)
		assert.NoError(t, err)
		assert.NotNil(t, fe)
		assert.Len(t, fe.(*inExpr).val, 2)
		assert.Equal(t, "x", fe.(*inExpr).val[0])
	})

	t.Run("IS NOT NULL", func(t *testing.T) {
		fe, err = DecodeFilter(`{"un": { "name": "url", "op": "IS NOT NULL"}}`)
		assert.NoError(t, err)
		assert.NotNil(t, fe)
		assert.Equal(t, OpIsNotNull, fe.(*unExpr).op)
		assert.Equal(t, "url", fe.(*unExpr).name)
		assert.Equal(t, "url IS NOT NULL", fe.String())
	})

	t.Run("BETWEEN", func(t *testing.T) {
		fe, err = DecodeFilter(`{"between":{"left":50,"name":"age","right":60}}`)
		assert.NoError(t, err)
		assert.NotNil(t, fe)
		assert.Equal(t, float64(50), fe.(BetweenExpression).Left())
		assert.Equal(t, float64(60), fe.(BetweenExpression).Right())
		assert.Equal(t, "age", fe.(BetweenExpression).Name())
		assert.Equal(t, "age BETWEEN 50 AND 60", fe.String())
	})
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
	assert.Equal(t, OpEq, fe.(*junctionExpr).sub[0].(*simpleExpr).op)

	data, err = EncodeFilter(
		Not(
			Junction(OpOr,
				SimpleExpr("name", OpEq, "John"),
				SimpleExpr("age", OpEq, 53),
			),
		),
	)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, `{"not":{"junction":{"op":"OR","sub":[{"simple":`+
		`{"name":"name","op":"=","val":"John"}},{"simple":{"name":"age","op":"=","val":53}}]}}}`,
		strings.TrimSpace(string(data)))

	t.Run("IN expression (num)", func(t *testing.T) {
		data, err = EncodeFilter(
			Not(
				In("user", []interface{}{"Bob", "Alice"}),
			),
		)
		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, `{"not":{"in":{"name":"user","val":["Bob","Alice"]}}}`,
			strings.TrimSpace(string(data)))

	})
	t.Run("simple expression with various operators", func(t *testing.T) {
		data, err = EncodeFilter(
			Junction(OpAnd,
				SimpleExpr("age", OpGt, 25),
				SimpleExpr("age", OpLt, 150),
			),
		)
		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, `{"junction":{"op":"AND","sub":[{"simple":{"name":"age","op":"\u003e","val":25}},`+
			`{"simple":{"name":"age","op":"\u003c","val":150}}]}}`,
			strings.TrimSpace(string(data)))
	})

	t.Run("unary expression IS NULL", func(t *testing.T) {
		data, err = EncodeFilter(
			UnaryExpr("url", OpIsNull),
		)
		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, `{"un":{"name":"url","op":"IS NULL"}}`, strings.TrimSpace(string(data)))
	})

	t.Run("between expression", func(t *testing.T) {
		data, err = EncodeFilter(
			BetweenExpr("age", 50, 60),
		)
		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, `{"between":{"left":50,"name":"age","right":60}}`, strings.TrimSpace(string(data)))
	})
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
	qry, err = FromParams(nil, nil, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, qry)
	req, _ = http.NewRequest(http.MethodGet, "", nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	s := `{"not":{"simple":{"op":"=","name":"name","val":"John"}}}`
	qry, err = FromParams(
		lo.ToPtr(31),
		lo.ToPtr(20),
		lo.ToPtr([]string{"name=asc", "age=desc"}),
		&s,
	)
	assert.NoError(t, err)
	assert.NotNil(t, qry)
	assert.Equal(t, OpEq, qry.Filter().(*notExpr).expr.(*simpleExpr).op)
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
