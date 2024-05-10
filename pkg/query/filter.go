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
	"github.com/samber/lo"
	"strings"
)

type simpleExpr struct {
	name string
	op   Op
	val  interface{}
}

func (s simpleExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"simple": map[string]interface{}{
			"name": s.name,
			"op":   s.op,
			"val":  s.val,
		},
	})
}

func wrapStr(v interface{}, q string) interface{} {
	if _, ok := v.(string); ok {
		return q + v.(string) + q
	}
	return v
}

func (s simpleExpr) String() string {
	return fmt.Sprintf(`%s %s %v`, s.name, s.op, wrapStr(s.val, "'"))
}

type junctionExpr struct {
	op  Op
	sub []FilterExpression
}

func (j junctionExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"junction": map[string]interface{}{
			"sub": j.sub,
			"op":  j.op,
		},
	})
}

func (j junctionExpr) String() string {
	var sb strings.Builder
	sb.WriteRune('(')
	sb.WriteString(strings.Join(lo.Map(j.sub, func(item FilterExpression, _ int) string {
		return "(" + item.String() + ")"
	}), " "+string(j.op)+" "))
	sb.WriteRune(')')

	return sb.String()
}

type notExpr struct {
	expr FilterExpression
}

func (n notExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"not": n.expr,
	})
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
