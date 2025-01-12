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

package crud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testBody = map[string]interface{}{
		"name": "my-name",
		"age":  30,
	}
	testEnt = "myentity"
	testId  = "id"
)

func TestCreateUpdateQuery(t *testing.T) {
	sql, vals := createUpdateQuery(testEnt, "ent_id", testBody)
	assert.Equal(t, "UPDATE `myentity` SET `age` = ?, `name` = ? WHERE `ent_id` = ? LIMIT 1", sql)
	assert.NotNil(t, vals)
	assert.Equal(t, 30, vals[0])
	assert.Equal(t, "my-name", vals[1])
}

func TestCreateInsertQuery(t *testing.T) {
	sql, vals := createInsertQuery(testEnt, testBody)
	assert.Equal(t, "INSERT INTO `myentity` (`age`,`name`) VALUES(?,?)", sql)
	assert.NotNil(t, vals)
}

func TestCreateDeleteQuery(t *testing.T) {
	sql := createSingleDeleteQuery(testEnt, testId)
	assert.Equal(t, "DELETE FROM `myentity` WHERE `id` = ? LIMIT 1", sql)
}

func TestCreateSingleSelectQuery(t *testing.T) {
	sql := createSingleSelectQuery(testEnt, testId)
	assert.Equal(t, "SELECT * FROM `myentity` WHERE `id` = ? LIMIT 1", sql)
}
