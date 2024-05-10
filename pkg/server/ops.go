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
	"encoding/json"
	"net/http"
	"slices"

	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/rkosegi/db2rest-bridge/pkg/crud"
	"github.com/rkosegi/db2rest-bridge/pkg/query"
	"github.com/samber/lo"
)

func (rs *restServer) ListBackends(w http.ResponseWriter, _ *http.Request) {
	bes := lo.Keys(rs.cfg.Backends)
	slices.Sort(bes)
	sendJson(w, bes)
}

func (rs *restServer) ListEntities(w http.ResponseWriter, r *http.Request, backend string) {
	rs.handleBackend(w, r, backend, func(c crud.Interface, writer http.ResponseWriter, _ *http.Request) {
		if entities, err := c.ListEntities(); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			slices.Sort(entities)
			sendJson(writer, entities)
		}
	})
}

func (rs *restServer) ListItems(w http.ResponseWriter, r *http.Request, backend string, entity string, params api.ListItemsParams) {
	rs.handleEntity(w, r, backend, entity, func(c crud.Interface, entity string, writer http.ResponseWriter, request *http.Request) {
		var (
			err   error
			qry   query.Interface
			items []crud.Untyped
		)
		if qry, err = query.FromParams(params); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if items, err = c.ListItems(entity, qry); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		sendJson(writer, items)
	})
}

func (rs *restServer) CreateItem(w http.ResponseWriter, r *http.Request, backend string, entity string) {
	rs.handleEntity(w, r, backend, entity, func(c crud.Interface, entity string, writer http.ResponseWriter, request *http.Request) {
		var err error
		body := make(crud.Untyped)
		if err = json.NewDecoder(request.Body).Decode(&body); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		} else {
			if body, err = c.Create(entity, body); err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			} else {
				sendJson(writer, body)
			}
		}
	})
}

func (rs *restServer) GetItemById(w http.ResponseWriter, r *http.Request, backend string, entity string, id string) {
	rs.handleItem(w, r, backend, entity, id, func(c crud.Interface, entity, id string, writer http.ResponseWriter, _ *http.Request) {
		if obj, err := c.Get(entity, id); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		} else {
			if obj != nil {
				sendJson(writer, obj)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	})
}

func (rs *restServer) ExistsItemById(w http.ResponseWriter, r *http.Request, backend string, entity string, id string) {
	rs.handleItem(w, r, backend, entity, id, func(c crud.Interface, entity, id string, writer http.ResponseWriter, _ *http.Request) {
		if exists, err := c.Exists(entity, id); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			if exists {
				writer.WriteHeader(http.StatusNoContent)
			} else {
				writer.WriteHeader(http.StatusNotFound)
			}
		}
	})
}

func (rs *restServer) UpdateItemById(w http.ResponseWriter, r *http.Request, backend string, entity string, id string) {
	rs.handleItem(w, r, backend, entity, id, func(c crud.Interface, entity, id string, writer http.ResponseWriter, req *http.Request) {
		var (
			err    error
			exists bool
		)
		body := make(crud.Untyped)
		if err = json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if exists, err = c.Exists(entity, id); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if !exists {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		if body, err = c.Update(entity, id, body); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		sendJson(w, body)
	})
}

func (rs *restServer) DeleteItemById(w http.ResponseWriter, r *http.Request, backend string, entity string, id string) {
	rs.handleItem(w, r, backend, entity, id, func(c crud.Interface, entity, id string, writer http.ResponseWriter, _ *http.Request) {
		if err := c.Delete(entity, id); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			writer.WriteHeader(http.StatusNoContent)
		}
	})
}
