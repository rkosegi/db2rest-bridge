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

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeConfig(t *testing.T) {
	var c *Config
	t.Run("page size limit from be", func(t *testing.T) {
		c = &Config{Backends: Backends{
			"A": &BackendConfig{
				PageSizeLimit: new(7),
			},
		}}
		assert.NoError(t, c.CheckAndNormalize())
		assert.Equal(t, 7, c.Backends["A"].ApplyPageSizeLimit("x", 100))
	})
	t.Run("page size limit from global", func(t *testing.T) {
		c = &Config{Backends: Backends{
			"A": &BackendConfig{},
		}}
		assert.NoError(t, c.CheckAndNormalize())
		assert.Equal(t, 20, c.Backends["A"].ApplyPageSizeLimit("x", 200))
	})
	t.Run("page size limit from caller", func(t *testing.T) {
		c = &Config{Backends: Backends{
			"A": &BackendConfig{
				Entities: new(map[string]*EntityConfig{
					"x": &EntityConfig{
						PageSizeLimit: new(7),
					},
				}),
			},
		}}
		assert.NoError(t, c.CheckAndNormalize())
		assert.Equal(t, 7, c.Backends["A"].ApplyPageSizeLimit("x", 5))
	})

}
