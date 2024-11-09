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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	assert.Equal(t, " WHERE name = 'Bob' ORDER BY `name` DESC LIMIT 120, 10", NewBuilder().
		Paging(120, 10).
		OrderBy("name", false).
		Filter(SimpleExpr("name", OpEq, "Bob")).
		Build().
		String())

	assert.Equal(t, " WHERE name = 'Bob' LIMIT 0, 30", NewBuilder().
		Paging(0, 30).
		Filter(SimpleExpr("name", OpEq, "Bob")).
		Build().
		String())

	assert.Equal(t, " WHERE name = 'Bob' LIMIT 0, 20", NewBuilder().
		Filter(SimpleExpr("name", OpEq, "Bob")).
		Build().
		String())
}
