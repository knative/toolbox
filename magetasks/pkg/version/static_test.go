package version_test

import (
	"knative.dev/toolbox/magetasks/pkg/version"
)

// staticResolver just return values from values given upfront.
type staticResolver struct {
	VersionString         string
	Tags                  []string
	FailOnInvalidReleases bool
}

func (r staticResolver) Version() string {
	return r.VersionString
}

func (r staticResolver) IsLatest(versionRange string) (bool, error) {
	return version.IsLatestGivenReleases(
		r.VersionString, versionRange, !r.FailOnInvalidReleases,
		func() []string {
			return r.Tags
		},
	)
}
