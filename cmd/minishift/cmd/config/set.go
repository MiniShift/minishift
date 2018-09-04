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
	"github.com/minishift/minishift/pkg/minikube/constants"
	"github.com/minishift/minishift/pkg/minishift/config"
	"github.com/minishift/minishift/pkg/util/os/atexit"
	"github.com/spf13/cobra"
)

var configSetCmd = &cobra.Command{
	Use:   "set PROPERTY_NAME PROPERTY_VALUE",
	Short: "Sets the value of a configuration property in the Minishift configuration file.",
	Long: `Sets the value of one or more configuration properties in the Minishift configuration file.
These values can be overwritten by flags or environment variables at runtime.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			atexit.ExitWithMessage(1, "usage: minishift config set PROPERTY_NAME PROPERTY_VALUE")
		}
		err := Set(args[0], args[1], true)
		if err != nil {
			atexit.ExitWithMessage(1, err.Error())
		}
	},
}

func init() {
	ConfigCmd.AddCommand(configSetCmd)
	configSetCmd.Flags().BoolVar(&global, "global", false, "Sets the value of a configuration property in the global configuration file.")
}

func Set(name string, value string, runCallback bool) error {
	s, err := findSetting(name)
	if err != nil {
		return err
	}
	// Validate the new value
	err = run(name, value, s.validations)
	if err != nil {
		return err
	}

	// Set the value
	confFile := constants.ConfigFile
	if global {
		confFile = constants.GlobalConfigFile
	}
	conf, err := config.ReadViperConfig(confFile)
	if err != nil {
		return err
	}
	err = s.set(conf, name, value)
	if err != nil {
		return err
	}

	if runCallback {
		// Run any callbacks for this property
		err = run(name, value, s.callbacks)
		if err != nil {
			return err
		}
	}

	// Write the value
	return config.WriteViperConfig(confFile, conf)
}
