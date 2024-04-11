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
