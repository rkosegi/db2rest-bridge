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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/samber/lo"
)

const (
	DefaultPageOffset = 0
	DefaultPageSize   = 20
	qpFilter          = "filter"
	qpPageSize        = "page-size"
	qpPageOffset      = "page-offset"
	qpOrder           = "order[]"
)

func ToParams(params Interface) (*api.ListItemsParams, error) {
	var (
		data []byte
		err  error
	)
	ret := &api.ListItemsParams{}
	if fe := params.Filter(); fe != nil {
		if data, err = json.Marshal(fe); err != nil {
			return nil, err
		}
		ret.Filter = lo.ToPtr(string(data))
	}
	if len(params.Orders()) > 0 {
		ret.Order = lo.ToPtr(lo.Map(params.Orders(), func(ord Order, _ int) string {
			r, _ := ord.(*order).MarshalText()
			return string(r)
		}))
	}
	if paging := params.Paging(); paging != nil {
		ret.PageOffset = lo.ToPtr(int(paging.Offset()))
		ret.PageSize = lo.ToPtr(paging.Size())
	}
	return ret, nil
}

func FromParams(params api.ListItemsParams) (Interface, error) {
	var (
		orders Orders
		err    error
		filter FilterExpression
	)
	filter = DefaultFilter
	pageOffset := DefaultPageOffset
	pageSize := DefaultPageSize
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}
	if params.PageOffset != nil {
		pageOffset = *params.PageOffset
	}
	if params.Order != nil {
		for _, orderStr := range *params.Order {
			var ord order
			if err = ord.UnmarshalText([]byte(orderStr)); err != nil {
				return nil, err
			}
			orders = append(orders, &ord)
		}
	}
	if params.Filter != nil {
		if filter, err = DecodeFilter(*params.Filter); err != nil {
			return nil, err
		}
	}
	qry := &qryData{
		paging: Page(uint64(pageOffset), pageSize),
		orders: orders,
		filter: filter,
	}
	return qry, nil
}

func EncodeRequest(req *http.Request, qry Interface) error {
	var (
		err error
	)
	q := req.URL.Query()
	if qry.Paging() != nil {
		q.Set(qpPageOffset, strconv.FormatUint(qry.Paging().Offset(), 10))
		q.Set(qpPageSize, strconv.Itoa(qry.Paging().Size()))
	}
	if qry.Orders() != nil {
		for _, ord := range qry.Orders() {
			var data []byte
			if data, err = ord.(*order).MarshalText(); err != nil {
				return err
			}
			q.Add(qpOrder, string(data))
		}
	}
	if qry.Filter() != nil {
		var buff strings.Builder
		if err = json.NewEncoder(&buff).Encode(qry.Filter()); err != nil {
			return err
		}
		q.Set(qpFilter, buff.String())
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
	return nil
}

func decodeExprFromMap(m map[string]interface{}) FilterExpression {
	var (
		sub map[string]interface{}
		ok  bool
	)
	if sub, ok = hasSubMap("simple", m); ok {
		return simpleExprFromMap(sub)
	}
	if sub, ok = hasSubMap("junction", m); ok {
		return junctionExprFromMap(sub)
	}
	if sub, ok = hasSubMap("not", m); ok {
		return notExprFromMap(sub)
	}
	if sub, ok = hasSubMap("in", m); ok {
		return inExprFromMap(sub)
	}
	return nil
}

func hasSubMap(key string, m map[string]interface{}) (inner map[string]interface{}, ok bool) {
	if _, ok = m[key]; ok {
		if inner, ok = m[key].(map[string]interface{}); ok && inner != nil {
			return inner, ok
		}
	}
	return nil, false
}

func junctionExprFromMap(m map[string]interface{}) FilterExpression {
	op := Op(m["op"].(string))
	sub := m["sub"].([]interface{})
	return Junction(op, lo.Map(sub, func(item interface{}, _ int) FilterExpression {
		return decodeExprFromMap(item.(map[string]interface{}))
	})...)
}

func simpleExprFromMap(m map[string]interface{}) FilterExpression {
	return SimpleExpr(
		fmt.Sprintf("%v", m["name"]),
		Op(m["op"].(string)),
		fmt.Sprintf("%v", m["val"]))
}

func inExprFromMap(m map[string]interface{}) FilterExpression {
	return In(fmt.Sprintf("%v", m["name"]), m["val"].([]interface{}))
}

func notExprFromMap(m map[string]interface{}) FilterExpression {
	return Not(decodeExprFromMap(m))
}

func DecodeFilter(str string) (FilterExpression, error) {
	var (
		m   map[string]interface{}
		err error
	)
	err = json.NewDecoder(strings.NewReader(str)).Decode(&m)
	return decodeExprFromMap(m), err
}

func EncodeFilter(fe FilterExpression) ([]byte, error) {
	var (
		data bytes.Buffer
		err  error
	)
	err = json.NewEncoder(&data).Encode(fe)
	return data.Bytes(), err
}
