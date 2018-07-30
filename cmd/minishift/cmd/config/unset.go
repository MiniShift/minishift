/*
Copyright (C) 2016 Red Hat, Inc.

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

package config

import (
	"fmt"
	"os"

	"github.com/minishift/minishift/pkg/minikube/constants"
	"github.com/minishift/minishift/pkg/minishift/config"
	"github.com/minishift/minishift/pkg/util/os/atexit"
	"github.com/spf13/cobra"
)

var configUnsetCmd = &cobra.Command{
	Use:   "unset PROPERTY_NAME",
	Short: "Clears the value of a configuration property in the Minishift configuration file.",
	Long:  "Clears the value of a configuration property in the Minishift configuration file. The value can be overwritten at runtime by flags or environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			atexit.ExitWithMessage(1, "usage: minishift config unset PROPERTY_NAME")
		}
		err := unset(args[0])
		if err != nil {
			fmt.Fprintln(os.Stdout, err)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(configUnsetCmd)
	configUnsetCmd.Flags().BoolVar(&global, "global", false, "Unset the value of a configuration property in the global configuration file.")
}

func unset(name string) error {
	_, err := findSetting(name)
	if err != nil {
		return err
	}
	confFile := constants.ConfigFile
	if global {
		confFile = constants.GlobalConfigFile
	}
	m, err := config.ReadViperConfig(confFile)
	if err != nil {
		return err
	}
	if m[name] == nil {
		return fmt.Errorf("Property name '%s' is not set", name)
	}
	delete(m, name)
	fmt.Printf("Property name '%s' successfully unset\n", name)
	return config.WriteViperConfig(confFile, m)
}
