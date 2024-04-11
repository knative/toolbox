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
	"knative.dev/toolbox/magetasks/pkg/environment"
	"knative.dev/toolbox/magetasks/pkg/git"
)

// NewVersionResolver creates a version.Resolver implementation directly
// targeting Knative project CI.
func NewVersionResolver(options ...VersionResolverOption) VersionResolver {
	r := VersionResolver{
		env: environment.VersionResolver{
			VersionKey: "TAG",
			IsApplicable: []environment.Check{
				{Key: "TAG_RELEASE", Value: "1"},
				{Key: "TAG"},
			},
		},
	}
	for _, option := range options {
		option(&r)
	}

	return r
}

// VersionResolverOption id option to customize version resolution.
type VersionResolverOption func(*VersionResolver)

// WithGit allows passing options for git.VersionResolver.
func WithGit(options ...git.VersionResolverOption) VersionResolverOption {
	return func(resolver *VersionResolver) {
		for _, option := range options {
			option(&resolver.git)
		}
	}
}

// WithEnvironmental allows passing options for environment.VersionResolver.
func WithEnvironmental(options ...environment.VersionResolverOption) VersionResolverOption {
	return func(resolver *VersionResolver) {
		for _, option := range options {
			option(&resolver.env)
		}
	}
}

// VersionResolver is a knative compatible version resolver.
type VersionResolver struct {
	git git.VersionResolver
	env environment.VersionResolver
}

func (v VersionResolver) Version() string {
	if ver := v.env.Version(); ver != "" {
		return ver
	}
	return v.git.Version()
}

func (v VersionResolver) IsLatest(versionRange string) (bool, error) {
	return git.ResolveIsLatest(v.git, v, versionRange)
}
