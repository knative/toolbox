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

package git

import (
	"errors"

	"knative.dev/toolbox/magetasks/pkg/version"
)

// TagBasedIsLatestStrategy is the default strategy, that uses a repository
// tags to determine if given version is latest within version range given.
func (r VersionResolver) TagBasedIsLatestStrategy(vr version.Resolver) func(versionRange string) (bool, error) {
	return func(versionRange string) (bool, error) {
		val, err := version.IsLatestGivenReleases(
			vr.Version(), versionRange, true,
			r.resolveTags,
		)
		if err != nil {
			return false, err
		}
		return val, nil
	}
}

// ResolveIsLatest will resolve is version given by version.Resolver is the
// latest one in given version range, using a logic of git.VersionResolver.
func ResolveIsLatest(r VersionResolver, vr version.Resolver, versionRange string) (bool, error) {
	strategy := r.TagBasedIsLatestStrategy
	if r.IsLatestStrategy != nil {
		strategy = r.IsLatestStrategy
	}
	fn := strategy(vr)
	latest, err := fn(versionRange)
	if err != nil {
		if !errors.Is(err, version.ErrVersionIsNotValid) {
			return false, err
		}
		return false, nil
	}
	return latest, nil
}
