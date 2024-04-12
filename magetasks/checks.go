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

	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/deps"
	"knative.dev/toolbox/magetasks/pkg/targets"
	"knative.dev/toolbox/magetasks/pkg/tasks"
)

// Check will run all lints checks.
func Check(ctx context.Context) {
	targets.Deps(ctx, deps.Install)
	t := tasks.Start("🔍", "Checking", len(config.Actual().Checks) > 0)
	for _, check := range config.Actual().Checks {
		p := t.Part(check.Name)
		pp := p.Starting()
		pp.Done(check.Operation(pp))
	}
	t.End()
}
