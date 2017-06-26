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

package registration

import (
	"fmt"
	"strings"

	"github.com/docker/machine/libmachine/provision"
	minishiftStrings "github.com/minishift/minishift/pkg/util/strings"
)

func init() {
	Register("Redhat", &RegisteredRegistrator{
		New: NewRedHatRegistrator,
	})
}

func NewRedHatRegistrator(c provision.SSHCommander) Registrator {
	return &RedHatRegistrator{
		SSHCommander: c,
	}
}

type RedHatRegistrator struct {
	provision.SSHCommander
}

func (registrator *RedHatRegistrator) CompatibleWithDistribution(osReleaseInfo *provision.OsRelease) bool {
	if osReleaseInfo.ID != "rhel" {
		return false
	}
	if _, err := registrator.SSHCommand("sudo -E subscription-manager"); err != nil {
		return false
	} else {
		return true
	}
}

func (registrator *RedHatRegistrator) Register(param *RegistrationParameters) error {
	output, err := registrator.SSHCommand("sudo -E subscription-manager version")
	if err != nil {
		return err
	}
	if strings.Contains(output, "not registered") {
		for i := 1; i < 4; i++ {
			if param.Username == "" {
				param.Username = param.GetUsernameInteractive("Red Hat Developers or Red Hat Subscription Management (RHSM) username")
			}
			if param.Password == "" {
				param.Password = param.GetPasswordInteractive("Red Hat Developers or Red Hat Subscription Management (RHSM) password")
			}
			subscriptionCommand := fmt.Sprintf("sudo -E subscription-manager register --auto-attach "+
				"--username %s "+
				"--password '%s' ", param.Username, minishiftStrings.EscapeSingleQuote(param.Password))
			_, err = registrator.SSHCommand(subscriptionCommand)
			if err == nil {
				return nil
			}
			if strings.Contains(err.Error(), "Invalid username or password") {
				fmt.Println("Invalid username or password Retry: ", i)
				param.Username = ""
				param.Password = ""
			} else {
				return err
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (registrator *RedHatRegistrator) Unregister(param *RegistrationParameters) error {
	if output, err := registrator.SSHCommand("sudo -E subscription-manager version"); err != nil {
		return err
	} else {
		if !strings.Contains(output, "not registered") {
			if _, err := registrator.SSHCommand(
				"sudo -E subscription-manager unregister"); err != nil {
				return err
			}
		}
	}
	return nil
}
