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
	"io"
	"net/http"

	"github.com/rkosegi/db2rest-bridge/pkg/crud"
	"github.com/rkosegi/db2rest-bridge/pkg/types"
)

type ItemHandler func(c crud.Interface, entity, id string, writer http.ResponseWriter, request *http.Request)
type EntityHandler func(c crud.Interface, entity string, writer http.ResponseWriter, request *http.Request)
type BackendHandler func(c crud.Interface, writer http.ResponseWriter, request *http.Request)

type Interface interface {
	io.Closer
	Run() error
}

func New(cfg *types.Config) Interface {
	return &restServer{
		cfg:    cfg,
		logger: configureLogging(cfg.LoggingConfig),
	}
}
