//go:build mage

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

package main

import (
	// mage:import
	"knative.dev/toolbox/magetasks"
	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/config/buildvars"
	"knative.dev/toolbox/magetasks/pkg/artifact"
	artifactimage "knative.dev/toolbox/magetasks/pkg/artifact/image"
	"knative.dev/toolbox/magetasks/pkg/artifact/platform"
	"knative.dev/toolbox/magetasks/pkg/checks"
	"knative.dev/toolbox/magetasks/pkg/image"
	"knative.dev/toolbox/magetasks/pkg/knative"
	"knative.dev/toolbox/magetasks/tests/testdata/build/overrides"
	"knative.dev/toolbox/magetasks/tests/testdata/pkg/metadata"
)

// Default target is set to Build
//
//goland:noinspection GoUnusedGlobalVariable
var Default = magetasks.Build

func init() { //nolint:gochecknoinits
	archs := []platform.Architecture{
		platform.AMD64, platform.ARM64, platform.S390X, platform.PPC64LE,
	}
	placeholder := artifact.Image{
		Metadata: config.Metadata{Name: "placeholder"},
		Labels: map[string]config.Resolver{
			"description": config.StaticResolver("A placeholder image description"),
		},
		Architectures: archs,
	}
	sampleimage := artifact.Image{
		Metadata:      config.Metadata{Name: "sampleimage"},
		Architectures: archs,
	}
	other := artifact.Binary{
		Metadata: config.Metadata{
			Name: "other",
			BuildVariables: buildvars.Assemble([]buildvars.Operator{
				image.InfluenceableReference{
					Path:        metadata.ImagePath(metadata.PlaceholderImage),
					EnvVariable: "MAGETASKS_EXAMPLE_PLACEHOLDER_IMAGE",
					Image:       placeholder,
				},
				image.InfluenceableReference{
					Path:        metadata.ImagePath(metadata.SampleImage),
					EnvVariable: "MAGETASKS_EXAMPLE_SAMPLE_IMAGE",
					Image:       sampleimage,
				},
			}),
		},
		Platforms: []artifact.Platform{
			{OS: platform.Linux, Architecture: platform.AMD64},
			{OS: platform.Linux, Architecture: platform.ARM64},
			{OS: platform.Mac, Architecture: platform.AMD64},
			{OS: platform.Mac, Architecture: platform.ARM64},
			{OS: platform.Windows, Architecture: platform.AMD64},
		},
	}
	magetasks.Configure(config.Config{
		Version: &config.Version{
			Path: metadata.VersionPath(), Resolver: knative.NewVersionResolver(),
		},
		Artifacts: []config.Artifact{
			placeholder, sampleimage, other,
		},
		Checks: []config.Task{
			checks.GolangCiLint(func(opts *checks.GolangCiLintOptions) {
				opts.Version = "v1.55.2"
			}),
			checks.Revive(),
			checks.Staticcheck(),
		},
		BuildVariables: map[string]config.Resolver{
			metadata.ImageBasenamePath():          artifactimage.BaseName,
			metadata.ImageBasenameSeparatorPath(): artifactimage.BaseNameSeparator,
		},
		Overrides: overrides.List,
	})
}
