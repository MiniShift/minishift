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
	"fmt"
	"github.com/golang/glog"
	"os"
	"runtime"

	hostfolderActions "github.com/minishift/minishift/pkg/minishift/hostfolder"
	"github.com/minishift/minishift/pkg/util/os/atexit"
	"github.com/spf13/cobra"
)

var (
	instanceOnly bool
	usersShare   bool
)

var hostfolderAddCmd = &cobra.Command{
	Use:   "add HOSTFOLDER_NAME",
	Short: "Add a host folder definition",
	Long:  `Add a host folder definition that can be mounted to a running cluster`,
	Run: func(cmd *cobra.Command, args []string) {

		var err error = nil
		if usersShare && runtime.GOOS == "windows" {
			// Windows-only (CIFS), all instances
			err = hostfolderActions.SetupUsers(true)
		} else {
			if len(args) < 1 {
				fmt.Fprintln(os.Stderr, "usage: minishift hostfolder add HOSTFOLDER_NAME")
				atexit.Exit(1)
			}
			err = hostfolderActions.Add(args[0], !instanceOnly)
		}

		if err != nil {
			glog.Errorln(err)
			atexit.Exit(1)
		}
	},
}

func init() {
	HostfolderCmd.AddCommand(hostfolderAddCmd)
	hostfolderAddCmd.Flags().BoolVarP(&instanceOnly, "instance-only", "", false, "Define host folder only for the current cluster instance.")

	// Windows-only
	if runtime.GOOS == "windows" {
		hostfolderAddCmd.Flags().BoolVarP(&usersShare, "users-share", "", false, "Define host folder for the Users share on a Windows host.")
	}
}
