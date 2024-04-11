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

package version

import (
	"fmt"

	"github.com/blang/semver/v4"
)

// AnyVersion is a range that matches any released version.
const AnyVersion = ">=0.0.0"

// IsLatestResolver is a func that returns true if given version is the latest
// one within given version range.
type IsLatestResolver func(version semver.Version, versionRange semver.Range) (bool, error)

// IsLatest does basic sanity checking on version and range, before executing
// provided resolver to answer if given version is the latest within the given
// version range.
func IsLatest(version, versionRange string, resolver IsLatestResolver) (bool, error) {
	sver, err := semver.ParseTolerant(version)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrVersionIsNotValid, err)
	}
	if versionRange == "" {
		versionRange = AnyVersion
	}
	verRange, err := semver.ParseRange(versionRange)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrRangeIsNotValid, err)
	}
	if !verRange(sver) {
		return false, fmt.Errorf("%w: %#v is outside of %#v",
			ErrVersionOutsideOfRange, version, versionRange)
	}
	if len(sver.Build) > 0 || len(sver.Pre) > 0 {
		return false, nil
	}
	return resolver(sver, verRange)
}

// IsLatestGivenReleases checks if given version is the latest within the given
// version range, using a list of provided releases.
func IsLatestGivenReleases(
	version, versionRange string,
	skipInvalidReleases bool,
	releases func() []string,
) (bool, error) {
	return IsLatest(version, versionRange, func(ver semver.Version, versionRange semver.Range) (bool, error) {
		for _, r := range releases() {
			sr, err := semver.ParseTolerant(r)
			if err != nil {
				if skipInvalidReleases {
					continue
				}
				return false, fmt.Errorf("%w: %w", ErrVersionIsNotValid, err)
			}
			if !versionRange(sr) {
				continue
			}
			if sr.GT(ver) {
				return false, nil
			}
		}
		return true, nil
	})
}
