//go:build mage

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
	"knative.dev/toolbox/magetasks/tests/example/build/overrides"
	"knative.dev/toolbox/magetasks/tests/example/pkg/metadata"
)

// Default target is set to Build
//
//goland:noinspection GoUnusedGlobalVariable
var Default = magetasks.Build

func init() { //nolint:gochecknoinits
	archs := []platform.Architecture{
		platform.AMD64, platform.ARM64, platform.S390X, platform.PPC64LE,
	}
	dummy := artifact.Image{
		Metadata: config.Metadata{Name: "dummy"},
		Labels: map[string]config.Resolver{
			"description": config.StaticResolver("A dummy image description"),
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
					Path:        metadata.ImagePath(metadata.DummyImage),
					EnvVariable: "MAGETASKS_EXAMPLE_DUMMY_IMAGE",
					Image:       dummy,
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
			dummy, sampleimage, other,
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
