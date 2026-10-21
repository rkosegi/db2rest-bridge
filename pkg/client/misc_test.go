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

package client

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/rkosegi/db2rest-bridge/pkg/query"
	"github.com/stretchr/testify/assert"
)

func TestExcludeProps(t *testing.T) {
	in := map[string]any{
		"created": "2024-11-01 22:32:21",
		"name":    "Hello",
		"age":     42,
	}
	out := excludeProps(in, []string{"created"})
	assert.Len(t, out, 2)
	assert.Equal(t, "Hello", out["name"])
	assert.Equal(t, 42, out["age"])
}

func TestOnlyProps(t *testing.T) {
	out := onlyProps(map[string]any{
		"id":   999,
		"name": "Alice",
		"age":  30,
	}, []string{"id"})
	assert.Len(t, out, 1)
	assert.Equal(t, 999, out["id"])
}

func TestEnsureResponseCode(t *testing.T) {
	assert.Error(t, ensureResponseCode(&http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader("Not found")),
	}, http.StatusOK, nil))

	assert.NoError(t, ensureResponseCode(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("Not found")),
	}, http.StatusOK, nil))
}

func TestLoad(t *testing.T) {
	c := NewInMemory[string](func(s string, a *string) bool {
		return true
	}, func() []*string {
		var out []*string
		for i := range 35 {
			out = append(out, new(fmt.Sprintf("%d", i)))
		}
		return out
	}())
	t.Run("All", func(t *testing.T) {
		res, err := LoadAll[string](t.Context(), c, nil, query.OrderBy("id", true))
		assert.NoError(t, err)
		assert.Len(t, res, 35)
	})
	t.Run("AllInBatch", func(t *testing.T) {
		res, err := LoadAllInBatch[string](t.Context(), c, nil, 20, query.OrderBy("id", false))
		assert.NoError(t, err)
		assert.Len(t, res, 35)
	})
	t.Run("AllInBatchWithCap", func(t *testing.T) {
		res, err := LoadAllInBatchWithCap[string](t.Context(), c, nil, 10, 33)
		assert.NoError(t, err)
		assert.Len(t, res, 33)
	})
}
