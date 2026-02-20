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

package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterExpressionWrapper(t *testing.T) {
	var wrapper *FilterExpressionWrapper
	t.Run("simple", func(t *testing.T) {
		wrapper = &FilterExpressionWrapper{
			Simple: &FilterSimpleExpression{
				Name: "name",
				Op:   string(OpLike),
				Val:  "Alice%",
			},
		}
		fe := wrapper.AsFilterExpression()
		assert.Equal(t, "name LIKE 'Alice%'", fe.String())
	})
	t.Run("junction+simple", func(t *testing.T) {
		wrapper = &FilterExpressionWrapper{
			Junction: &FilterJunctionExpression{
				Op: string(OpAnd),
				Sub: []FilterExpressionWrapper{
					{
						Simple: &FilterSimpleExpression{
							Name: "age",
							Op:   string(OpGt),
							Val:  50,
						},
					},
					{
						Simple: &FilterSimpleExpression{
							Name: "salary",
							Op:   string(OpLt),
							Val:  100,
						},
					},
				},
			},
		}
		fe := wrapper.AsFilterExpression()
		assert.Equal(t, "((age > 50) AND (salary < 100))", fe.String())
	})
	t.Run("not+in", func(t *testing.T) {
		wrapper = &FilterExpressionWrapper{
			Not: &FilterNotExpression{
				Sub: FilterExpressionWrapper{
					In: &FilterInExpression{
						Name: "user_id",
						Val:  []interface{}{1, 2, 3},
					},
				},
			},
		}
		fe := wrapper.AsFilterExpression()
		assert.Equal(t, "NOT (user_id IN (1,2,3))", fe.String())

	})

	t.Run("not null+between", func(t *testing.T) {
		wrapper = &FilterExpressionWrapper{
			Junction: &FilterJunctionExpression{
				Op: string(OpAnd),
				Sub: []FilterExpressionWrapper{
					{
						Un: &FilterUnExpression{
							Name: "department",
							Op:   string(OpIsNotNull),
						},
					},
					{
						Between: &FilterBetweenExpression{
							Left:  70,
							Name:  "age",
							Right: 100,
						},
					},
				},
			},
		}
		fe := wrapper.AsFilterExpression()
		assert.Equal(t, "((department IS NOT NULL) AND (age BETWEEN 70 AND 100))", fe.String())
	})

	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, (&FilterExpressionWrapper{}).AsFilterExpression())
	})
}
