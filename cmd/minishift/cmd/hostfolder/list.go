/*
Copyright (C) 2017 Red Hat, Inc.

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

package hostfolder

import (
	"github.com/golang/glog"

	"github.com/docker/machine/libmachine"
	"github.com/minishift/minishift/pkg/minikube/constants"
	hostfolderActions "github.com/minishift/minishift/pkg/minishift/hostfolder"
	"github.com/minishift/minishift/pkg/util/os/atexit"
	"github.com/spf13/cobra"
)

var hostfolderListCmd = &cobra.Command{
	Use:   "list",
	Short: "List an overview of defined host folders",
	Long:  `List an overview of defined host folders that can be mounted to a running cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		api := libmachine.NewClient(constants.Minipath, constants.MakeMiniPath("certs"))
		defer api.Close()
		host, err := api.Load(constants.MachineName)
		if err != nil {
			glog.Errorln("Error: ", err)
			atexit.Exit(1)
		}

		isRunning := isHostRunning(host.Driver)
		err = hostfolderActions.List(host.Driver, isRunning)
		if err != nil {
			glog.Errorln("Error: ", err)
			atexit.Exit(1)
		}
	},
}

func init() {
	HostfolderCmd.AddCommand(hostfolderListCmd)
}
