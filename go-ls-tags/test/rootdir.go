/*
Copyright 2022 The Knative Authors

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

package test

import (
	"path"
	"runtime"
	"testing"
)

// Rootdir returns a path to a root directory of toolbox repo.
func Rootdir() string {
	if !testing.Testing() {
		panic("test.WithDirectory can be used only in tests")
	}
	_, curfile, _, _ := runtime.Caller(0) //nolint:dogsled
	return path.Dir(path.Dir(path.Dir(curfile)))
}
