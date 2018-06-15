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

package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/minishift/minishift/pkg/minishift/constants"
)

type OpenShiftCommand struct {
	*defaultCommand
}

func NewOpenShiftCommand(command string, ignoreError bool, outputVariable string) *OpenShiftCommand {
	defaultCommand := &defaultCommand{rawCommand: command, ignoreError: ignoreError, outputVariable: outputVariable}
	openShiftCommand := &OpenShiftCommand{defaultCommand}
	defaultCommand.fn = openShiftCommand.doExecute
	return openShiftCommand
}

func (c *OpenShiftCommand) doExecute(ec *ExecutionContext, ignoreError bool, outputVariable string) error {
	// split off the actual 'openshift' command. We are using origin container to run those commands
	cmd := strings.Replace(c.rawCommand, "openshift ", "", 1)
	cmd = ec.Interpolate(cmd)
	fmt.Print(".")

	commander := ec.GetDockerCommander()
	output, err := commander.Exec("-t", constants.OpenshiftContainerName, constants.OpenshiftOcExec, ec.Interpolate(cmd))
	if err != nil {
		return errors.New(fmt.Sprintf("Error executing command '%s':", err.Error()))
	}

	if outputVariable != "" {
		ec.AddToContext(outputVariable, strings.TrimSpace(output))
	}

	return nil
}
