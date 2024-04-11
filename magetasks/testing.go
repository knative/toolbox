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

package magetasks

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/magefile/mage/sh"
	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/files"
	"knative.dev/toolbox/magetasks/pkg/ldflags"
	"knative.dev/toolbox/magetasks/pkg/targets"
	"knative.dev/toolbox/magetasks/pkg/tasks"
)

// Test will execute regular unit tests.
func Test(ctx context.Context) {
	targets.Deps(ctx, Check)
	t := tasks.Start("✅", "Testing", true)
	args := []string{
		"--format", "testname",
		"--jsonfile", fmt.Sprintf("%s/tests-output.json", files.ReportsDir()),
		"--junitfile", fmt.Sprintf("%s/tests-results.junit.xml", files.ReportsDir()),
	}
	if color.NoColor {
		args = append(args, "--no-color")
	}
	args = append(args, "--",
		"-covermode=atomic", fmt.Sprintf("-coverprofile=%s/tests-coverage.out", files.ReportsDir()),
		"-race",
		"-count=1",
		"-short",
	)
	args = append(appendBuildVariables(args), "./...")
	cmd := fmt.Sprintf("%s/tools/gotestsum", files.BuildDir())
	err := sh.RunV(cmd, args...)
	t.End(err)
}

func appendBuildVariables(args []string) []string {
	c := config.Actual()
	if c.Version != nil || len(c.BuildVariables) > 0 {
		builder := ldflags.NewBuilder()
		if c.Version != nil {
			builder.Add(c.Version.Path, c.Version.Resolver.Version)
		}
		for key, resolver := range c.BuildVariables {
			builder.Add(key, resolver)
		}
		args = builder.BuildOnto(args)
	}
	return args
}
