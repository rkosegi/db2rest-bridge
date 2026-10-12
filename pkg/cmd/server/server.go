/*
Copyright 2026 Richard Kosegi

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
	"errors"
	"log/slog"
	"net/http"

	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/rkosegi/db2rest-bridge/pkg/server"
	xlog "github.com/rkosegi/slog-config"
	"github.com/rkosegi/yaml-toolkit/fluent"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cfgFile := "config.yaml"
	xc := xlog.MustNew("info", xlog.LogFormatLogFmt)
	var l *slog.Logger
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cfgFile == "" {
				return errors.New("must specify config file (--config)")
			}

			l = xc.Logger()
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cfg := fluent.NewConfigHelper[api.Config]().Load(cfgFile).Result()
			if err = cfg.CheckAndNormalize(); err != nil {
				return err
			}
			srv := server.New(cfg, server.WithLogger(l))
			defer func(srv server.Interface) {
				_ = srv.Close()
			}(srv)
			if err = srv.Run(cmd.Context()); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return err
			}
			return nil
		},
	}

	xc.AddPFlags(cmd.Flags())
	cmd.Flags().StringVar(&cfgFile, "config", cfgFile, "config file")

	return cmd
}
