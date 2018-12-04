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

package os

import (
	"github.com/kardianos/osext"
	"os/user"
	"runtime"
	"strings"
)

type OS string

const (
	LINUX   OS = "linux"
	DARWIN  OS = "darwin"
	WINDOWS OS = "windows"
)

func (t OS) String() string {
	return string(t)
}

func CurrentOS() OS {
	switch runtime.GOOS {
	case "windows":
		return WINDOWS
	case "darwin":
		return DARWIN
	case "linux":
		return LINUX
	}
	panic("Unexpected OS type")
}

func CurrentExecutable() (string, error) {
	currentExec, err := osext.Executable()
	if err != nil {
		return "", err
	}
	return currentExec, nil
}

func CurrentUser() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	tokens := strings.Split(user.Username, `\`) // user.Name returns Domain\Username
	return tokens[1], nil
}
