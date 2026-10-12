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

package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/rkosegi/db2rest-bridge/pkg/client"
	"github.com/rkosegi/yaml-toolkit/fluent"
	"github.com/spf13/cobra"
)

type cmdData struct {
	profile    string
	configFile string
	lc         *api.Config
	entityName string
	backend    string

	// get, update, delete
	id string

	// list
	offset int
	size   int
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get single entity instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}

func backendsCmd(cd *cmdData) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backends",
		Short: "List backends",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(cd.profile) == 0 {
				return errors.New("profile name is required (--profile <name>)")
			}
			if len(cd.configFile) == 0 {
				return errors.New("config file name is required (--config <config file>)")
			}
			cd.lc = fluent.NewConfigHelper[api.Config]().Load(cd.configFile).Result()
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if cd.lc.ClientProfiles == nil {
				return errors.New("no client profiles found in the config file")
			}
			if cd.lc.ActiveClientProfile == nil {
				cd.lc.ActiveClientProfile = new("default")
			}
			cl, err := client.NewApi((*cd.lc.ClientProfiles)[*cd.lc.ActiveClientProfile].Url)
			if err != nil {
				return err
			}

			lbr, err := cl.ListBackendsWithResponse(cmd.Context())
			if err != nil {
				return err
			}

			if lbr.StatusCode() != http.StatusOK {
				return errors.New(lbr.Status())
			}

			for _, name := range *lbr.JSON200 {
				fmt.Println(name)
			}

			return nil
		},
	}
	return cmd
}

func listCmd() *cobra.Command {
	d := &cmdData{size: 20}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List entity instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.Flags().IntVar(&d.offset, "offset", d.offset, "Offset from begging of the view")
	cmd.Flags().IntVar(&d.size, "size", d.size, "How many instances to list. Actual value is also limited on server side")
	cmd.Flags().StringVar(&d.backend, "backend", "", "Name of backend to interact with")
	cmd.Flags().StringVar(&d.entityName, "entity", "", "Name of entity to act on")

	return cmd
}

func New() *cobra.Command {
	cd := &cmdData{
		profile:    "default",
		configFile: "config.yaml",
	}

	cmd := &cobra.Command{
		Use:   "client",
		Short: "Runs a client",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(cd.profile) == 0 {
				return errors.New("profile name is required (--profile <name>)")
			}
			if len(cd.configFile) == 0 {
				return errors.New("config file name is required (--config <config file>)")
			}
			cd.lc = fluent.NewConfigHelper[api.Config]().Load(cd.configFile).Result()
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&cd.profile, "profile", "p", cd.profile, "name of client profile within the configuration")
	cmd.PersistentFlags().StringVarP(&cd.configFile, "config", "c", cd.configFile, "Path to the config file")
	cmd.AddCommand(listCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(backendsCmd(cd))
	return cmd
}
