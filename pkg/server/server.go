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
	"github.com/rkosegi/go-http-commons/middlewares"
	"github.com/rkosegi/go-http-commons/openapi"
	"github.com/rkosegi/go-http-commons/output"
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

	out = output.NewBuilder().Build()
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
	rs.logger.Info("starting server", "listen address", rs.cfg.Server.ListenAddress)
	mws := []api.MiddlewareFunc{middlewares.NewLoggingBuilder().WithLogger(rs.logger).Build()}

	cors := handlers.CORS(
		handlers.AllowedMethods([]string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		}),
		handlers.AllowedOrigins(rs.cfg.Server.Cors.AllowedOrigins),
		handlers.MaxAge(rs.cfg.Server.Cors.MaxAge),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	r := mux.NewRouter()
	r.HandleFunc("/spec/opeanapi.v1.json", openapi.SpecHandler(api.PathToRawSpec))

	if rs.cfg.Server.Telemetry.Enabled {
		mws = append(mws, middlewares.NewInterceptorBuilder().
			WithRequestFilter(func(r *http.Request) bool {
				return strings.HasPrefix(r.URL.Path, *rs.cfg.Server.APIPrefix)
			}).
			WithCallback(func(resp middlewares.InterceptedResponse) {
				req := resp.Request()
				p := req.URL.Path
				p = strings.TrimPrefix(p, *rs.cfg.Server.APIPrefix)
				p = strings.TrimPrefix(p, "/")
				start := time.Now()
				if parts := strings.SplitN(p, "/", 3); len(parts) > 1 {
					httpDuration.WithLabelValues(parts[0], parts[1], req.Method,
						strconv.Itoa(resp.Status())).Observe(start.Sub(time.Now()).Seconds())
					if req.ContentLength > 0 {
						httpRequestBytes.WithLabelValues(parts[0], parts[1], req.Method,
							strconv.Itoa(resp.Status())).Add(float64(req.ContentLength))
					}
					httpResponseBytes.WithLabelValues(parts[0], parts[1], req.Method,
						strconv.Itoa(resp.Status())).Add(float64(resp.Written()))
				}
			}).
			Build())
		r.Handle(*rs.cfg.Server.Telemetry.Path, promhttp.Handler())
	}

	rs.server = &http.Server{
		Addr: rs.cfg.Server.ListenAddress,
		Handler: cors(api.HandlerWithOptions(rs, api.GorillaServerOptions{
			BaseURL:     *rs.cfg.Server.APIPrefix,
			BaseRouter:  r,
			Middlewares: mws,
		})),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return rs.cfg.Server.RunForever(rs.server)
}
