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

package constants

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/minishift/minishift/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestDefaultMinishiftHome(t *testing.T) {
	os.Unsetenv("MINISHIFT_HOME")
	expectedMiniPath := filepath.Join(util.HomeDir(), ".minishift")
	actualMiniPath := GetMinishiftHomeDir()
	assert.Equal(t, expectedMiniPath, actualMiniPath)
}

func TestMinishiftHomeViaEnvironment(t *testing.T) {
	expectedMiniPath := "/tmp/minishift-test-home"
	os.Setenv("MINISHIFT_HOME", expectedMiniPath)
	defer os.Unsetenv("MINISHIFT_HOME")

	actualMiniPath := GetMinishiftHomeDir()
	assert.Equal(t, expectedMiniPath, actualMiniPath)
}
