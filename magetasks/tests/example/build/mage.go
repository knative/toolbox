//go:build ignored

/*
 Copyright 2024 The Knative Authors

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

package main

import (
	"os"
	"path"
	"runtime"

	"knative.dev/toolbox/magetasks/entrypoint"
)

func main() {
	bdir := builddir()
	os.Exit(entrypoint.Execute(entrypoint.Context{
		Directories: entrypoint.Directories{
			BuildDir:   bdir,
			ProjectDir: path.Dir(bdir),
			CacheDir:   path.Join(bdir, "_output"),
		},
	}))
}

func builddir() string {
	_, file, _, _ := runtime.Caller(0) //nolint:dogsled
	return path.Dir(file)
}
