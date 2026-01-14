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
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExcludeProps(t *testing.T) {
	in := map[string]interface{}{
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
	out := onlyProps(map[string]interface{}{
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
