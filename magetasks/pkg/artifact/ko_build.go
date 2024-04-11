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
	"os"
	"path"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/ko/pkg/build"
	"github.com/google/ko/pkg/commands"
	"github.com/google/ko/pkg/commands/options"
	"golang.org/x/mod/modfile"
	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/files"
	"knative.dev/toolbox/magetasks/pkg/ldflags"
	"knative.dev/toolbox/magetasks/pkg/output/color"
)

const (
	koImportPath  = "ko.import.path"
	koBuildResult = "ko.build.result"
)

// ErrKoFailed when th Google's ko fails to build.
var ErrKoFailed = errors.New("ko failed")

// KoBuilder builds images with Google's KO.
type KoBuilder struct{}

type imageReference struct {
	ociName name.Reference
	tags    []string
}

func (kb KoBuilder) Accepts(artifact config.Artifact) bool {
	_, ok := artifact.(Image)
	return ok
}

func (kb KoBuilder) Build(artifact config.Artifact, notifier config.Notifier) config.Result {
	image, ok := artifact.(Image)
	if !ok {
		return config.Result{Error: ErrInvalidArtifact}
	}
	importPath, err := imageImportPath(image)
	if err != nil {
		return resultErrKoFailed(err)
	}
	bo := &options.BuildOptions{
		Platforms: buildPlatforms(image),
		Labels:    buildLabels(image, importPath),
	}
	fillInLdflags(bo, importPath, image)
	ctx := config.Actual().Context
	builder, err := commands.NewBuilder(ctx, bo)
	if err != nil {
		return resultErrKoFailed(err)
	}
	result, err := builder.Build(ctx, importPath)
	if err != nil {
		return resultErrKoFailed(err)
	}
	ref, err := calculateImageReference(result, artifact)
	if err != nil {
		return resultErrKoFailed(err)
	}
	notifier.Notify(fmt.Sprintf("ko built image: %s", color.Blue(ref.ociName)))
	notifier.Notify(fmt.Sprintf("image tags: %s", color.Blue(fmt.Sprintf("%+q", ref.tags))))
	return config.Result{Info: map[string]interface{}{
		imageReferenceKey: ref.ociName.String(),
		koBuildResult:     result,
	}}
}

func fillInLdflags(bo *options.BuildOptions, importPath string, image Image) {
	c := config.Actual()
	args := make([]string, 0)
	if c.Version != nil || len(c.BuildVariables) > 0 || len(image.BuildVariables) > 0 {
		builder := ldflags.NewBuilder()
		for key, resolver := range c.BuildVariables {
			builder.Add(key, resolver)
		}
		if c.Version != nil {
			builder.Add(c.Version.Path, c.Version.Resolver.Version)
		}
		for key, resolver := range image.BuildVariables {
			builder.Add(key, resolver)
		}
		args = builder.Build()
	}
	if len(args) > 0 {
		bo.BuildConfigs = map[string]build.Config{
			importPath: {
				ID:      "ldflags-config",
				Ldflags: args,
			},
		}
	}
}

func buildLabels(image Image, importPath string) []string {
	labels := make([]string, 0, len(image.Labels))
	if version := config.Actual().Version; version != nil {
		labels = append(labels, fmt.Sprintf("version=%s", version.Resolver.Version()))
	}
	labels = append(labels, fmt.Sprintf("%s=%s", koImportPath, importPath))
	for key, resolver := range image.Labels {
		labels = append(labels, fmt.Sprintf("%s=%s", key, resolver()))
	}
	return labels
}

func buildPlatforms(im Image) []string {
	platforms := make([]string, len(im.Architectures))
	for i, architecture := range im.Architectures {
		platforms[i] = fmt.Sprintf("linux/%s", architecture)
	}
	return platforms
}

func calculateImageReference(result build.Result, artifact config.Artifact) (*imageReference, error) {
	kp := KoPublisher{}
	po, err := kp.publishOptions()
	if err != nil {
		return nil, err
	}
	po.Push = false
	po.TarballFile = ""
	po.OCILayoutPath = ""
	po.Local = false
	publisher, err := commands.NewPublisher(po)
	if err != nil {
		return nil, err
	}
	defer closePublisher(publisher)
	ctx := config.Actual().Context
	ref, err := publisher.Publish(ctx, result, artifact.GetName())
	if err != nil {
		return nil, err
	}
	return &imageReference{
		ociName: ref, tags: po.Tags,
	}, nil
}

func resultErrKoFailed(err error) config.Result {
	return config.Result{
		Error: fmt.Errorf("%w: %w", ErrKoFailed, err),
	}
}

func imageImportPath(image Image) (string, error) {
	binDir := fullBinaryDirectory(image.GetName())
	rs, err := lookForGoModule(binDir)
	if err != nil {
		return "", err
	}
	importPath := rs.resolve(binDir)
	if resolver, ok := image.Labels[koImportPath]; ok {
		importPath = resolver()
	}
	return importPath, nil
}

func lookForGoModule(dir string) (lookupGoModuleResult, error) {
	rs := lookupGoModuleResult{}
	for i := 0; i < 10_000; i++ {
		modFile := path.Join(dir, "go.mod")
		if files.DontExists(modFile) {
			dir = path.Dir(dir)
			rs.directoryDistance++
			continue
		}
		bytes, err := os.ReadFile(modFile)
		if err != nil {
			return rs, err
		}
		file, err := modfile.Parse(modFile, bytes, nil)
		if err != nil {
			return rs, err
		}
		rs.module = file
		return rs, nil
	}
	return rs, fmt.Errorf("%w: can't find go module", ErrKoFailed)
}

type lookupGoModuleResult struct {
	module            *modfile.File
	directoryDistance int
}

func (r lookupGoModuleResult) resolve(dir string) string {
	root := dir
	for i := 0; i < r.directoryDistance; i++ {
		root = path.Dir(root)
	}
	p := strings.Replace(dir, root, "", 1)
	return path.Join(r.module.Module.Mod.Path, p)
}
