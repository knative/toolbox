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

package magetasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/magefile/mage/mg"
	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/artifact"
	"knative.dev/toolbox/magetasks/pkg/tasks"
)

// ErrNoPublisherForArtifact when no publisher for artifact is found.
var ErrNoPublisherForArtifact = errors.New("no publisher for artifact found")

// Publish will publish built artifacts to remote site.
func Publish(ctx context.Context) {
	mg.CtxDeps(ctx, Build)
	artifacts := config.Actual().Artifacts
	t := tasks.Start("📤", "Publishing", len(artifacts) > 0)
	for _, art := range artifacts {
		p := t.Part(fmt.Sprintf("%s %s", art.GetType(), art.GetName()))
		pp := p.Starting()

		publishArtifact(art, pp)
	}
	t.End()
}

func publishArtifact(art config.Artifact, pp tasks.PartProcessing) {
	found := false
	for _, publisher := range config.Actual().Publishers {
		if !publisher.Accepts(art) {
			continue
		}
		found = true
		result := publisher.Publish(art, pp)
		if result.Failed() {
			pp.Done(result.Error)
			return
		}
		config.WithContext(func(ctx context.Context) context.Context {
			return context.WithValue(ctx, artifact.PublishKey(art), result)
		})
	}
	var err error
	if !found {
		err = ErrNoPublisherForArtifact
	}
	pp.Done(err)
}
