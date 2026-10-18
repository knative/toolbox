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

package config

import (
	"context"
)

var (
	// DefaultBuilders is a list of default builders.
	DefaultBuilders = make([]Builder, 0)
	// DefaultPublishers is a list of default publishers.
	DefaultPublishers = make([]Publisher, 0)
)

// FillInDefaultValues in provided config and returns a filled one.
func FillInDefaultValues(cfg Config) Config {
	if len(cfg.BuildDirPath) == 0 {
		cfg.BuildDirPath = []string{"build", "_output"}
	}
	if cfg.Dependencies.Golang == nil {
		cfg.Dependencies.Golang = NewDependencies("gotest.tools/gotestsum@latest")
	}
	if cfg.Dependencies.Binaries == nil {
		cfg.Dependencies.Binaries = NewBinaries()
	}
	if cfg.Context == nil {
		cfg.Context = context.TODO()
	}
	if cfg.Artifacts == nil {
		cfg.Artifacts = make([]Artifact, 0)
	}
	if len(cfg.Builders) == 0 {
		cfg.Builders = append(cfg.Builders, DefaultBuilders...)
	}
	if len(cfg.Publishers) == 0 {
		cfg.Publishers = append(cfg.Publishers, DefaultPublishers...)
	}
	return cfg
}
