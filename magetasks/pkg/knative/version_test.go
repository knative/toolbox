package knative_test

import (
	"strconv"
	"testing"

	"gotest.tools/v3/assert"
	"knative.dev/toolbox/magetasks/pkg/environment"
	"knative.dev/toolbox/magetasks/pkg/git"
	"knative.dev/toolbox/magetasks/pkg/knative"
	"knative.dev/toolbox/magetasks/pkg/strings"
	"knative.dev/toolbox/magetasks/pkg/testing/errors"
	"knative.dev/toolbox/magetasks/pkg/version"
)

func TestVersionResolver(t *testing.T) {
	tests := []testCase{{}, {
		environment: environment.New("TAG=v4.6.23", "TAG_RELEASE=1"),
		version:     "v4.6.23",
		latest:      true,
	}, {
		environment: environment.New("TAG=v6.23.1", "TAG_RELEASE=1"),
		version:     "v6.23.1",
		tags:        strings.NewSet("v6.23.0", "v7.0.0"),
		latest:      false,
	}, {
		environment: environment.New("TAG=v23.1.2", "TAG_RELEASE=1"),
		version:     "v23.1.2",
		tags:        strings.NewSet("v6.23.0", "v7.0.0"),
		latest:      true,
	}}
	for i, tc := range tests {
		tc := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			resolver := tc.resolver()
			assert.Equal(t, resolver.Version(), tc.version)
			if tc.version != "" {
				latest, err := resolver.IsLatest(version.AnyVersion)
				errors.Check(t, err, tc.err)
				assert.Equal(t, latest, tc.latest)
			}
		})
	}
}

type testCase struct {
	environment environment.Values
	describe    string
	tags        strings.Set
	version     string
	latest      bool
	err         error
}

func (tc testCase) resolver() version.Resolver {
	return knative.NewTestableVersionResolver(
		git.StaticRepository{DescribeString: tc.describe, TagsSet: tc.tags},
		func() environment.Values {
			return tc.environment
		},
	)
}
