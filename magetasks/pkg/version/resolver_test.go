package version_test

import (
	"strconv"
	"testing"

	"knative.dev/toolbox/magetasks/pkg/strings"
	"knative.dev/toolbox/magetasks/pkg/testing/errors"
	"knative.dev/toolbox/magetasks/pkg/version"
)

func TestCompatibleRanges(t *testing.T) {
	tests := []compatibleRangesTestCase{{
		version: "v0.5.2-2-g8cc3513",
		want:    []string{},
	}, {
		version: "v0.5.3",
		want:    []string{"v0.5", "v0"},
	}, {
		version: "v0.5.3",
		tags:    []string{"v0.5.2", "v0.6.0"},
		want:    []string{"v0.5"},
	}, {
		version: "0.5.3",
		want:    []string{"0.5", "0"},
	}, {
		version: "af134dd",
		err:     version.ErrVersionIsNotValid,
	}}
	for i, tc := range tests {
		tr := tc.resolver()
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := version.CompatibleRanges(tr)
			errors.Check(t, err, tc.err)
			if !equal(got, tc.want) {
				t.Fatalf("want: %#v, got: %#v", tc.want, got)
			}
		})
	}
}

func equal(a, b []string) bool {
	return strings.NewSet(a...).Equal(strings.NewSet(b...))
}

type compatibleRangesTestCase struct {
	version string
	tags    []string
	want    []string
	err     error
}

func (tc compatibleRangesTestCase) resolver() version.Resolver {
	return staticResolver{
		VersionString: tc.version,
		Tags:          tc.tags,
	}
}
