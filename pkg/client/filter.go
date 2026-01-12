/*
Copyright 2026 Richard Kosegi

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

package client

import (
	"github.com/google/go-cmp/cmp"
	"github.com/rkosegi/db2rest-bridge/pkg/query"
	"github.com/rkosegi/yaml-toolkit/dom"
)

func extractItemProperty[T any](propName string, item *T) interface{} {
	if node := dom.DecodeAnyToNode(item).AsContainer().Child(propName); node != nil && node.IsLeaf() {
		return node.AsLeaf().Value()
	}
	return nil
}

func InExpressionAsPredicate[T any](ie query.InExpression, item *T) bool {
	v := extractItemProperty(ie.Name(), item)
	for _, val := range ie.Values() {
		if cmp.Equal(val, v) {
			return true
		}
	}
	return false
}

func SimpleExpressionPredicate[T any](se query.SimpleExpression, item *T) bool {
	v := extractItemProperty(se.Name(), item)
	switch se.Op() {
	case query.OpEq:
		return cmp.Equal(se.Value(), v)
	case "!=", "<>":
		return !cmp.Equal(se.Value(), v)
	default:
		panic("unsupported op for binary expression: " + se.Op())
	}
}

func JunctionExpressionPredicate[T any](je query.JunctionExpression, item *T) bool {
	switch je.Op() {
	case query.OpAnd:
		out := true
		for _, se := range je.Sub() {
			out = out && FilterExpressionAsPredicate[T](se, item)
		}
		return out
	case query.OpOr:
		out := false
		for _, se := range je.Sub() {
			out = out || FilterExpressionAsPredicate[T](se, item)
		}
		return out
	}
	panic("unsupported op for junction: " + je.Op())
}

func NotExpressionAsPredicate[T any](ne query.NotExpression, item *T) bool {
	return !FilterExpressionAsPredicate[T](ne.Sub(), item)
}

func FilterExpressionAsPredicate[T any](fe query.FilterExpression, item *T) bool {
	if fe == nil {
		return true
	}
	switch fe.(type) {
	case query.InExpression:
		return InExpressionAsPredicate(fe.(query.InExpression), item)
	case query.SimpleExpression:
		return SimpleExpressionPredicate[T](fe.(query.SimpleExpression), item)
	case query.NotExpression:
		return NotExpressionAsPredicate(fe.(query.NotExpression), item)
	case query.JunctionExpression:
		return JunctionExpressionPredicate[T](fe.(query.JunctionExpression), item)
	default:
		panic("unsupported filter expression: " + fe.String())
	}
}

func ApplyPaging[T any](items []*T, paging query.Paging) []*T {
	offset := int(paging.Offset())
	if offset > len(items) {
		return nil
	}
	size := min(paging.Size(), offset+len(items))
	return items[offset : offset+size]
}
