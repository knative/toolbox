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

package artifact

import (
	"errors"
	"fmt"
	"path"

	"github.com/magefile/mage/sh"
	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/files"
	"knative.dev/toolbox/magetasks/pkg/ldflags"
	"knative.dev/toolbox/magetasks/pkg/output/color"
)

const (
	// ResultBuilt is used to cache results of building the artifacts.
	ResultBuilt = "built"

	// ResultPublication is used to cache results of publication of artifacts.
	ResultPublication = "publication"

	// BinariesByPlatform is used to assign built binary to a platform.
	BinariesByPlatform = "binary.by-platform"
)

var (
	// ErrGoBuildFailed when go fails to build the project.
	ErrGoBuildFailed = errors.New("go build failed")

	// ErrInvalidArtifact when invalid type of artifact is given.
	ErrInvalidArtifact = errors.New("invalid artifact")
)

// Binary represents a binary that will be built.
type Binary struct {
	config.Metadata
	Platforms []Platform
}

func (b Binary) GetType() string {
	return "📦"
}

// BinaryBuilder is a regular binary Golang builder.
type BinaryBuilder struct{}

func (bb BinaryBuilder) Accepts(artifact config.Artifact) bool {
	_, ok := artifact.(Binary)
	return ok
}

func (bb BinaryBuilder) Build(artifact config.Artifact, notifier config.Notifier) config.Result {
	b, ok := artifact.(Binary)
	if !ok {
		return config.Result{Error: ErrInvalidArtifact}
	}
	info := make(map[string]interface{})
	artifacts := make([]string, 0, len(b.Platforms))
	byPlatform := make(map[Platform]string)
	var err error
	for _, pltfrm := range b.Platforms {
		var bin string
		bin, err = b.buildForPlatform(pltfrm, artifact.GetName(), notifier)
		if err == nil {
			artifacts = append(artifacts, bin)
			byPlatform[pltfrm] = bin
		} else {
			break
		}
	}
	info[ArtifactsBuilt] = artifacts
	info[BinariesByPlatform] = byPlatform
	return config.Result{Error: err, Info: info}
}

func (b Binary) buildForPlatform(
	platform Platform,
	name string,
	notifier config.Notifier,
) (string, error) {
	args := []string{
		"build",
	}
	c := config.Actual()
	if c.Version != nil || len(c.BuildVariables) > 0 || len(b.BuildVariables) > 0 {
		builder := ldflags.NewBuilder()
		for key, resolver := range c.BuildVariables {
			builder.Add(key, resolver)
		}
		if c.Version != nil {
			builder.Add(c.Version.Path, c.Version.Resolver.Version)
		}
		for key, resolver := range b.BuildVariables {
			builder.Add(key, resolver)
		}
		args = builder.BuildOnto(args)
	}
	binary := fullBinaryName(platform, name)
	args = append(args, "-o", binary, fullBinaryDirectory(name))
	env := map[string]string{
		"GOOS":   string(platform.OS),
		"GOARCH": string(platform.Architecture),
	}
	notifier.Notify(fmt.Sprintf("go build (%s)",
		color.Blue(fmt.Sprintf("%s-%s", platform.OS, platform.Architecture)),
	))
	err := sh.RunWithV(env, "go", args...)
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrGoBuildFailed, err)
	}
	return binary, err
}

func fullBinaryName(platform Platform, name string) string {
	return path.Join(files.BuildDir(), "bin",
		fmt.Sprintf("%s-%s-%s", name, platform.OS, platform.Architecture))
}

func fullBinaryDirectory(name string) string {
	return path.Join(files.ProjectDir(), "cmd", name)
}
