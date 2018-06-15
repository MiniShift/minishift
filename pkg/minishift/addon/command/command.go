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

// Command defines a single command to be executed as part of an addon definition.
// Minishift supports various types of commands as part of its addon DSL, eg oc, openshift or docker commands.
type Command interface {
	// Executes the command
	Execute(ec *ExecutionContext) error

	// String returns a string representation of the command
	String() string
}

type doExecute func(ec *ExecutionContext, ignoreError bool, outputVariable string) error

type defaultCommand struct {
	Command

	rawCommand     string
	fn             doExecute
	ignoreError    bool
	outputVariable string
}

func (c *defaultCommand) Execute(ec *ExecutionContext) error {
	err := c.fn(ec, c.ignoreError, c.outputVariable)
	if c.ignoreError {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *defaultCommand) String() string {
	return c.rawCommand
}

func (c *defaultCommand) IgnoreError() bool {
	return c.ignoreError
}
