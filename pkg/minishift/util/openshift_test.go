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

package util

import "testing"

func TestValidateOpenshiftMinVersion(t *testing.T) {
	verList := map[string]bool{
		"v1.1.0":         false,
		"v1.2.2":         false,
		"v1.2.3-beta":    false,
		"v1.3.1":         false,
		"v1.3.5-alpha":   false,
		"v1.4.1":         false,
		"v1.5.0-alpha.0": false,
		"v1.5.1-beta.0":  true,
		"v1.5.0-rc.0":    true,
		"v1.5.0":         true,
		"v1.6.0":         true,
	}
	for ver, val := range verList {
		if ValidateOpenshiftMinVersion(ver) != val {
			t.Fatalf("Expected '%t' Got '%t' for %s", val, ValidateOpenshiftMinVersion(ver), ver)
		}
	}
}
