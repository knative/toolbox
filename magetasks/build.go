package magetasks

import (
	"context"
	"errors"
	"fmt"

	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/artifact"
	"knative.dev/toolbox/magetasks/pkg/targets"
	"knative.dev/toolbox/magetasks/pkg/tasks"
)

// ErrNoBuilderForArtifact when no builder for artifact is found.
var ErrNoBuilderForArtifact = errors.New("no builder for artifact found")

// Build will build project artifacts, binaries and images.
func Build(ctx context.Context) error {
	targets.Deps(ctx, Test)
	t := tasks.Start("🔨", "Building", len(config.Actual().Artifacts) > 0)
	for _, art := range config.Actual().Artifacts {
		p := t.Part(fmt.Sprintf("%s %s", art.GetType(), art.GetName()))
		pp := p.Starting()

		err := buildArtifact(art, pp)
		if err != nil {
			t.End(err)
			return err
		}
	}
	t.End()

	return nil
}

func buildArtifact(art config.Artifact, pp tasks.PartProcessing) error {
	found := false
	for _, builder := range config.Actual().Builders {
		if !builder.Accepts(art) {
			continue
		}
		found = true
		result := builder.Build(art, pp)
		if result.Failed() {
			err := result.Error
			pp.Done(err)
			return err
		}
		config.WithContext(func(ctx context.Context) context.Context {
			return context.WithValue(ctx, artifact.BuildKey(art), result)
		})
	}
	var err error
	if !found {
		err = ErrNoBuilderForArtifact
	}
	pp.Done(err)
	return err
}
