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

package main

import (
	"context"
	"errors"
	"flag"
	"net/http"

	"github.com/rkosegi/db2rest-bridge/pkg/server"
	"github.com/rkosegi/db2rest-bridge/pkg/types"
	"github.com/rkosegi/yaml-toolkit/fluent"
)

const (
	name = "db2rest_bridge"
)

func main() {
	var (
		cfgFile string
		err     error
	)
	flag.StringVar(&cfgFile, "config", "config.yaml", "config file")
	flag.Parse()

	cfg := fluent.NewConfigHelper[types.Config]().Load(cfgFile).Result()
	if err = cfg.CheckAndNormalize(); err != nil {
		panic(err)
	}
	srv := server.New(cfg)
	if err = srv.Run(context.Background()); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
