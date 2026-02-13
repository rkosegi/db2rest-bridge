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
	"fmt"
	"strings"

	"github.com/samber/lo"
)

const (
	OpAnd       = Op("AND")
	OpOr        = Op("OR")
	OpNot       = Op("NOT")
	OpEq        = Op("=")
	OpIsNull    = Op("IS NULL")
	OpIsNotNull = Op("IS NOT NULL")
	OpGt        = Op(">")
	OpLt        = Op("<")
	OpNe        = Op("<>")
	OpNe2       = Op("!=")
	OpIn        = Op("IN")
	OpLike      = Op("LIKE")
	OpNotLike   = Op("NOT LIKE")
)

var (
	DefaultPaging           = Page(DefaultPageOffset, DefaultPageSize)
	DefaultFilter           = SimpleExpr("1", OpAnd, "1")
	DefaultQuery  Interface = &qryData{
		paging: DefaultPaging,
	}
)

type order struct {
	name string
	asc  bool
}

func ord2str(asc bool) string {
	if asc {
		return "asc"
	}
	return "desc"
}

func (o *order) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf("%s=%s", o.name, ord2str(o.asc))), nil
}

func (o *order) UnmarshalText(text []byte) error {
	parts := strings.Split(string(text), "=")
	o.name = parts[0]
	if len(parts) > 1 && strings.ToLower(parts[1]) == "asc" {
		o.asc = true
	}
	return nil
}

func (o *order) Name() string {
	return o.name
}

func (o *order) Asc() bool {
	return o.asc
}

func (o Orders) String() string {
	return strings.Join(lo.Map(o, func(item Order, _ int) string {
		dir := "ASC"
		if !item.Asc() {
			dir = "DESC"
		}
		return fmt.Sprintf("%s %s", wrapStr(item.Name(), "`"), dir)
	}), ", ")
}

func OrderBy(name string, asc bool) Order {
	return &order{name: name, asc: asc}
}

type qryData struct {
	orders Orders
	paging Paging
	filter FilterExpression
}

func (q *qryData) Orders() Orders {
	return q.orders
}

func (q *qryData) Paging() Paging {
	return q.paging
}

func (q *qryData) Filter() FilterExpression {
	return q.filter
}

func (q *qryData) String() string {
	var sb strings.Builder
	if q.filter != nil {
		sb.WriteString(" WHERE ")
		sb.WriteString(q.filter.String())
	}
	if len(q.orders) > 0 {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(q.orders.String())
	}
	if q.paging != nil {
		sb.WriteString(" LIMIT ")
		sb.WriteString(q.paging.String())
	}
	return sb.String()
}

type page struct {
	offset uint64
	size   int
}

func (p page) Offset() uint64 {
	return p.offset
}

func (p page) Size() int {
	return p.size
}

func (p page) String() string {
	return fmt.Sprintf("%d, %d", p.offset, p.size)
}

func Page(offset uint64, count int) Paging {
	return &page{offset, count}
}
