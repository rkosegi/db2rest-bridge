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
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type simpleExpr struct {
	name string
	op   Op
	val  interface{}
}

func (s simpleExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"simple": map[string]interface{}{
			"name": s.Name(),
			"op":   s.Op(),
			"val":  s.Value(),
		},
	})
}

func (s simpleExpr) Name() string {
	return s.name
}

func (s simpleExpr) Op() Op {
	return s.op
}

func (s simpleExpr) Value() interface{} {
	return s.val
}

func wrapStr(v interface{}, q string) interface{} {
	if _, ok := v.(string); ok {
		return q + v.(string) + q
	}
	return v
}

func (s simpleExpr) String() string {
	return fmt.Sprintf(`%s %s %v`, s.Name(), s.Op(), wrapStr(s.Value(), "'"))
}

type junctionExpr struct {
	op  Op
	sub []FilterExpression
}

func (j junctionExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"junction": map[string]interface{}{
			"sub": j.Sub(),
			"op":  j.Op(),
		},
	})
}

func (j junctionExpr) Op() Op {
	return j.op
}

func (j junctionExpr) Sub() []FilterExpression {
	return j.sub
}

func (j junctionExpr) String() string {
	var sb strings.Builder
	sb.WriteRune('(')
	sb.WriteString(strings.Join(lo.Map(j.Sub(), func(item FilterExpression, _ int) string {
		return "(" + item.String() + ")"
	}), " "+string(j.Op())+" "))
	sb.WriteRune(')')

	return sb.String()
}

type notExpr struct {
	expr FilterExpression
}

func (n notExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"not": n.Sub(),
	})
}

func (n notExpr) Sub() FilterExpression {
	return n.expr
}

func (n notExpr) String() string {
	return fmt.Sprintf("NOT (%s)", n.expr.String())
}

func SimpleExpr(name string, op Op, val interface{}) FilterExpression {
	return &simpleExpr{name: name, op: op, val: val}
}

func Junction(op Op, sub ...FilterExpression) FilterExpression {
	return &junctionExpr{op: op, sub: sub}
}

func Not(expr FilterExpression) FilterExpression {
	return &notExpr{expr: expr}
}

type inExpr struct {
	name string
	val  []interface{}
}

func (i inExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"in": map[string]interface{}{
			"name": i.Name(),
			"val":  i.Values(),
		},
	})
}

func (i inExpr) Name() string {
	return i.name
}

func (i inExpr) Values() []interface{} {
	return i.val
}

func (i inExpr) String() string {
	return fmt.Sprintf("%s IN (%s)", i.name, strings.Join(lo.Map(i.val, func(item interface{}, _ int) string {
		return fmt.Sprintf("%v", wrapStr(item, "'"))
	}), ","))
}

func In(name string, vals []interface{}) FilterExpression {
	return &inExpr{
		name: name,
		val:  vals,
	}
}

type unExpr struct {
	name string
	op   Op
}

func (u unExpr) Name() string {
	return u.name
}

func (u unExpr) Op() Op {
	return u.op
}

func (u unExpr) String() string {
	return fmt.Sprintf("%v %s", u.name, u.op)
}

func (u unExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(FilterExpressionWrapper{
		Un: &FilterUnExpression{
			Name: u.Name(),
			Op:   string(u.Op()),
		},
	})
}

func UnaryExpr(name string, op Op) FilterExpression {
	return &unExpr{
		name: name,
		op:   op,
	}
}

type betweenExpr struct {
	name        string
	left, right interface{}
}

func (b betweenExpr) Name() string {
	return b.name
}

func (b betweenExpr) Left() interface{} {
	return b.left
}

func (b betweenExpr) Right() interface{} {
	return b.right
}

func (b betweenExpr) String() string {
	return fmt.Sprintf("%v BETWEEN %v AND %v", b.name, wrapStr(b.left, "'"), wrapStr(b.right, "'"))
}

func (b betweenExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(FilterExpressionWrapper{
		Between: &FilterBetweenExpression{
			Left:  b.Left(),
			Name:  b.Name(),
			Right: b.Right(),
		},
	})
}
func BetweenExpr(name string, left interface{}, right interface{}) FilterExpression {
	return &betweenExpr{name: name, left: left, right: right}
}

type FilterBetweenExpression struct {
	Left  interface{} `json:"left"`
	Name  string      `json:"name"`
	Right interface{} `json:"right"`
}

type FilterExpressionWrapper struct {
	Between  *FilterBetweenExpression  `json:"between,omitempty"`
	In       *FilterInExpression       `json:"in,omitempty"`
	Junction *FilterJunctionExpression `json:"junction,omitempty"`
	Not      *FilterNotExpression      `json:"not,omitempty"`
	Simple   *FilterSimpleExpression   `json:"simple,omitempty"`
	Un       *FilterUnExpression       `json:"un,omitempty"`
}

type FilterInExpression struct {
	Name string        `json:"name"`
	Val  []interface{} `json:"val"`
}

type FilterJunctionExpression struct {
	Op  string                    `json:"op"`
	Sub []FilterExpressionWrapper `json:"sub"`
}

type FilterNotExpression struct {
	Sub FilterExpressionWrapper `json:"sub"`
}

type FilterSimpleExpression struct {
	Name string      `json:"name"`
	Op   string      `json:"op"`
	Val  interface{} `json:"val"`
}

type FilterUnExpression struct {
	Name string `json:"name"`
	Op   string `json:"op"`
}

func (e *FilterSimpleExpression) AsFilterExpression() FilterExpression {
	return SimpleExpr(e.Name, Op(e.Op), e.Val)
}

func (e *FilterJunctionExpression) AsFilterExpression() FilterExpression {
	return Junction(Op(e.Op), lo.Map(e.Sub, func(item FilterExpressionWrapper, _ int) FilterExpression {
		return item.AsFilterExpression()
	})...)
}

func (e *FilterNotExpression) AsFilterExpression() FilterExpression {
	return Not(e.Sub.AsFilterExpression())
}

func (e *FilterInExpression) AsFilterExpression() FilterExpression {
	return In(e.Name, e.Val)
}

func (e *FilterBetweenExpression) AsFilterExpression() FilterExpression {
	return BetweenExpr(e.Name, e.Left, e.Right)
}

func (e *FilterUnExpression) AsFilterExpression() FilterExpression {
	return UnaryExpr(e.Name, Op(e.Op))
}

func (few *FilterExpressionWrapper) AsFilterExpression() FilterExpression {
	if few.Simple != nil {
		return few.Simple.AsFilterExpression()
	}
	if few.Junction != nil {
		return few.Junction.AsFilterExpression()
	}
	if few.Not != nil {
		return few.Not.AsFilterExpression()
	}
	if few.In != nil {
		return few.In.AsFilterExpression()
	}
	if few.Between != nil {
		return few.Between.AsFilterExpression()
	}
	if few.Un != nil {
		return few.Un.AsFilterExpression()
	}
	return nil
}
