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
	"os"

	"github.com/docker/machine/libmachine"
	"github.com/minishift/minishift/cmd/minishift/state"
	"github.com/minishift/minishift/pkg/minikube/cluster"
	"github.com/minishift/minishift/pkg/util/os/atexit"
	"github.com/spf13/cobra"
)

var (
	follow bool
	tail   int64
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Gets the logs of the running OpenShift cluster.",
	Long:  `Gets the logs of the running OpenShift cluster. The logs do not contain information about your application code.`,
	Run: func(cmd *cobra.Command, args []string) {
		api := libmachine.NewClient(state.InstanceDirs.Home, state.InstanceDirs.Certs)
		defer api.Close()
		s, err := cluster.GetHostLogs(api, follow, tail)
		if err != nil {
			atexit.ExitWithMessage(1, fmt.Sprintf("Error getting logs: %s", err.Error()))
		}
		fmt.Fprintln(os.Stdout, s)
	},
}

func init() {
	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Continuously print the logs entries")
	logsCmd.Flags().Int64VarP(&tail, "tail", "t", -1, "Number of lines to show from the end of the logs")
	RootCmd.AddCommand(logsCmd)
}
