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

	"encoding/json"
	"fmt"

	"github.com/docker/machine/libmachine/provision"
	"github.com/minishift/minishift/cmd/minishift/state"
	"github.com/minishift/minishift/pkg/minikube/cluster"
	"github.com/minishift/minishift/pkg/minishift/docker"
	"github.com/minishift/minishift/pkg/minishift/openshift"
	"github.com/minishift/minishift/pkg/util/os/atexit"
)

const (
	targetFlag = "target"
	patchFlag  = "patch"

	unknownPatchTargetError = "Unkown patch target. Only 'master', 'node' and 'kube' are supported."
	emptyPatchError         = "You must specify a patch using the --patch flag."
	invalidJSONError        = "The patch must be a valid JSON file."
)

var (
	target string
	patch  string
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Patches the OpenShift configuration resource with the specified patch.",
	Long:  "Patches the OpenShift configuration resource with the specified patch. The patch must be a valid JSON file.",
	Run:   runPatch,
}

func init() {
	setCmd.Flags().StringVar(&target, targetFlag, "master", "Target configuration to patch. Options are 'master', 'node' and 'kube'.")
	setCmd.Flags().StringVar(&patch, patchFlag, "", "The patch to apply.")
	configCmd.AddCommand(setCmd)
}

func runPatch(cmd *cobra.Command, args []string) {
	patchTarget := determineTarget(target)
	if patchTarget == openshift.GetOpenShiftPatchTarget("unknown") {
		atexit.ExitWithMessage(1, unknownPatchTargetError)
	}

	validatePatch(patch)

	api := libmachine.NewClient(state.InstanceDirs.Home, state.InstanceDirs.Certs)
	defer api.Close()

	host, err := cluster.CheckIfApiExistsAndLoad(api)
	if err != nil {
		atexit.ExitWithMessage(1, nonExistentMachineError)
	}

	sshCommander := provision.GenericSSHCommander{Driver: host.Driver}
	dockerCommander := docker.NewVmDockerCommander(sshCommander)

	_, err = openshift.Patch(patchTarget, patch, dockerCommander)
	if err != nil {
		atexit.ExitWithMessage(1, fmt.Sprintf("Error patching the OpenShift configuration: %s", err.Error()))
	}
}

func determineTarget(target string) openshift.OpenShiftPatchTarget {
	switch target {
	case "master":
		return openshift.GetOpenShiftPatchTarget("master")
	case "node":
		return openshift.GetOpenShiftPatchTarget("node")
	case "kube":
		return openshift.GetOpenShiftPatchTarget("kube")
	default:
		return openshift.GetOpenShiftPatchTarget("unknown")
	}
}

func validatePatch(patch string) {
	if len(patch) == 0 {
		atexit.ExitWithMessage(1, emptyPatchError)
	}

	if !isJSON(patch) {
		atexit.ExitWithMessage(1, invalidJSONError)
	}

}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}
