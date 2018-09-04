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

package openshift

import (
	"github.com/spf13/cobra"

	"github.com/docker/machine/libmachine"

	"fmt"

	"github.com/docker/machine/libmachine/provision"
	"github.com/minishift/minishift/cmd/minishift/state"
	"github.com/minishift/minishift/pkg/minikube/cluster"
	"github.com/minishift/minishift/pkg/minishift/docker"
	"github.com/minishift/minishift/pkg/minishift/openshift"
	"github.com/minishift/minishift/pkg/util/os/atexit"
)

const (
	configTargetFlag         = "target"
	unknownConfigTargetError = "Unknown configuration target. Only 'master', 'node' and 'kube' are supported."
)

var (
	configTarget string
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Displays the specified OpenShift configuration resource.",
	Long:  "Displays the specified OpenShift configuration resource.",
	Run:   runViewConfig,
}

func init() {
	viewCmd.Flags().StringVar(&configTarget, configTargetFlag, "master", "Target configuration to display. Options are 'master', 'node' and 'kube'.")
	configCmd.AddCommand(viewCmd)
}

func runViewConfig(cmd *cobra.Command, args []string) {
	configFileTarget := determineTarget(configTarget)
	if configFileTarget == openshift.GetOpenShiftPatchTarget("unknown") {
		atexit.ExitWithMessage(1, unknownConfigTargetError)
	}

	api := libmachine.NewClient(state.InstanceDirs.Home, state.InstanceDirs.Certs)
	defer api.Close()

	host, err := cluster.CheckIfApiExistsAndLoad(api)
	if err != nil {
		atexit.ExitWithMessage(1, nonExistentMachineError)
	}

	sshCommander := provision.GenericSSHCommander{Driver: host.Driver}
	dockerCommander := docker.NewVmDockerCommander(sshCommander)

	out, err := openshift.ViewConfig(configFileTarget, dockerCommander)
	if err != nil {
		atexit.ExitWithMessage(1, fmt.Sprintf("Cannot display the OpenShift configuration: %s", err.Error()))
	}

	fmt.Println(out)
}
