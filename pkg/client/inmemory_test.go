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
	"testing"

	"github.com/rkosegi/db2rest-bridge/pkg/query"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

type employee struct {
	Name       string
	Department string
	Company    string
	Age        int
	Salary     int
	ExtID      int
}

var testdata = []*employee{
	{
		Name:       "Alice",
		Age:        50,
		Salary:     1000,
		ExtID:      14,
		Company:    "acme ltd.",
		Department: "HR",
	},
	{
		Name:       "Bob",
		Age:        42,
		Salary:     500,
		ExtID:      12,
		Company:    "acme ltd.",
		Department: "Engineering",
	},
	{
		Name:       "Charlie",
		Age:        48,
		Salary:     3500,
		ExtID:      86,
		Company:    "acme ltd.",
		Department: "Management",
	},
}

func TestInMemoryCrud(t *testing.T) {
	var (
		items []*employee
		total int
		emp   *employee
		err   error
	)
	ic := NewInMemory[employee](func(id string, em *employee) bool {
		return em.Name == id
	}, testdata)

	t.Run("get all items", func(t *testing.T) {
		items, total, err = ic.List(t.Context(), query.NewBuilder().Build())
		assert.NoError(t, err)
		assert.Equal(t, 3, total)
		assert.Len(t, items, 3)
	})

	t.Run("get all items using LoadAll", func(t *testing.T) {
		items, err = LoadAll(t.Context(), ic, nil)
		assert.NoError(t, err)
		assert.Len(t, items, 3)
	})

	t.Run("create/get/remove item", func(t *testing.T) {
		items, total, err = ic.List(t.Context(), query.NewBuilder().Build())

		emp, err = ic.Create(t.Context(), &employee{
			Name:    "John",
			Company: "Umbrella corp",
		})
		assert.NoError(t, err)
		assert.Equal(t, "John", emp.Name)
		total2 := total + 1

		// after create
		items, _, err = ic.List(t.Context(), query.NewBuilder().Build())
		assert.Len(t, items, total2)

		items, total, err = ic.List(t.Context(), query.NewBuilder().Filter(
			query.SimpleExpr("Name", query.OpEq, "John")).Build(),
		)

		assert.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Equal(t, "John", items[0].Name)

		assert.NoError(t, ic.Delete(t.Context(), "John"))
		total2 = total2 - 1

		items, total, err = ic.List(t.Context(), query.NewBuilder().Build())
		assert.Len(t, items, total2)
	})

	t.Run("get non-existent", func(t *testing.T) {
		emp, err = ic.Get(t.Context(), "Nobody")
		assert.NoError(t, err)
		assert.Nil(t, emp)
	})

	t.Run("filter - in expression", func(t *testing.T) {
		items, total, err = ic.List(t.Context(), query.NewBuilder().
			Filter(query.In("ExtID", []interface{}{12, 14, 16, 18})).
			Build())
		assert.NoError(t, err)
		assert.Equal(t, 2, total)
	})

	t.Run("filter - not expression", func(t *testing.T) {
		items, total, err = ic.List(t.Context(), query.NewBuilder().
			Filter(query.Not(query.In("ExtID", []interface{}{12, 14, 16, 18}))).
			Build())
		assert.NoError(t, err)
		assert.Equal(t, 1, total)
	})

	t.Run("filter - junction expression", func(t *testing.T) {
		items, total, err = ic.List(t.Context(), query.NewBuilder().
			Filter(query.Junction(
				query.OpOr,
				query.SimpleExpr("Salary", query.OpEq, 500),
				query.SimpleExpr("Age", query.OpEq, 48),
			)).
			Build())
		assert.NoError(t, err)
		assert.Equal(t, 2, total)

		items, total, err = ic.List(t.Context(), query.NewBuilder().
			Filter(query.Junction(
				query.OpAnd,
				query.SimpleExpr("Salary", "!=", 500),
				query.SimpleExpr("Age", "!=", 50),
			)).
			Build())
		assert.NoError(t, err)
		assert.Equal(t, 1, total)

	})

	t.Run("get/update item", func(t *testing.T) {
		emp, err = ic.Get(t.Context(), "Alice")
		assert.NoError(t, err)
		assert.Equal(t, "Alice", emp.Name)
		assert.Equal(t, 14, emp.ExtID)

		emp, err = ic.Update(t.Context(), "Alice", &employee{ExtID: 99})
		assert.NoError(t, err)
		assert.Equal(t, 99, emp.ExtID)
	})
}

func TestApplyPaging(t *testing.T) {
	var out []*int
	arr := []*int{lo.ToPtr(0), lo.ToPtr(1), lo.ToPtr(2), lo.ToPtr(3), lo.ToPtr(4), lo.ToPtr(5),
		lo.ToPtr(6), lo.ToPtr(7), lo.ToPtr(8), lo.ToPtr(9), lo.ToPtr(10), lo.ToPtr(1)}

	out = ApplyPaging[int](arr, query.Page(5, 3))
	assert.Len(t, out, 3)
	assert.Equal(t, 5, *out[0])
	assert.Equal(t, 6, *out[1])
	assert.Equal(t, 7, *out[2])

	out = ApplyPaging[int](arr, query.Page(100, 10))
	assert.Len(t, out, 0)
}

func TestExtractProperty(t *testing.T) {
	emp := &employee{
		Name: "Alice",
	}
	assert.Equal(t, "Alice", extractItemProperty[employee]("Name", emp))
	assert.Nil(t, extractItemProperty[employee]("non sense", emp))
}
