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

package image_test

import (
	"strconv"
	"testing"

	"knative.dev/toolbox/magetasks/pkg/artifact/image"
	"knative.dev/toolbox/magetasks/pkg/environment"
	"knative.dev/toolbox/magetasks/pkg/git"
	"knative.dev/toolbox/magetasks/pkg/knative"
	"knative.dev/toolbox/magetasks/pkg/strings"
	"knative.dev/toolbox/magetasks/pkg/testing/errors"
	"knative.dev/toolbox/magetasks/pkg/version"
)

func TestTags(t *testing.T) {
	tests := []tagsTestCase{{
		version: "v1.5.2-2-g8cc3513",
		want:    []string{"v1.5.2-2-g8cc3513"},
	}, {
		version: "v1.5.3",
		want:    []string{"v1.5.3", "v1.5", "v1", "latest"},
	}, {
		version: "v1.5.3",
		tags:    strings.NewSet("v1.5.2", "v1.6.0"),
		want:    []string{"v1.5", "v1.5.3"},
	}, {
		version: "v1.5.3",
		tags:    strings.NewSet("wrong", "v1.5.2", "v1.5.4", "v1.6.0"),
		want:    []string{"v1.5.3"},
	}, {
		version: "1.5.3",
		want:    []string{"1.5.3", "1.5", "1", "latest"},
	}, {
		version: "af134dd",
		want:    []string{"af134dd"},
	}}
	for i, tc := range tests {
		resolver := tc.resolver()
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := image.Tags(resolver)
			errors.Check(t, err, tc.err)
			if !equal(got, tc.want) {
				t.Fatalf("want: %#v, got: %#v", tc.want, got)
			}
		})
	}
}

type tagsTestCase struct {
	version string
	tags    strings.Set
	want    []string
	err     error
}

func (tc tagsTestCase) resolver() version.Resolver {
	return knative.NewTestableVersionResolver(
		git.StaticRepository{DescribeString: tc.version, TagsSet: tc.tags},
		func() environment.Values {
			return nil
		},
	)
}

func equal(a, b []string) bool {
	return strings.NewSet(a...).Equal(strings.NewSet(b...))
}
