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
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/rkosegi/db2rest-bridge/pkg/crud"
	"github.com/rkosegi/db2rest-bridge/pkg/types"
)

var (
	httpDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "db2rest",
		Name:      "http_duration_seconds",
		Help:      "Duration of HTTP requests.",
	}, []string{"backend", "entity", "method", "status"})
	httpRequestBytes = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "db2rest",
		Name:      "http_request_bytes",
		Help:      "Total bytes of HTTP requests.",
	}, []string{"backend", "entity", "method", "status"})
	httpResponseBytes = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "db2rest",
		Name:      "http_response_bytes",
		Help:      "Total bytes of HTTP responses.",
	}, []string{"backend", "entity", "method", "status"})
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
	middlewares := []api.MiddlewareFunc{loggingMiddleware(rs.logger.With("type", "access log"))}

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

	if rs.cfg.Server.Telemetry.Enabled {
		middlewares = append(middlewares, func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				p := r.URL.Path
				if strings.HasPrefix(p, *rs.cfg.Server.APiPrefix) {
					p = strings.TrimPrefix(p, *rs.cfg.Server.APiPrefix)
					p = strings.TrimPrefix(p, "/")
					start := time.Now()
					ir := &interceptedResp{delegate: w}
					next.ServeHTTP(ir, r)
					if parts := strings.SplitN(p, "/", 3); len(parts) > 1 {
						httpDuration.WithLabelValues(parts[0], parts[1], r.Method, strconv.Itoa(ir.Status())).Observe(start.Sub(time.Now()).Seconds())
						if r.ContentLength > 0 {
							httpRequestBytes.WithLabelValues(parts[0], parts[1], r.Method, strconv.Itoa(ir.Status())).Add(float64(r.ContentLength))
						}
						httpResponseBytes.WithLabelValues(parts[0], parts[1], r.Method, strconv.Itoa(ir.Status())).Add(float64(ir.Written()))
					}
				} else {
					next.ServeHTTP(w, r)
				}
			})
		})
		r.Handle(rs.cfg.Server.Telemetry.Path, promhttp.Handler())
	}

	rs.server = &http.Server{
		Addr: rs.cfg.Server.HTTPListenAddress,
		Handler: cors(api.HandlerWithOptions(rs, api.GorillaServerOptions{
			BaseURL:     *rs.cfg.Server.APiPrefix,
			BaseRouter:  r,
			Middlewares: middlewares,
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
