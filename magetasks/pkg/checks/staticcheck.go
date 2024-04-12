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

package checks

import (
	"fmt"
	"path"

	"github.com/magefile/mage/sh"
	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/files"
)

// Staticcheck will configure staticcheck in the build.
func Staticcheck() config.Task {
	return config.Task{
		Name:      "staticcheck",
		Operation: staticcheck,
		Overrides: []config.Configurator{
			config.NewDependencies("honnef.co/go/tools/cmd/staticcheck@latest"),
		},
	}
}

func staticcheck(notifier config.Notifier) error {
	configFile := "staticcheck.conf"
	c := path.Join(files.ProjectDir(), configFile)
	if files.DontExists(c) {
		skipBecauseOfMissingConfig(notifier, configFile)
		return nil
	}
	cmd := fmt.Sprintf("%s/tools/staticcheck", files.BuildDir())
	return sh.RunV(cmd, "-f", "stylish", "./...")
}
