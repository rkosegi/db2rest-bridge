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
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/rkosegi/db2rest-bridge/pkg/crud"
	"github.com/rkosegi/db2rest-bridge/pkg/types"
)

type restServer struct {
	cfg     *types.Config
	server  *http.Server
	crudMap crud.NameToCrudMap
	logger  *slog.Logger
}

func (rs *restServer) Close() error {
	return rs.crudMap.Close()
}

func (rs *restServer) specHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if data, err := api.PathToRawSpec(r.URL.Path)[r.URL.Path](); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(data)
		}
	}
}

func (rs *restServer) Run() (err error) {
	rs.crudMap = make(crud.NameToCrudMap)
	for n, be := range rs.cfg.Backends {
		rs.logger.Debug("Opening backend", "name", n)
		if err = be.Open(); err != nil {
			rs.logger.Error("Unable to open backend", "name", n)
			return err
		}
		rs.crudMap[n] = crud.New(be, n, rs.logger)
	}
	rs.logger.Info("starting server", "listen address", rs.cfg.Server.HTTPListenAddress)

	cors := handlers.CORS(
		handlers.AllowedMethods([]string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		}),
		handlers.AllowedOrigins(rs.cfg.Server.CorsConfig.AllowedOrigins),
		handlers.MaxAge(rs.cfg.Server.CorsConfig.MaxAge),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	r := mux.NewRouter()
	r.HandleFunc("/spec/opeanapi.v1.json", rs.specHandler())

	rs.server = &http.Server{
		Addr: rs.cfg.Server.HTTPListenAddress,
		Handler: cors(api.HandlerWithOptions(rs, api.GorillaServerOptions{
			BaseURL:    "/api/v1",
			BaseRouter: r,
			Middlewares: []api.MiddlewareFunc{
				loggingMiddleware(rs.logger.With("type", "access log")),
			},
		})),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	if rs.cfg.Server.HTTPTLSConfig != nil {
		return rs.server.ListenAndServeTLS(
			rs.cfg.Server.HTTPTLSConfig.TLSCertPath,
			rs.cfg.Server.HTTPTLSConfig.TLSKeyPath,
		)
	} else {
		return rs.server.ListenAndServe()
	}
}
