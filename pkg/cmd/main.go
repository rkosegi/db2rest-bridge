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
	"errors"
	"flag"
	"github.com/rkosegi/db2rest-bridge/pkg/server"
	"github.com/rkosegi/db2rest-bridge/pkg/types"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
)

const (
	name = "db2rest_bridge"
)

func main() {
	var cfgFile string
	flag.StringVar(&cfgFile, "config", "config.yaml", "config file")
	flag.Parse()
	cfg, err := loadConfig(cfgFile)
	if err != nil {
		panic(err)
	}
	err = cfg.CheckAndNormalize()
	if err != nil {
		panic(err)
	}
	srv := server.New(cfg)
	if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func loadConfig(cfgFile string) (*types.Config, error) {
	var (
		cfg  types.Config
		err  error
		data []byte
	)

	if data, err = os.ReadFile(cfgFile); err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if err = cfg.CheckAndNormalize(); err != nil {
		return nil, err
	}
	return &cfg, nil
}
