/*
Copyright (C) 2018 Red Hat, Inc.

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
	"fmt"
	"os"

	minishiftConstants "github.com/minishift/minishift/pkg/minishift/constants"
	"github.com/spf13/cobra"
)

// version command represent current running openshift version and available one.
var componentListCmd = &cobra.Command{
	Use:   "list [component-name]",
	Short: "List valid components that can be added to an OpenShift cluster (Works only with OpenShift version >= 3.10.x)",
	Long:  "List valid components that can be added to an OpenShift cluster (Works only with OpenShift version >= 3.10.x)",
	Run:   runComponentList,
}

func runComponentList(cmd *cobra.Command, args []string) {
	fmt.Fprint(os.Stdout, "The following OpenShift components are available: \n")
	for _, component := range minishiftConstants.ValidComponents {
		fmt.Fprintf(os.Stdout, "\t- %s\n", component)
	}
}

func init() {
	componentCmd.AddCommand(componentListCmd)
}
