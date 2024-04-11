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

package image

import (
	"errors"

	"knative.dev/toolbox/magetasks/pkg/version"
)

// Tags will return a list of tags for an OCI image based on the version
// information given by resolver.
func Tags(resolver version.Resolver) ([]string, error) {
	ranges, err := version.CompatibleRanges(resolver)
	if err != nil {
		if !errors.Is(err, version.ErrVersionIsNotValid) {
			return nil, err
		}
		ranges = make([]string, 0)
	}
	tags := append([]string{resolver.Version()}, ranges...)
	var latest bool
	latest, err = resolver.IsLatest(version.AnyVersion)
	if err != nil && !errors.Is(err, version.ErrVersionIsNotValid) {
		return nil, err
	}
	if latest {
		tags = append(tags, "latest")
	}
	return tags, nil
}
