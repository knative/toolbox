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

package knative

import (
	"knative.dev/toolbox/magetasks/pkg/cache"
	"knative.dev/toolbox/magetasks/pkg/environment"
	"knative.dev/toolbox/magetasks/pkg/git"
	"knative.dev/toolbox/magetasks/pkg/version"
)

// NewTestableVersionResolver creates an instance of version.Resolver that can
// be easily tested.
func NewTestableVersionResolver(
	repo git.Repository,
	env func() environment.Values,
) version.Resolver {
	return NewVersionResolver(
		WithGit(
			git.WithCache(cache.NoopCache{}),
			git.WithRepository(repo),
		),
		WithEnvironmental(environment.WithValuesSupplier(env)),
	)
}
