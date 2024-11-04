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
	"context"
	"net/http"
	"testing"

	. "github.com/jarcoal/httpmock"
	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

type mockType struct {
	Name string
	Age  int
}

type mockDoer struct{}

func (m *mockDoer) Do(req *http.Request) (*http.Response, error) {
	return DefaultTransport.RoundTrip(req)
}

func TestOpList(t *testing.T) {
	var (
		err error
		cl  GenericInterface[mockType]
		res []*mockType
	)
	Activate()
	defer DeactivateAndReset()
	RegisterResponder("GET", "http://loopback/dummy/mock",
		NewJsonResponderOrPanic(http.StatusOK, api.PagedResult{
			Data: lo.ToPtr([]api.UntypedDto{
				map[string]interface{}{
					"name": "Alice",
					"age":  42,
				},
			})}))

	cl, err = New[mockType]("http://loopback", "dummy", "mock",
		WithClientOptions[mockType](api.WithHTTPClient(&mockDoer{})))
	assert.NoError(t, err)
	assert.NotNil(t, cl)
	res, err = cl.List(context.Background(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, "Alice", res[0].Name)
	assert.Equal(t, 42, res[0].Age)
}

func TestOpCreate(t *testing.T) {
	var (
		err error
		cl  GenericInterface[mockType]
		res *mockType
	)
	dto := api.UntypedDto{
		"name": "Bob",
		"age":  43,
		"id":   1,
	}
	Activate()
	defer DeactivateAndReset()
	RegisterResponder("POST", "http://loopback/dummy/mock", NewJsonResponderOrPanic(http.StatusCreated, dto))

	cl, err = New[mockType]("http://loopback", "dummy", "mock",
		WithClientOptions[mockType](api.WithHTTPClient(&mockDoer{})))
	assert.NoError(t, err)
	assert.NotNil(t, cl)
	res, err = cl.Create(context.Background(), &mockType{Name: "Bob"})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Bob", res.Name)
	assert.Equal(t, 43, res.Age)
}
