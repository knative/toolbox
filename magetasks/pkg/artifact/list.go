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

	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/files"
	"knative.dev/toolbox/magetasks/pkg/output/color"
)

const (
	// ArtifactsBuilt is used to list artifacts that was built.
	ArtifactsBuilt     = "artifacts.built"
	allowReadAllOsPerm = 0o644
)

// ErrMisconfiguration when the project configuration faulty.
var ErrMisconfiguration = errors.New("project configuration is faulty")

// ListPublisher will output built artifacts as a "\n" delimited list in
// a result file.
type ListPublisher struct {
	FilePath         string
	ResultsRetriever func(config.Artifact) *config.Result
}

func (l ListPublisher) Accepts(artifact config.Artifact) bool {
	return l.resultsRetriever()(artifact) != nil
}

func (l ListPublisher) Publish(artifact config.Artifact, notifier config.Notifier) config.Result {
	reportFilename := l.FilePath
	if reportFilename == "" {
		reportFilename = "artifacts.list"
	}
	reportPath := path.Join(files.BuildDir(), reportFilename)
	result := l.resultsRetriever()(artifact)
	if result == nil && !result.Failed() {
		return config.Result{Error: fmt.Errorf(
			"%w: can't find result for %v", ErrMisconfiguration, artifact)}
	}
	artifactsList, ok := result.Info[ArtifactsBuilt].([]string)
	if !ok {
		return config.Result{Error: fmt.Errorf(
			"%w: can't find result for %v", ErrMisconfiguration, artifact)}
	}

	//nolint:nosnakecase
	flgs := os.O_APPEND | os.O_WRONLY | os.O_CREATE
	f, err := os.OpenFile(reportPath, flgs, allowReadAllOsPerm)
	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()
	if err != nil {
		return config.Result{Error: fmt.Errorf("%w: %w", ErrMisconfiguration, err)}
	}

	for _, artifactPath := range artifactsList {
		if _, err = f.WriteString(fmt.Sprintf("%s\n", artifactPath)); err != nil {
			return config.Result{Error: fmt.Errorf("%w: %w", ErrMisconfiguration, err)}
		}
	}

	notifier.Notify(fmt.Sprintf("artifact list: %s", color.Blue(reportPath)))

	return config.Result{}
}

func (l ListPublisher) resultsRetriever() func(config.Artifact) *config.Result {
	if l.ResultsRetriever != nil {
		return l.ResultsRetriever
	}
	return func(artifact config.Artifact) *config.Result {
		if !BinaryBuilder.Accepts(BinaryBuilder{}, artifact) {
			return nil
		}
		key := BuildKey(artifact)
		result, ok := config.Actual().Context.Value(key).(config.Result)
		if !ok {
			return nil
		}
		return &result
	}
}
