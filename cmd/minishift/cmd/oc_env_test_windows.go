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

package cmd

import (
	"testing"

	"github.com/docker/machine/libmachine"
	"github.com/stretchr/testify/assert"
)

func Test_windows_oc_path(t *testing.T) {
	api := libmachine.NewClient("foo", "foo")
	defer api.Close()
	shellConfig, err := getOcShellConfig(api, "C:\\Users\\john\\.minishift\\cache\\oc\\v1.5.0\\oc.exe", "", false)

	assert.NoError(t, err)
	expectedOcDirPath := "C:\\Users\\john\\.minishift\\cache\\oc\\v1.5.0"
	assert.Equal(t, shellConfig.OcDirPath, expectedOcDirPath)
}
