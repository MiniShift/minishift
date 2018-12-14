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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/minishift/minishift/cmd/testing/cli"
	pgkTesting "github.com/minishift/minishift/pkg/testing"
	"github.com/stretchr/testify/assert"

	"github.com/minishift/minishift/cmd/minishift/state"
	"github.com/minishift/minishift/pkg/util/filehelper"
	"github.com/minishift/minishift/pkg/util/os/atexit"
	"k8s.io/client-go/tools/clientcmd"
)

func Test_clear_cache_user_confirms(t *testing.T) {
	tmpMinishiftHomeDir, tee := prepareCacheDir(t)
	defer cli.TearDown(tmpMinishiftHomeDir, tee)

	// create a canned user response for stdin
	origStdin, tmpFile := cli.PrepareStdinResponse("y", t)
	defer cli.ResetStdin(origStdin, tmpFile)

	atexit.RegisterExitHandler(cli.PreventAtExit(t))
	clearCache()

	assert.False(t, filehelper.Exists(state.InstanceDirs.Cache),
		fmt.Sprintf("Expected cache dir '%s' to be deleted", state.InstanceDirs.Cache))
}

func Test_clear_cache_user_aborts(t *testing.T) {
	tmpMinishiftHomeDir, tee := prepareCacheDir(t)
	defer cli.TearDown(tmpMinishiftHomeDir, tee)

	// create a canned user response for stdin
	origStdin, tmpFile := cli.PrepareStdinResponse("n", t)
	defer cli.ResetStdin(origStdin, tmpFile)

	clearCache()

	assert.DirExists(t, state.InstanceDirs.Cache,
		fmt.Sprintf("Expected cache dir '%s' to still exist", state.InstanceDirs.Cache))
}

func Test_clear_cache_forced(t *testing.T) {
	tmpMinishiftHomeDir, tee := prepareCacheDir(t)
	defer cli.TearDown(tmpMinishiftHomeDir, tee)

	atexit.RegisterExitHandler(cli.PreventAtExit(t))
	// simulate setting the force flag
	forceFlag = true
	clearCache()

	assert.False(t, filehelper.Exists(state.InstanceDirs.Cache), "Expected cache dir '%s' to be deleted", state.InstanceDirs.Cache)
}

func Test_remove_entries_for_cluster(t *testing.T) {
	kubeConfigPath := filepath.Join("..", "..", "..", "test", "testdata", "kubeconfig")

	dirtyKubeConfigPath := filepath.Join("..", "..", "..", "test", "testdata", "kubeconfig_dirty")
	cleanKubeConfigPath := filepath.Join("..", "..", "..", "test", "testdata", "kubeconfig_clean")

	dirtyKubeConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
	assert.NoError(t, err, "Error loading dirty kubeconfig file")

	os.Rename(kubeConfigPath, dirtyKubeConfigPath)
	os.Rename(cleanKubeConfigPath, kubeConfigPath)

	cleanKubeConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
	assert.NoError(t, err, "Error loading clean kubeconfig file")

	os.Rename(kubeConfigPath, cleanKubeConfigPath)
	os.Rename(dirtyKubeConfigPath, kubeConfigPath)

	// Cluster present in kubeconfig
	actualKubeConfig, err := removeEntriesForCluster("192.168.42.28", dirtyKubeConfig)
	assert.NoError(t, err, "Error removing entries from kubeconfig.")
	assert.Equal(t, cleanKubeConfig, actualKubeConfig, "Existing cluster test.")

	// Non existent cluster
	actualKubeConfig, err = removeEntriesForCluster("192.168.102.88", dirtyKubeConfig)
	assert.NoError(t, err, "Error removing entries from kubeconfig.")
	assert.Equal(t, dirtyKubeConfig, actualKubeConfig, "Non Existent cluster test")
}

func Test_delete_succeeds_for_non_existing_vm(t *testing.T) {
	tmpMinishiftHomeDir := cli.SetupTmpMinishiftHome(t)
	tee := cli.CreateTee(t, false)
	defer cli.TearDown(tmpMinishiftHomeDir, tee)

	// register a custom exit handler which verifies the exit code as well as the content of the cache
	atexit.RegisterExitHandler(cli.VerifyExitCodeAndMessage(t, tee, 0, ""))

	runDelete(nil, nil)
}

func prepareCacheDir(t *testing.T) (string, *pgkTesting.Tee) {
	// setup the test Minishift home directory and make sure we have a cache directory
	tmpMinishiftHomeDir := cli.SetupTmpMinishiftHome(t)
	os.MkdirAll(state.InstanceDirs.Cache, os.ModePerm)

	// create a tee to keep the test silent and make sure that the temporary Minishift home directory gets cleaned up
	tee := cli.CreateTee(t, true)

	return tmpMinishiftHomeDir, tee
}
