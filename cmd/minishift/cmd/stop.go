/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package cmd

import (
	"fmt"

	"github.com/docker/machine/libmachine"
	registrationUtil "github.com/minishift/minishift/cmd/minishift/cmd/registration"
	"github.com/minishift/minishift/cmd/minishift/cmd/util"
	"github.com/minishift/minishift/cmd/minishift/state"
	"github.com/minishift/minishift/pkg/minikube/cluster"
	"github.com/minishift/minishift/pkg/minikube/constants"
	"github.com/minishift/minishift/pkg/util/os/atexit"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the running local OpenShift cluster.",
	Long: `Stops the running local OpenShift cluster. This command stops the Minishift
VM but does not delete any associated files. To start the cluster again, use the 'minishift start' command.`,
	Run: runStop,
}

func runStop(cmd *cobra.Command, args []string) {
	api := libmachine.NewClient(state.InstanceDirs.Home, state.InstanceDirs.Certs)
	defer api.Close()

	// if VM does not exist, exit with error
	util.ExitIfUndefined(api, constants.MachineName)

	hostVm, err := api.Load(constants.MachineName)
	if err != nil {
		atexit.ExitWithMessage(1, err.Error())
	}

	// check if VM is already in stopped state
	if util.IsHostStopped(hostVm.Driver) {
		atexit.ExitWithMessage(0, fmt.Sprintf("The '%s' VM is already stopped.", constants.MachineName))
	}

	fmt.Println("Stopping the OpenShift cluster...")

	if hostVm.Driver.DriverName() == "generic" {
		if err := util.OcClusterDown(hostVm); err != nil {
			atexit.ExitWithMessage(1, err.Error())
		}
	} else {
		// Unregister, allow to be skipped and force deletion is ignored
		registrationUtil.UnregisterHost(api, true, false)

		if err := cluster.StopHost(api); err != nil {
			atexit.ExitWithMessage(1, fmt.Sprintf("Error stopping cluster: %s", err.Error()))
		}
	}
	fmt.Println("Cluster stopped.")
}

func init() {
	stopCmd.Flags().BoolVar(&registrationUtil.SkipUnRegistration, "skip-unregistration", false, "Skip the virtual machine unregistration.")
	RootCmd.AddCommand(stopCmd)
}
