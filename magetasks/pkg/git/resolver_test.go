package git_test

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
	"knative.dev/toolbox/magetasks/pkg/cache"
	"knative.dev/toolbox/magetasks/pkg/git"
	"knative.dev/toolbox/magetasks/pkg/strings"
	"knative.dev/toolbox/magetasks/pkg/testing/errors"
	"knative.dev/toolbox/magetasks/pkg/version"
)

func TestResolver(t *testing.T) {
	tests := []testCase{{
		version:      "v1.5.2-2-g8cc3513",
		versionRange: version.AnyVersion,
	}, {
		version:      "v1.5.3",
		versionRange: version.AnyVersion,
		latest:       true,
	}, {
		version:      "v1.5.3",
		tags:         strings.NewSet("v1.5.2", "v1.6.0"),
		versionRange: version.AnyVersion,
	}, {
		version:      "v1.5.3",
		tags:         strings.NewSet("wrong", "v1.5.2", "v1.5.4", "v1.6.0"),
		versionRange: version.AnyVersion,
	}, {
		version:      "1.5.3",
		versionRange: version.AnyVersion,
		latest:       true,
	}, {
		version:      "af134dd",
		versionRange: version.AnyVersion,
	}}
	for _, tc := range tests {
		tc := tc
		resolver := git.VersionResolver{
			Cache: cache.NoopCache{},
			Repository: git.StaticRepository{
				DescribeString: tc.version,
				TagsSet:        tc.tags,
			},
		}
		t.Run(tc.String(), func(t *testing.T) {
			assert.Equal(t, resolver.Version(), tc.version)
			latest, err := resolver.IsLatest(tc.versionRange)
			errors.Check(t, err, tc.err)
			assert.Equal(t, latest, tc.latest)
		})
	}
}

type testCase struct {
	version      string
	tags         strings.Set
	versionRange string
	latest       bool
	err          error
}

func (tc testCase) String() string {
	name := tc.version
	if tc.tags.Len() > 0 {
		name = fmt.Sprintf("%s-%v", name, tc.tags)
	}
	if tc.versionRange != version.AnyVersion {
		name = fmt.Sprintf("%s-%s", name, tc.versionRange)
	}
	return name
}
