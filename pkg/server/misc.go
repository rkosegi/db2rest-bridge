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

package server

import (
	"fmt"
	"net/http"

	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/rkosegi/db2rest-bridge/pkg/crud"
)

func (rs *restServer) handleBackend(writer http.ResponseWriter, request *http.Request, backend string, handler BackendHandler) {
	if c, ok := rs.crudMap[backend]; !ok {
		http.Error(writer, fmt.Sprintf("no such backend: %s", backend), http.StatusBadRequest)
		return
	} else {
		handler(c, writer, request)
	}
}

func (rs *restServer) handleEntity(writer http.ResponseWriter, request *http.Request, backend, entity string, handler EntityHandler) {
	rs.handleBackend(writer, request, backend, func(c crud.Interface, writer http.ResponseWriter, request *http.Request) {
		handler(c, entity, writer, request)
	})
}

func (rs *restServer) handleItem(writer http.ResponseWriter, request *http.Request, backend, entity, item string, handler ItemHandler) {
	rs.handleEntity(writer, request, backend, entity, func(c crud.Interface, entity string, writer http.ResponseWriter, request *http.Request) {
		handler(c, entity, item, writer, request)
	})
}

func extractIds(objs []api.UntypedDto, idCol string) ([]interface{}, error) {
	var ids []interface{}
	for _, obj := range objs {
		if id, ok := obj[idCol]; ok {
			ids = append(ids, id)
		} else {
			return nil, fmt.Errorf("can't extract ID column (%s) value from object %v", idCol, obj)
		}
	}
	return ids, nil
}
